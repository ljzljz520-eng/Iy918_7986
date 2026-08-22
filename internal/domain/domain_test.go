package domain

import "testing"

func TestRecordTransitions(t *testing.T) {
	record, err := NewRecord("r1", "vm1", "冷柜", "alice", []string{" Cold ", "cold"})
	if err != nil {
		t.Fatal(err)
	}
	if len(record.Tags) != 1 || record.Tags[0] != "cold" {
		t.Fatalf("tags=%v", record.Tags)
	}
	if err = record.SubmitReview(); err != nil {
		t.Fatal(err)
	}
	if err = record.Approve(); err != nil {
		t.Fatal(err)
	}
	if err = record.Archive(); err != nil {
		t.Fatal(err)
	}
	if record.Status != StatusArchived {
		t.Fatalf("status=%s", record.Status)
	}
}

func TestAttachmentAndWorkflowValidation(t *testing.T) {
	attachment, err := NewAttachment("a1", "r1", "manual.pdf", "hash", "manual.pdf", 1)
	if err != nil || !attachment.IsSafeName() || attachment.Extension() != "pdf" {
		t.Fatal(err)
	}
	workflow, err := NewWorkflow("w1", "r1", WorkflowReview, 10)
	if err != nil || !workflow.IsDue(10) {
		t.Fatalf("workflow=%+v err=%v", workflow, err)
	}
}
