package report

import (
	"fmt"
	"sort"
	"strings"

	"guardpanel.local/guardpanel/internal/domain"
)

type RecordReport struct {
	Total    int
	ByStatus map[string]int
	Lines    []string
}

func BuildRecordReport(records []domain.Record) RecordReport {
	result := RecordReport{Total: len(records), ByStatus: map[string]int{}}
	copyRecords := append([]domain.Record(nil), records...)
	sort.Slice(copyRecords, func(i, j int) bool { return copyRecords[i].ID < copyRecords[j].ID })
	for _, record := range copyRecords {
		result.ByStatus[record.Status]++
		result.Lines = append(result.Lines, fmt.Sprintf("%s|%s|%s", record.ID, record.MachineID, record.Status))
	}
	return result
}

func Render(report RecordReport) string { return strings.Join(report.Lines, "\n") }

func Filter(records []domain.Record, machineID, status, query string) []domain.Record {
	result := make([]domain.Record, 0)
	for _, record := range records {
		if record.Matches(machineID, status, query) {
			result = append(result, record)
		}
	}
	sort.Slice(result, func(i, j int) bool { return result[i].ID < result[j].ID })
	return result
}

func GroupByOwner(records []domain.Record) map[string][]domain.Record {
	result := map[string][]domain.Record{}
	for _, record := range records {
		result[record.Owner] = append(result[record.Owner], record)
	}
	for owner := range result {
		sort.Slice(result[owner], func(i, j int) bool { return result[owner][i].ID < result[owner][j].ID })
	}
	return result
}
