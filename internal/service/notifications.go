package service

import (
	"context"
	"errors"
	"fmt"

	"guardpanel.local/guardpanel/internal/domain"
	"guardpanel.local/guardpanel/internal/notifier"
)

type NotificationSummary struct {
	RecordID string
	Attempts int
	Outcome  string
	Message  string
}

func (s *Service) NotifyWithDeadline(ctx context.Context, recordID, message string, deadline int64) (NotificationSummary, error) {
	if _, err := s.store.RequireRecord(recordID); err != nil {
		return NotificationSummary{}, err
	}
	result := s.NotifyUpgrade(ctx, recordID, message, deadline)
	if notifier.IsTimeout(result) {
		return NotificationSummary{RecordID: recordID, Attempts: result.Attempts, Outcome: result.Outcome, Message: result.String()}, notifier.ErrDeadlineExceeded
	}
	if result.Err != nil {
		return NotificationSummary{RecordID: recordID, Attempts: result.Attempts, Outcome: result.Outcome, Message: result.String()}, result.Err
	}
	return NotificationSummary{RecordID: recordID, Attempts: result.Attempts, Outcome: result.Outcome, Message: result.String()}, nil
}

func (s *Service) CompleteReview(recordID string) error {
	workflow, err := s.ResolveWorkflow(recordID, domain.WorkflowReview)
	if err != nil {
		return err
	}
	if err = workflow.Complete(); err != nil {
		return err
	}
	workflow.UpdatedAt = s.clock.Now()
	if err = s.store.PutWorkflow(workflow); err != nil {
		return err
	}
	return s.recordEvent(recordID, "reviewer", "complete_review", "ok")
}

func (s *Service) ExpireReview(recordID string) error {
	workflow, err := s.ResolveWorkflow(recordID, domain.WorkflowReview)
	if err != nil {
		return err
	}
	if err = workflow.Expire(); err != nil {
		return err
	}
	workflow.UpdatedAt = s.clock.Now()
	if err = s.store.PutWorkflow(workflow); err != nil {
		return err
	}
	return s.recordEvent(recordID, "system", "expire_review", "timeout")
}

func FormatDeadline(deadline int64) string {
	if deadline == 0 {
		return "无截止时间"
	}
	return fmt.Sprintf("deadline=%d", deadline)
}

func deadlineError(err error) bool {
	return errors.Is(err, notifier.ErrDeadlineExceeded)
}
