package usecase

import (
	"fmt"
	"math"
	"time"

	"dashboard-cs-be/entities"
	"dashboard-cs-be/repository/interfaces"
)

const dateLayout = "2006-01-02"

type dashboardUsecase struct {
	repo interfaces.DashboardRepository
}

func NewDashboardUsecase(repo interfaces.DashboardRepository) DashboardUsecase {
	return &dashboardUsecase{repo: repo}
}

// ─── GetRealtime ──────────────────────────────────────────────────────────────

func (uc *dashboardUsecase) GetRealtime() (*entities.RealtimeResponse, error) {
	rt, err := uc.repo.GetRealtime()
	if err != nil {
		return nil, fmt.Errorf("usecase GetRealtime: %w", err)
	}

	// Delta created: hari ini vs kemarin
	yesterday := time.Now().AddDate(0, 0, -1).Format(dateLayout)
	today := time.Now().Format(dateLayout)

	summaryToday, err := uc.repo.GetSummary(today, today)
	if err != nil {
		return nil, fmt.Errorf("usecase GetRealtime summary today: %w", err)
	}
	summaryYesterday, err := uc.repo.GetSummary(yesterday, yesterday)
	if err != nil {
		return nil, fmt.Errorf("usecase GetRealtime summary yesterday: %w", err)
	}
	createdDelta := summaryToday.TotalTickets - summaryYesterday.TotalTickets

	return &entities.RealtimeResponse{
		SLAToday: entities.SLAToday{
			Percentage: rt.SLATodayPercentage,
			Delta:      rt.SLATodayDelta,
		},
		CreatedToday: entities.Created{
			Total: rt.CreatedTodayTotal,
			Delta: createdDelta,
		},
		IncidentsActive: 0,
	}, nil
}

// ─── GetDashboard ─────────────────────────────────────────────────────────────

func (uc *dashboardUsecase) GetDashboard(from, to string) (*entities.DashboardResponse, error) {
	if err := validateDateRange(from, to); err != nil {
		return nil, err
	}

	// Daily trend: jika from == to, ambil 7 hari ke belakang
	trendFrom := from
	if from == to {
		t, _ := time.Parse(dateLayout, to)
		trendFrom = t.AddDate(0, 0, -6).Format(dateLayout)
	}

	summary, err := uc.repo.GetSummary(from, to)
	if err != nil {
		return nil, fmt.Errorf("usecase GetDashboard summary: %w", err)
	}

	trend, err := uc.repo.GetDailyTrend(trendFrom, to)
	if err != nil {
		return nil, fmt.Errorf("usecase GetDashboard trend: %w", err)
	}

	hourly, err := uc.repo.GetTicketsPerHour(from, to)
	if err != nil {
		return nil, fmt.Errorf("usecase GetDashboard hourly: %w", err)
	}

	priority, err := uc.repo.GetPrioritySummary(from, to)
	if err != nil {
		return nil, fmt.Errorf("usecase GetDashboard priority: %w", err)
	}

	channelSLA, err := uc.repo.GetChannelSLA(from, to)
	if err != nil {
		return nil, fmt.Errorf("usecase GetDashboard channelSLA: %w", err)
	}

	return &entities.DashboardResponse{
		Filter: entities.FilterMeta{From: from, To: to},
		Summary: entities.SummaryAggregate{
			TotalTickets:   summary.TotalTickets,
			Open:           summary.Open,
			Closed:         summary.Closed,
			CSATPercentage: summary.CSATPercentage,
			CSATScore:      summary.CSATScore,
			Unassigned:     summary.Unassigned,
		},
		DailyTrend:      mapTrend(trend),
		TicketsPerHour:  mapHourly(hourly),
		PrioritySummary: mapPriority(priority),
		Channels:        mapChannels(channelSLA),
	}, nil
}

// ─── GetChannelDetail ─────────────────────────────────────────────────────────

