package usecase

import (
	"fmt"
	"time"

	"dashboard-cs-be/entities"
	"dashboard-cs-be/repository"
)

const dateLayout = "2006-01-02"

type dashboardUsecase struct {
	repo repository.DashboardRepository
}

func NewDashboardUsecase(repo repository.DashboardRepository) DashboardUsecase {
	return &dashboardUsecase{repo: repo}
}

// ─────────────────────────────────────────────────────────────────────────────
// GetRealtime
// ─────────────────────────────────────────────────────────────────────────────

func (uc *dashboardUsecase) GetRealtime() (*entities.RealtimeResponse, error) {
	rt, err := uc.repo.GetRealtime()
	if err != nil {
		return nil, fmt.Errorf("usecase GetRealtime: %w", err)
	}

	// Delta created: hari ini vs kemarin — query terpisah
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
		OpenTickets:     rt.OpenTickets,
		Unassigned:      rt.Unassigned,
		IncidentsActive: 0,
	}, nil
}

// ─────────────────────────────────────────────────────────────────────────────
// GetDashboard
// ─────────────────────────────────────────────────────────────────────────────

func (uc *dashboardUsecase) GetDashboard(from, to string) (*entities.DashboardResponse, error) {
	if err := validateDateRange(from, to); err != nil {
		return nil, err
	}

	// Run all queries (could be parallelised with goroutines later)
	summary, err := uc.repo.GetSummary(from, to)
	if err != nil {
		return nil, fmt.Errorf("usecase GetDashboard summary: %w", err)
	}

	trend, err := uc.repo.GetDailyTrend(from, to)
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

	topCorp, err := uc.repo.GetTopCorporate(from, to)
	if err != nil {
		return nil, fmt.Errorf("usecase GetDashboard topCorp: %w", err)
	}

	topKIP, err := uc.repo.GetTopKIP(from, to)
	if err != nil {
		return nil, fmt.Errorf("usecase GetDashboard topKIP: %w", err)
	}

	return &entities.DashboardResponse{
		Filter: entities.FilterMeta{From: from, To: to},
		Summary: entities.SummaryAggregate{
			TotalTickets:   summary.TotalTickets,
			Open:           summary.Open,
			Closed:         summary.Closed,
			CSATPercentage: summary.CSATPercentage,
			CSATScore:      summary.CSATScore,
		},
		DailyTrend:      mapTrend(trend),
		TicketsPerHour:  mapHourly(hourly),
		PrioritySummary: mapPriority(priority),
		Channels:        assembleChannels(channelSLA, topCorp, topKIP),
	}, nil
}

// ─────────────────────────────────────────────────────────────────────────────
// Mappers
// ─────────────────────────────────────────────────────────────────────────────

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
		Roaming:    p.Roaming,
		ExtraQuota: p.ExtraQuota,
		CC:         p.CC,
		VIP:        p.VIP,
		P1:         p.P1,
		Urgent:     p.Urgent,
	}
}

func assembleChannels(
	slaRows []entities.ChannelSLARow,
	corpRows []entities.TopCorporateRow,
	kipRows []entities.TopKIPRow,
) []entities.ChannelDetail {
	details := make([]entities.ChannelDetail, 0, len(slaRows))

	for _, ch := range slaRows {
		detail := entities.ChannelDetail{
			Channel: ch.Channel,
			SLA:     ch.SLA,
			Open:    ch.Open,
			Closed:  ch.Closed,
		}

		// attach top_corporate for this channel
		for _, c := range corpRows {
			if c.Channel == ch.Channel {
				detail.TopCorporate = append(detail.TopCorporate, entities.TopCorporate{
					CompanyName:   c.CompanyName,
					Interactions:  c.Interactions,
					Tickets:       c.Tickets,
					FCRPercentage: c.FCRPercentage,
				})
			}
		}

		// attach top_kip for this channel
		for _, k := range kipRows {
			if k.Channel == ch.Channel {
				detail.TopKIP = append(detail.TopKIP, entities.TopKIP{
					Topic:         k.Topic,
					Interactions:  k.Interactions,
					Tickets:       k.Tickets,
					FCRPercentage: k.FCRPercentage,
				})
			}
		}

		details = append(details, detail)
	}
	return details
}

// ─────────────────────────────────────────────────────────────────────────────
// Validation
// ─────────────────────────────────────────────────────────────────────────────

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
