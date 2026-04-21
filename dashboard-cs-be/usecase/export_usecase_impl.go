package usecase

import (
	"fmt"
	"io"
	"strings"

	"github.com/xuri/excelize/v2"

	"dashboard-cs-be/entities"
	"dashboard-cs-be/repository/interfaces"
)

const (
	colorHeaderBg   = "1F4E79"
	colorHeaderFont = "FFFFFF"
	colorSectionBg  = "2E75B6"
	colorSectionFont = "FFFFFF"
	colorMetaBg     = "D6E4F0"
	colorAltRow     = "EBF3FB"
	colorGoodSLA    = "E2EFDA"
	colorWarnSLA    = "FFF2CC"
	colorBadSLA     = "FCE4D6"
	colorFontDark   = "1F1F1F"
)

type exportUsecase struct {
	repo interfaces.ExportRepository
}

func NewExportUsecase(repo interfaces.ExportRepository) ExportUsecase {
	return &exportUsecase{repo: repo}
}

func (uc *exportUsecase) ExportExcel(filter entities.ExportFilter, w io.Writer) (string, error) {
	summary, err := uc.repo.GetExportSummary(filter)
	if err != nil {
		return "", fmt.Errorf("export summary: %w", err)
	}
	channels, err := uc.repo.GetExportChannels(filter)
	if err != nil {
		return "", fmt.Errorf("export channels: %w", err)
	}
	customers, err := uc.repo.GetExportCustomers(filter)
	if err != nil {
		return "", fmt.Errorf("export customers: %w", err)
	}
	topics, err := uc.repo.GetExportTopics(filter)
	if err != nil {
		return "", fmt.Errorf("export topics: %w", err)
	}
	priorities, err := uc.repo.GetExportPriorities(filter)
	if err != nil {
		return "", fmt.Errorf("export priorities: %w", err)
	}

	payload := &entities.ExportPayload{
		Filter:     filter,
		Summary:    *summary,
		Channels:   channels,
		Customers:  customers,
		Topics:     topics,
		Priorities: priorities,
	}

	f := excelize.NewFile()
	defer f.Close()

	sg := &sheetGen{f: f}

	if err := sg.writeSheet1Summary(payload); err != nil {
		return "", err
	}
	if err := sg.writeSheet2Channels(payload); err != nil {
		return "", err
	}
	if err := sg.writeSheet3Customers(payload); err != nil {
		return "", err
	}
	if err := sg.writeSheet4Topics(payload); err != nil {
		return "", err
	}
	if err := sg.writeSheet5Priorities(payload); err != nil {
		return "", err
	}

	f.DeleteSheet("Sheet1")

	if err := f.Write(w); err != nil {
		return "", fmt.Errorf("write excel: %w", err)
	}

	return buildFilename(filter), nil
}

// ─────────────────────────────────────────────────────────────────────────────
// sheetGen
// ─────────────────────────────────────────────────────────────────────────────

type sheetGen struct {
	f      *excelize.File
	styles map[string]int
}

func (sg *sheetGen) style(key string, s *excelize.Style) int {
	if sg.styles == nil {
		sg.styles = make(map[string]int)
	}
	if id, ok := sg.styles[key]; ok {
		return id
	}
	id, _ := sg.f.NewStyle(s)
	sg.styles[key] = id
	return id
}

func (sg *sheetGen) sHeader() int {
	return sg.style("header", &excelize.Style{
		Fill:      excelize.Fill{Type: "pattern", Color: []string{colorHeaderBg}, Pattern: 1},
		Font:      &excelize.Font{Bold: true, Color: colorHeaderFont, Size: 10},
		Alignment: &excelize.Alignment{Horizontal: "center", Vertical: "center", WrapText: true},
		Border:    thinBorder(),
	})
}

func (sg *sheetGen) sSection() int {
	return sg.style("section", &excelize.Style{
		Fill:      excelize.Fill{Type: "pattern", Color: []string{colorSectionBg}, Pattern: 1},
		Font:      &excelize.Font{Bold: true, Color: colorSectionFont, Size: 11},
		Alignment: &excelize.Alignment{Vertical: "center"},
	})
}

