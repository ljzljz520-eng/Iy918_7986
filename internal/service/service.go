package service

import (
	"context"
	"errors"
	"fmt"
	"sort"

	"guardpanel.local/guardpanel/internal/domain"
	"guardpanel.local/guardpanel/internal/importer"
	"guardpanel.local/guardpanel/internal/notifier"
	"guardpanel.local/guardpanel/internal/report"
	"guardpanel.local/guardpanel/internal/store"
)

type Service struct {
	store    *store.Store
	notifier *notifier.Service
	clock    domain.Clock
}

func New(st *store.Store, notice *notifier.Service, clock domain.Clock) (*Service, error) {
	if st == nil || notice == nil || clock == nil {
		return nil, errors.New("service dependencies are required")
	}
	return &Service{store: st, notifier: notice, clock: clock}, nil
}

func (s *Service) CreateRecord(id, machineID, title, owner string, tags []string) (domain.Record, error) {
	record, err := domain.NewRecord(id, machineID, title, owner, tags)
	if err != nil {
		return domain.Record{}, err
	}
	record.CreatedAt = s.clock.Now()
	record.UpdatedAt = record.CreatedAt
	if err := s.store.PutRecord(record); err != nil {
		return domain.Record{}, err
	}
	return record, s.recordEvent(record.ID, "system", "create", "ok")
}

func (s *Service) SubmitReview(id string, deadline int64) (domain.Workflow, error) {
	record, err := s.store.GetRecord(id)
	if err != nil {
		return domain.Workflow{}, err
	}
	if err = record.SubmitReview(); err != nil {
		return domain.Workflow{}, err
	}
	record.UpdatedAt = s.clock.Now()
	if err = s.store.PutRecord(record); err != nil {
		return domain.Workflow{}, err
	}
	workflow, err := domain.NewWorkflow(id+"-review", id, domain.WorkflowReview, deadline)
	if err != nil {
		return domain.Workflow{}, err
	}
	workflow.UpdatedAt = s.clock.Now()
	if err = s.store.PutWorkflow(workflow); err != nil {
		return domain.Workflow{}, err
	}
	return workflow, s.recordEvent(id, "operator", "submit_review", "ok")
}

func (s *Service) NotifyUpgrade(ctx context.Context, id, message string, deadline int64) notifier.Result {
	result := s.notifier.Notify(ctx, notifier.Request{RecordID: id, Message: message, Deadline: deadline, MaxAttempts: 5})
	outcome := result.Outcome
	if result.Err != nil {
		outcome = "timeout"
	}
	_ = s.recordEvent(id, "notifier", "upgrade", outcome)
	return result
}

func (s *Service) Approve(id string) error {
	record, err := s.store.GetRecord(id)
	if err != nil {
		return err
	}
	if err = record.Approve(); err != nil {
		return err
	}
	record.UpdatedAt = s.clock.Now()
	if err = s.store.PutRecord(record); err != nil {
		return err
	}
	return s.recordEvent(id, "reviewer", "approve", "approved")
}

func (s *Service) Reject(id string) error {
	record, err := s.store.GetRecord(id)
	if err != nil {
		return err
	}
	if err = record.Reject(); err != nil {
		return err
	}
	record.UpdatedAt = s.clock.Now()
	if err = s.store.PutRecord(record); err != nil {
		return err
	}
	return s.recordEvent(id, "reviewer", "reject", "rejected")
}

func (s *Service) Archive(id string) error {
	record, err := s.store.GetRecord(id)
	if err != nil {
		return err
	}
	if err = record.Archive(); err != nil {
		return err
	}
	record.UpdatedAt = s.clock.Now()
	if err = s.store.PutRecord(record); err != nil {
		return err
	}
	workflow, err := s.store.GetWorkflow(id + "-review")
	if err == nil && workflow.State == domain.WorkflowPending {
		_ = workflow.Complete()
		workflow.UpdatedAt = s.clock.Now()
		_ = s.store.PutWorkflow(workflow)
	}
	return s.recordEvent(id, "operator", "archive", "ok")
}

func (s *Service) UpdateTitle(id, title, actor string) error {
	if title == "" {
		return errors.New("title is required")
	}
	record, err := s.store.GetRecord(id)
	if err != nil {
		return err
	}
	record.Title = title
	record.UpdatedAt = s.clock.Now()
	if err = s.store.PutRecord(record); err != nil {
		return err
	}
	return s.recordEvent(id, actor, "update_title", "ok")
}

func (s *Service) Search(machineID, status, query string) ([]domain.Record, error) {
	records, err := s.store.ListRecords()
	if err != nil {
		return nil, err
	}
	return report.Filter(records, machineID, status, query), nil
}

func (s *Service) Publish(id string) error {
	record, err := s.store.GetRecord(id)
	if err != nil {
		return err
	}
	if record.Status != domain.StatusApproved && record.Status != domain.StatusArchived {
		return fmt.Errorf("record %s is not publishable", id)
	}
	return s.recordEvent(id, "publisher", "publish", "ok")
}

func (s *Service) ImportCSV(input string) (importer.Report, error) {
	batch, err := importer.ParseCSV(input)
	if err != nil {
		return importer.Report{}, err
	}
	if err = importer.ValidateBatch(batch); err != nil {
		return importer.Report{}, err
	}
	for _, row := range batch.Rows {
		row.Record.CreatedAt = s.clock.Now()
		row.Record.UpdatedAt = row.Record.CreatedAt
		if err = s.store.PutRecord(row.Record); err != nil {
			return importer.Report{}, err
		}
		if err = s.store.PutWorkflow(row.Workflow); err != nil {
			return importer.Report{}, err
		}
		if err = s.store.PutAttachment(row.Attachment); err != nil {
			return importer.Report{}, err
		}
		if err = s.recordEvent(row.Record.ID, "import", "import", "ok"); err != nil {
			return importer.Report{}, err
		}
	}
	return importer.BuildReport(batch), nil
}

func (s *Service) QueryReport(machineID, status, query string) (report.RecordReport, error) {
	records, err := s.Search(machineID, status, query)
	if err != nil {
		return report.RecordReport{}, err
	}
	return report.BuildRecordReport(records), nil
}

func (s *Service) AuditReport() ([]string, error) {
	events, err := s.store.ListEvents()
	if err != nil {
		return nil, err
	}
	return report.AuditLines(events), nil
}

func (s *Service) recordEvent(recordID, actor, action, outcome string) error {
	id := fmt.Sprintf("%s-%d-%s", recordID, s.clock.Now(), action)
	event, err := domain.NewAuditEvent(id, recordID, actor, action, outcome, s.clock.Now())
	if err != nil {
		return err
	}
	return s.store.PutEvent(event)
}

func (s *Service) Snapshot() (store.Snapshot, error) { return s.store.Export() }

func (s *Service) Restore(snapshot store.Snapshot) error { return s.store.Replace(snapshot) }

func SortRecords(records []domain.Record) []domain.Record {
	copyRecords := append([]domain.Record(nil), records...)
	sort.Slice(copyRecords, func(i, j int) bool { return copyRecords[i].ID < copyRecords[j].ID })
	return copyRecords
}
