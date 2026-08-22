package service

import (
	"context"
	"path/filepath"
	"testing"

	"guardpanel.local/guardpanel/internal/domain"
	"guardpanel.local/guardpanel/internal/notifier"
	"guardpanel.local/guardpanel/internal/store"
)

func newTestService(t *testing.T) *Service {
	t.Helper()
	st, err := store.Open(filepath.Join(t.TempDir(), "service.db"))
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { st.Close() })
	notice, _ := notifier.New(notifier.SenderFunc(func(context.Context, notifier.Request) error { return nil }))
	svc, err := New(st, notice, domain.FixedClock{Value: 100})
	if err != nil {
		t.Fatal(err)
	}
	return svc
}

func TestWorkflowCreateReviewArchive(t *testing.T) {
	svc := newTestService(t)
	if _, err := svc.CreateRecord("r1", "vm1", "冷柜资料", "alice", []string{"cold"}); err != nil {
		t.Fatal(err)
	}
	if _, err := svc.SubmitReview("r1", 200); err != nil {
		t.Fatal(err)
	}
	if err := svc.Approve("r1"); err != nil {
		t.Fatal(err)
	}
	if err := svc.Archive("r1"); err != nil {
		t.Fatal(err)
	}
	records, err := svc.Search("vm1", domain.StatusArchived, "冷柜")
	if err != nil || len(records) != 1 {
		t.Fatalf("records=%v err=%v", records, err)
	}
}

func TestWorkflowSearchUpdatePublish(t *testing.T) {
	svc := newTestService(t)
	if _, err := svc.CreateRecord("r1", "vm1", "old", "alice", nil); err != nil {
		t.Fatal(err)
	}
	if _, err := svc.SubmitReview("r1", 200); err != nil {
		t.Fatal(err)
	}
	if err := svc.Approve("r1"); err != nil {
		t.Fatal(err)
	}
	if err := svc.UpdateTitle("r1", "new", "bob"); err != nil {
		t.Fatal(err)
	}
	if err := svc.Publish("r1"); err != nil {
		t.Fatal(err)
	}
	report, err := svc.QueryReport("vm1", domain.StatusApproved, "new")
	if err != nil || report.Total != 1 {
		t.Fatalf("report=%+v err=%v", report, err)
	}
}

func TestWorkflowImportReport(t *testing.T) {
	svc := newTestService(t)
	input := "id,machine,title,owner,tags,kind,name,deadline\nr1,vm1,冷柜,alice,cold|store,import,manual.pdf,200\n"
	result, err := svc.ImportCSV(input)
	if err != nil || result.Accepted != 1 {
		t.Fatalf("result=%+v err=%v", result, err)
	}
	lines, err := svc.AuditReport()
	if err != nil || len(lines) != 1 {
		t.Fatalf("lines=%v err=%v", lines, err)
	}
}