func (sg *sheetGen) sMeta() int {
	return sg.style("meta", &excelize.Style{
		Fill:      excelize.Fill{Type: "pattern", Color: []string{colorMetaBg}, Pattern: 1},
		Font:      &excelize.Font{Bold: true, Color: colorFontDark, Size: 10},
		Alignment: &excelize.Alignment{Vertical: "center"},
		Border:    thinBorder(),
	})
}

func (sg *sheetGen) sMetaValue() int {
	return sg.style("metaval", &excelize.Style{
		Font:      &excelize.Font{Color: colorFontDark, Size: 10},
		Alignment: &excelize.Alignment{Vertical: "center"},
		Border:    thinBorder(),
	})
}

func (sg *sheetGen) sData(alt bool) int {
	key := "data"
	bg := "FFFFFF"
	if alt {
		key = "dataAlt"
		bg = colorAltRow
	}
	return sg.style(key, &excelize.Style{
		Fill:      excelize.Fill{Type: "pattern", Color: []string{bg}, Pattern: 1},
		Font:      &excelize.Font{Color: colorFontDark, Size: 10},
		Alignment: &excelize.Alignment{Vertical: "center"},
		Border:    thinBorder(),
	})
}

func (sg *sheetGen) sDataCenter(alt bool) int {
	key := "dataC"
	bg := "FFFFFF"
	if alt {
		key = "dataCalt"
		bg = colorAltRow
	}
	return sg.style(key, &excelize.Style{
		Fill:      excelize.Fill{Type: "pattern", Color: []string{bg}, Pattern: 1},
		Font:      &excelize.Font{Color: colorFontDark, Size: 10},
		Alignment: &excelize.Alignment{Horizontal: "center", Vertical: "center"},
		Border:    thinBorder(),
	})
}

func (sg *sheetGen) sSLAColor(pct float64, alt bool) int {
	bg := colorGoodSLA
	if pct < 50 {
		bg = colorBadSLA
	} else if pct < 80 {
		bg = colorWarnSLA
	}
	key := fmt.Sprintf("sla%.0f_%v", pct, alt)
	return sg.style(key, &excelize.Style{
		Fill:      excelize.Fill{Type: "pattern", Color: []string{bg}, Pattern: 1},
		Font:      &excelize.Font{Bold: true, Color: colorFontDark, Size: 10},
		Alignment: &excelize.Alignment{Horizontal: "center", Vertical: "center"},
		Border:    thinBorder(),
	})
}

func thinBorder() []excelize.Border {
	return []excelize.Border{
		{Type: "left", Color: "BFBFBF", Style: 1},
		{Type: "right", Color: "BFBFBF", Style: 1},
		{Type: "top", Color: "BFBFBF", Style: 1},
		{Type: "bottom", Color: "BFBFBF", Style: 1},
	}
}

// ─────────────────────────────────────────────────────────────────────────────
// Sheet 1 — Summary Overview
// CSAT = 100% jika ada tiket masuk, 0% jika tidak ada
// ─────────────────────────────────────────────────────────────────────────────

