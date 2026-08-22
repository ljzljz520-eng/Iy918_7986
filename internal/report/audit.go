package report

import (
	"fmt"
	"sort"

	"guardpanel.local/guardpanel/internal/domain"
)

func AuditLines(events []domain.AuditEvent) []string {
	copyEvents := append([]domain.AuditEvent(nil), events...)
	sort.Slice(copyEvents, func(i, j int) bool {
		if copyEvents[i].CreatedAt == copyEvents[j].CreatedAt {
			return copyEvents[i].ID < copyEvents[j].ID
		}
		return copyEvents[i].CreatedAt < copyEvents[j].CreatedAt
	})
	lines := make([]string, 0, len(copyEvents))
	for _, event := range copyEvents {
		lines = append(lines, fmt.Sprintf("%d %s %s %s", event.CreatedAt, event.RecordID, event.Action, event.Outcome))
	}
	return lines
}

func OutcomeCounts(events []domain.AuditEvent) map[string]int {
	result := map[string]int{}
	for _, event := range events {
		result[event.Outcome]++
	}
	return result
}
