package report

import (
	"guardpanel.local/guardpanel/internal/domain"
	"testing"
)

func TestReportsAreDeterministic(t *testing.T) {
	first, _ := domain.NewRecord("b", "vm", "B", "bob", nil)
	first.Status = domain.StatusApproved
	second, _ := domain.NewRecord("a", "vm", "A", "alice", nil)
	report := BuildRecordReport([]domain.Record{first, second})
	if report.Lines[0] != "a|vm|draft" || Render(report) == "" {
		t.Fatalf("report=%+v", report)
	}
	grouped := GroupByOwner([]domain.Record{first, second})
	if len(grouped) != 2 {
		t.Fatal("group mismatch")
	}
}