func (sg *sheetGen) writeSheet1Summary(p *entities.ExportPayload) error {
	name := "Summary Overview"
	sg.f.NewSheet(name)

	for col, w := range map[string]float64{"A": 28, "B": 22, "C": 18, "D": 18} {
		sg.f.SetColWidth(name, col, col, w)
	}

	row := 1

	sg.f.SetCellStyle(name, cell("A", row), cell("A", row), sg.sMeta())
	sg.f.SetCellStyle(name, cell("B", row), cell("D", row), sg.sMetaValue())
	sg.f.SetCellValue(name, cell("A", row), "Period")
	sg.f.MergeCell(name, cell("B", row), cell("D", row))
	sg.f.SetCellValue(name, cell("B", row), fmt.Sprintf("%s  s/d  %s", p.Filter.From, p.Filter.To))
	row++

	sg.f.SetCellStyle(name, cell("A", row), cell("A", row), sg.sMeta())
	sg.f.SetCellStyle(name, cell("B", row), cell("D", row), sg.sMetaValue())
	sg.f.SetCellValue(name, cell("A", row), "Channel Filter")
	sg.f.MergeCell(name, cell("B", row), cell("D", row))
	sg.f.SetCellValue(name, cell("B", row), channelLabel(p.Filter.Channels))
	row++
	row++ // blank

	sg.f.SetCellStyle(name, cell("A", row), cell("D", row), sg.sSection())
	sg.f.MergeCell(name, cell("A", row), cell("D", row))
	sg.f.SetCellValue(name, cell("A", row), "  OVERALL PERFORMANCE")
	sg.f.SetRowHeight(name, row, 22)
	row++

	// SLA dari channel data
	slaPct := 0.0
	if len(p.Channels) > 0 {
		totalT := 0
		totalAchieved := 0.0
		for _, ch := range p.Channels {
			totalT += ch.TotalTickets
			totalAchieved += ch.SLAPercent * float64(ch.TotalTickets)
		}
		if totalT > 0 {
			slaPct = totalAchieved / float64(totalT)
		}
	}
	slaBreached := 0.0
	if slaPct > 0 {
		slaBreached = 100 - slaPct
	}

	// CSAT = 100% berbasis interaksi masuk
	csatPct := 0.0
	if p.Summary.TotalTickets > 0 {
		csatPct = 100.0
	}

	metrics := [][]interface{}{
		{"Total Tickets", p.Summary.TotalTickets},
		{"Open Tickets", p.Summary.Open},
		{"Closed Tickets", p.Summary.Closed},
		{"SLA Achievement", fmt.Sprintf("%.2f%%", slaPct)},
		{"SLA Breached", fmt.Sprintf("%.2f%%", slaBreached)},
		{"CSAT", fmt.Sprintf("%.2f%%", csatPct)},
	}

	for _, m := range metrics {
		sg.f.SetCellStyle(name, cell("A", row), cell("A", row), sg.sMeta())
		sg.f.SetCellStyle(name, cell("B", row), cell("D", row), sg.sMetaValue())
		sg.f.SetCellValue(name, cell("A", row), m[0])
		sg.f.MergeCell(name, cell("B", row), cell("D", row))
		sg.f.SetCellValue(name, cell("B", row), m[1])
		sg.f.SetRowHeight(name, row, 18)
		row++
	}

	return nil
}

// ─────────────────────────────────────────────────────────────────────────────
// Sheet 2 — Channel Performance (hapus kolom FCR%)
// ─────────────────────────────────────────────────────────────────────────────

func (sg *sheetGen) writeSheet2Channels(p *entities.ExportPayload) error {
	name := "Channel Performance"
	sg.f.NewSheet(name)

	colWidths := map[string]float64{
		"A": 18, "B": 16, "C": 10, "D": 12, "E": 10,
	}
	for col, w := range colWidths {
		sg.f.SetColWidth(name, col, col, w)
	}

	headers := []string{"Channel", "Total Tickets", "Open", "Closed", "SLA%"}
	for i, h := range headers {
		c := colLetter(i)
		sg.f.SetCellStyle(name, cell(c, 1), cell(c, 1), sg.sHeader())
		sg.f.SetCellValue(name, cell(c, 1), h)
	}
	sg.f.SetRowHeight(name, 1, 22)

	for i, ch := range p.Channels {
		r := i + 2
		alt := i%2 == 1
		sg.f.SetCellStyle(name, cell("A", r), cell("A", r), sg.sData(alt))
		sg.f.SetCellStyle(name, cell("B", r), cell("D", r), sg.sDataCenter(alt))
		sg.f.SetCellStyle(name, cell("E", r), cell("E", r), sg.sSLAColor(ch.SLAPercent, alt))

		sg.f.SetCellValue(name, cell("A", r), labelChannel(ch.Channel))
		sg.f.SetCellValue(name, cell("B", r), ch.TotalTickets)
		sg.f.SetCellValue(name, cell("C", r), ch.Open)
		sg.f.SetCellValue(name, cell("D", r), ch.Closed)
		sg.f.SetCellValue(name, cell("E", r), fmt.Sprintf("%.2f%%", ch.SLAPercent))
		sg.f.SetRowHeight(name, r, 18)
	}

	return nil
}

