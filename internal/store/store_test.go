package store

import (
	"path/filepath"
	"testing"

	"guardpanel.local/guardpanel/internal/domain"
)

func TestPersistenceSurvivesReopen(t *testing.T) {
	path := filepath.Join(t.TempDir(), "guard.db")
	st, err := Open(path)
	if err != nil {
		t.Fatal(err)
	}
	record, _ := domain.NewRecord("r1", "vm1", "冷柜", "alice", []string{"cold"})
	event, _ := domain.NewAuditEvent("e1", "r1", "alice", "create", "ok", 1)
	workflow, _ := domain.NewWorkflow("w1", "r1", domain.WorkflowReview, 5)
	attachment, _ := domain.NewAttachment("a1", "r1", "manual.pdf", "hash", "manual.pdf", 1)
	for _, putErr := range []error{st.PutRecord(record), st.PutEvent(event), st.PutWorkflow(workflow), st.PutAttachment(attachment)} {
		if putErr != nil {
			t.Fatal(putErr)
		}
	}
	if err = st.Close(); err != nil {
		t.Fatal(err)
	}
	st, err = Open(path)
	if err != nil {
		t.Fatal(err)
	}
	defer st.Close()
	if _, err = st.GetRecord("r1"); err != nil {
		t.Fatal(err)
	}
	if _, err = st.GetEvent("e1"); err != nil {
		t.Fatal(err)
	}
	if _, err = st.GetWorkflow("w1"); err != nil {
		t.Fatal(err)
	}
	if _, err = st.GetAttachment("a1"); err != nil {
		t.Fatal(err)
	}
}

func TestStoreSnapshot(t *testing.T) {
	st, err := Open(filepath.Join(t.TempDir(), "snapshot.db"))
	if err != nil {
		t.Fatal(err)
	}
	defer st.Close()
	record, _ := domain.NewRecord("r1", "vm1", "title", "owner", nil)
	if err = st.PutRecord(record); err != nil {
		t.Fatal(err)
	}
	snapshot, err := st.Export()
	if err != nil || len(snapshot.Records) != 1 {
		t.Fatalf("snapshot=%+v err=%v", snapshot, err)
	}
	data, err := st.ExportJSON()
	if err != nil {
		t.Fatal(err)
	}
	decoded, err := DecodeSnapshot(data)
	if err != nil || len(decoded.Records) != 1 {
		t.Fatalf("decoded=%+v err=%v", decoded, err)
	}
}