func (uc *dashboardUsecase) GetChannelDetail(f entities.ChannelDetailFilter) (*entities.ChannelDetailResponse, error) {
	if err := validateDateRange(f.From, f.To); err != nil {
		return nil, err
	}

	// Sanitize pagination
	if f.Page < 1 {
		f.Page = 1
	}
	if f.Limit < 1 || f.Limit > 100 {
		f.Limit = 10
	}

	// Fetch top corporate & top KIP secara paralel
	type corpResult struct {
		rows []entities.TopCorporateRow
		err  error
	}
	type kipResult struct {
		rows []entities.TopKIPRow
		err  error
	}

	corpCh := make(chan corpResult, 1)
	kipCh := make(chan kipResult, 1)

	go func() {
		rows, err := uc.repo.GetTopCorporate(f)
		corpCh <- corpResult{rows, err}
	}()
	go func() {
		rows, err := uc.repo.GetTopKIP(f)
		kipCh <- kipResult{rows, err}
	}()

	corpRes := <-corpCh
	kipRes := <-kipCh

	if corpRes.err != nil {
		return nil, fmt.Errorf("usecase GetChannelDetail corporate: %w", corpRes.err)
	}
	if kipRes.err != nil {
		return nil, fmt.Errorf("usecase GetChannelDetail kip: %w", kipRes.err)
	}

	return &entities.ChannelDetailResponse{
		Filter: f,
		TopCorporate: entities.PaginatedCorporate{
			Data:       mapCorporate(corpRes.rows),
			Pagination: buildPagination(f.Page, f.Limit, totalCount(corpRes.rows)),
		},
		TopKIP: entities.PaginatedKIP{
			Data:       mapKIP(kipRes.rows),
			Pagination: buildPagination(f.Page, f.Limit, totalCount2(kipRes.rows)),
		},
	}, nil
}

// ─── Helpers ─────────────────────────────────────────────────────────────────

func mapTrend(rows []entities.TrendRow) []entities.DailyTrendPoint {
	out := make([]entities.DailyTrendPoint, len(rows))
	for i, r := range rows {
		out[i] = entities.DailyTrendPoint{Date: r.Date, Created: r.Created, Solved: r.Solved}
	}
	return out
}

func mapHourly(rows []entities.HourlyRow) []entities.TicketHourPoint {
	out := make([]entities.TicketHourPoint, len(rows))
	for i, r := range rows {
		out[i] = entities.TicketHourPoint{Hour: r.Hour, Created: r.Created, Solved: r.Solved}
	}
	return out
}

func mapPriority(p *entities.PriorityRow) entities.PrioritySummary {
	return entities.PrioritySummary{
		Roaming: p.Roaming, ExtraQuota: p.ExtraQuota,
		CC: p.CC, VIP: p.VIP, P1: p.P1, Urgent: p.Urgent,
	}
}

func mapChannels(rows []entities.ChannelSLARow) []entities.ChannelSummary {
	out := make([]entities.ChannelSummary, len(rows))
	for i, r := range rows {
		out[i] = entities.ChannelSummary{
			Channel: r.Channel,
			SLA:     r.SLA,
			Open:    r.Open,
			Closed:  r.Closed,
		}
	}
	return out
}

func mapCorporate(rows []entities.TopCorporateRow) []entities.TopCorporate {
	out := make([]entities.TopCorporate, len(rows))
	for i, r := range rows {
		out[i] = entities.TopCorporate{
			CompanyName:   r.CompanyName,
			Interactions:  r.Interactions,
			Tickets:       r.Tickets,
			FCRPercentage: r.FCRPercentage,
		}
	}
	return out
}

func mapKIP(rows []entities.TopKIPRow) []entities.TopKIP {
	out := make([]entities.TopKIP, len(rows))
	for i, r := range rows {
		out[i] = entities.TopKIP{
			Topic:         r.Topic,
			Interactions:  r.Interactions,
			Tickets:       r.Tickets,
			FCRPercentage: r.FCRPercentage,
		}
	}
	return out
}

func buildPagination(page, limit, total int) entities.Pagination {
	totalPages := int(math.Ceil(float64(total) / float64(limit)))
	if totalPages < 1 {
		totalPages = 1
	}
	return entities.Pagination{
		Page:       page,
		Limit:      limit,
		TotalItems: total,
		TotalPages: totalPages,
	}
}

func totalCount(rows []entities.TopCorporateRow) int {
	if len(rows) == 0 {
		return 0
	}
	return rows[0].Total
}

func totalCount2(rows []entities.TopKIPRow) int {
	if len(rows) == 0 {
		return 0
	}
	return rows[0].Total
}

func validateDateRange(from, to string) error {
	f, err := time.Parse(dateLayout, from)
	if err != nil {
		return fmt.Errorf("invalid from date %q: use YYYY-MM-DD", from)
	}
	t, err := time.Parse(dateLayout, to)
	if err != nil {
		return fmt.Errorf("invalid to date %q: use YYYY-MM-DD", to)
	}
	if f.After(t) {
		return fmt.Errorf("from date must be <= to date")
	}
	return nil
}