// ─────────────────────────────────────────────────────────────────────────────
// Sheet 3 — Consumer Report (tidak ada perubahan)
// ─────────────────────────────────────────────────────────────────────────────

func (sg *sheetGen) writeSheet3Customers(p *entities.ExportPayload) error {
	name := "Consumer Report"
	sg.f.NewSheet(name)

	colWidths := map[string]float64{
		"A": 16, "B": 30, "C": 16, "D": 12,
		"E": 14, "F": 10, "G": 12, "H": 14, "I": 14,
	}
	for col, w := range colWidths {
		sg.f.SetColWidth(name, col, col, w)
	}

	headers := []string{
		"Customer ID", "Name", "Phone", "Type",
		"Total Tickets", "Open", "Closed", "SLA Achieved", "SLA Breached",
	}
	for i, h := range headers {
		c := colLetter(i)
		sg.f.SetCellStyle(name, cell(c, 1), cell(c, 1), sg.sHeader())
		sg.f.SetCellValue(name, cell(c, 1), h)
	}
	sg.f.SetRowHeight(name, 1, 22)

	for i, cust := range p.Customers {
		r := i + 2
		alt := i%2 == 1
		ds := sg.sData(alt)
		dc := sg.sDataCenter(alt)

		sg.f.SetCellStyle(name, cell("A", r), cell("D", r), ds)
		sg.f.SetCellStyle(name, cell("E", r), cell("I", r), dc)

		sg.f.SetCellValue(name, cell("A", r), cust.CustomerID)
		sg.f.SetCellValue(name, cell("B", r), cust.Name)
		sg.f.SetCellValue(name, cell("C", r), cust.Phone)
		sg.f.SetCellValue(name, cell("D", r), strings.Title(cust.CustomerType))
		sg.f.SetCellValue(name, cell("E", r), cust.Total)
		sg.f.SetCellValue(name, cell("F", r), cust.Open)
		sg.f.SetCellValue(name, cell("G", r), cust.Closed)
		sg.f.SetCellValue(name, cell("H", r), cust.SLAAchieved)
		sg.f.SetCellValue(name, cell("I", r), cust.SLABreached)
		sg.f.SetRowHeight(name, r, 18)
	}

	return nil
}

// ─────────────────────────────────────────────────────────────────────────────
// Sheet 4 — Topic/KIP Report (hapus kolom FCR%)
// ─────────────────────────────────────────────────────────────────────────────

func (sg *sheetGen) writeSheet4Topics(p *entities.ExportPayload) error {
	name := "Topic - KIP Report"
	sg.f.NewSheet(name)

	colWidths := map[string]float64{
		"A": 38, "B": 16, "C": 10, "D": 10, "E": 12,
	}
	for col, w := range colWidths {
		sg.f.SetColWidth(name, col, col, w)
	}

	headers := []string{"Topic", "Channel", "Total", "Open", "Closed"}
	for i, h := range headers {
		c := colLetter(i)
		sg.f.SetCellStyle(name, cell(c, 1), cell(c, 1), sg.sHeader())
		sg.f.SetCellValue(name, cell(c, 1), h)
	}
	sg.f.SetRowHeight(name, 1, 22)

	for i, t := range p.Topics {
		r := i + 2
		alt := i%2 == 1
		sg.f.SetCellStyle(name, cell("A", r), cell("B", r), sg.sData(alt))
		sg.f.SetCellStyle(name, cell("C", r), cell("E", r), sg.sDataCenter(alt))

		sg.f.SetCellValue(name, cell("A", r), t.Topic)
		sg.f.SetCellValue(name, cell("B", r), labelChannel(t.Channel))
		sg.f.SetCellValue(name, cell("C", r), t.Total)
		sg.f.SetCellValue(name, cell("D", r), t.Open)
		sg.f.SetCellValue(name, cell("E", r), t.Closed)
		sg.f.SetRowHeight(name, r, 18)
	}

	return nil
}

