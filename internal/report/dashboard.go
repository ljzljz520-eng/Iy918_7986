package report

import (
	"fmt"
	"sort"
	"strings"

	"guardpanel.local/guardpanel/internal/domain"
)

type Dashboard struct {
	Total    int
	Draft    int
	Review   int
	Approved int
	Archived int
	Rejected int
	Owners   []string
}

func BuildDashboard(records []domain.Record) Dashboard {
	dashboard := Dashboard{Total: len(records)}
	owners := map[string]bool{}
	for _, record := range records {
		owners[record.Owner] = true
		switch record.Status {
		case domain.StatusDraft:
			dashboard.Draft++
		case domain.StatusReview:
			dashboard.Review++
		case domain.StatusApproved:
			dashboard.Approved++
		case domain.StatusArchived:
			dashboard.Archived++
		case domain.StatusRejected:
			dashboard.Rejected++
		}
	}
	for owner := range owners {
		dashboard.Owners = append(dashboard.Owners, owner)
	}
	sort.Strings(dashboard.Owners)
	return dashboard
}

func (d Dashboard) CompletionRate() float64 {
	if d.Total == 0 {
		return 0
	}
	return float64(d.Archived) / float64(d.Total)
}

func (d Dashboard) NeedsAttention() int {
	return d.Review + d.Rejected
}

func (d Dashboard) Render() string {
	lines := []string{
		fmt.Sprintf("total=%d", d.Total),
		fmt.Sprintf("draft=%d", d.Draft),
		fmt.Sprintf("review=%d", d.Review),
		fmt.Sprintf("approved=%d", d.Approved),
		fmt.Sprintf("archived=%d", d.Archived),
		fmt.Sprintf("rejected=%d", d.Rejected),
		fmt.Sprintf("owners=%s", strings.Join(d.Owners, ",")),
	}
	return strings.Join(lines, "\n")
}

func Timeline(events []domain.AuditEvent) []string {
	lines := AuditLines(events)
	if len(lines) == 0 {
		return []string{"no events"}
	}
	return lines
}

func StatusDistribution(records []domain.Record) []string {
	counts := BuildRecordReport(records).ByStatus
	keys := make([]string, 0, len(counts))
	for key := range counts {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	result := make([]string, 0, len(keys))
	for _, key := range keys {
		result = append(result, fmt.Sprintf("%s=%d", key, counts[key]))
	}
	return result
}