// ─────────────────────────────────────────────────────────────────────────────
// Sheet 5 — Priority Report (tidak ada perubahan)
// ─────────────────────────────────────────────────────────────────────────────

func (sg *sheetGen) writeSheet5Priorities(p *entities.ExportPayload) error {
	name := "Priority Report"
	sg.f.NewSheet(name)

	colWidths := map[string]float64{
		"A": 14, "B": 10, "C": 10, "D": 12,
		"E": 14, "F": 14, "G": 22,
	}
	for col, w := range colWidths {
		sg.f.SetColWidth(name, col, col, w)
	}

	headers := []string{
		"Priority", "Total", "Open", "Closed",
		"SLA Achieved", "SLA Breached", "Avg Resolution Time",
	}
	for i, h := range headers {
		c := colLetter(i)
		sg.f.SetCellStyle(name, cell(c, 1), cell(c, 1), sg.sHeader())
		sg.f.SetCellValue(name, cell(c, 1), h)
	}
	sg.f.SetRowHeight(name, 1, 22)

	for i, pr := range p.Priorities {
		r := i + 2
		alt := i%2 == 1
		dc := sg.sDataCenter(alt)

		sg.f.SetCellStyle(name, cell("A", r), cell("G", r), dc)

		sg.f.SetCellValue(name, cell("A", r), strings.ToUpper(pr.Priority))
		sg.f.SetCellValue(name, cell("B", r), pr.Total)
		sg.f.SetCellValue(name, cell("C", r), pr.Open)
		sg.f.SetCellValue(name, cell("D", r), pr.Closed)
		sg.f.SetCellValue(name, cell("E", r), pr.SLAAchieved)
		sg.f.SetCellValue(name, cell("F", r), pr.SLABreached)
		sg.f.SetCellValue(name, cell("G", r), formatDuration(pr.AvgResolutionM))
		sg.f.SetRowHeight(name, r, 18)
	}

	return nil
}

// ─────────────────────────────────────────────────────────────────────────────
// Helpers
// ─────────────────────────────────────────────────────────────────────────────

func cell(col string, row int) string {
	return fmt.Sprintf("%s%d", col, row)
}

func colLetter(i int) string {
	return string(rune('A' + i))
}

func channelLabel(channels []string) string {
	if len(channels) == 0 {
		return "All Channels"
	}
	labels := make([]string, len(channels))
	for i, ch := range channels {
		labels[i] = labelChannel(ch)
	}
	return strings.Join(labels, ", ")
}

func labelChannel(ch string) string {
	m := map[string]string{
		"email":        "Email",
		"whatsapp":     "WhatsApp",
		"social_media": "Social Media",
		"live_chat":    "Live Chat",
		"call_center":  "Call Center",
	}
	if v, ok := m[ch]; ok {
		return v
	}
	return ch
}

func formatDuration(minutes float64) string {
	if minutes < 0 {
		return "-"
	}
	if minutes < 60 {
		return fmt.Sprintf("%.0f menit", minutes)
	}
	hours := minutes / 60
	if hours < 24 {
		return fmt.Sprintf("%.1f jam", hours)
	}
	days := hours / 24
	return fmt.Sprintf("%.1f hari", days)
}

func buildFilename(f entities.ExportFilter) string {
	chPart := "AllChannels"
	if len(f.Channels) > 0 {
		parts := make([]string, len(f.Channels))
		for i, ch := range f.Channels {
			parts[i] = labelChannel(ch)
		}
		chPart = strings.ReplaceAll(strings.Join(parts, "-"), " ", "")
	}
	return fmt.Sprintf("Dashboard_Report_%s_%s_%s.xlsx", f.From, f.To, chPart)
}