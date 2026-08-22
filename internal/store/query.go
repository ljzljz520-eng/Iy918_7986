package store

import (
	"errors"
	"strings"

	"guardpanel.local/guardpanel/internal/domain"
)

type Query struct {
	MachineID string
	Status    string
	Owner     string
	Text      string
	Limit     int
}

func (q Query) normalized() Query {
	q.MachineID = strings.TrimSpace(q.MachineID)
	q.Status = strings.TrimSpace(q.Status)
	q.Owner = strings.TrimSpace(q.Owner)
	q.Text = strings.TrimSpace(q.Text)
	if q.Limit < 0 {
		q.Limit = 0
	}
	return q
}

func (s *Store) QueryRecords(query Query) ([]domain.Record, error) {
	query = query.normalized()
	records, err := s.ListRecords()
	if err != nil {
		return nil, err
	}
	filtered := make([]domain.Record, 0, len(records))
	for _, record := range records {
		if query.MachineID != "" && record.MachineID != query.MachineID {
			continue
		}
		if query.Status != "" && record.Status != query.Status {
			continue
		}
		if query.Owner != "" && record.Owner != query.Owner {
			continue
		}
		if query.Text != "" && !record.Matches("", "", query.Text) {
			continue
		}
		filtered = append(filtered, record)
		if query.Limit > 0 && len(filtered) >= query.Limit {
			break
		}
	}
	return filtered, nil
}

func (s *Store) RequireRecord(id string) (domain.Record, error) {
	record, err := s.GetRecord(id)
	if err != nil {
		return domain.Record{}, err
	}
	if err = record.Validate(); err != nil {
		return domain.Record{}, err
	}
	return record, nil
}

func (s *Store) RequireWorkflow(id string) (domain.Workflow, error) {
	workflow, err := s.GetWorkflow(id)
	if err != nil {
		return domain.Workflow{}, err
	}
	if err = workflow.Validate(); err != nil {
		return domain.Workflow{}, err
	}
	return workflow, nil
}

func (s *Store) RecordsForMachine(machineID string) ([]domain.Record, error) {
	if strings.TrimSpace(machineID) == "" {
		return nil, errors.New("machine id is required")
	}
	return s.QueryRecords(Query{MachineID: machineID})
}

func (s *Store) EventsForRecord(recordID string) ([]domain.AuditEvent, error) {
	events, err := s.ListEvents()
	if err != nil {
		return nil, err
	}
	filtered := make([]domain.AuditEvent, 0)
	for _, event := range events {
		if event.RecordID == recordID {
			filtered = append(filtered, event)
		}
	}
	return filtered, nil
}

func (s *Store) WorkflowsForRecord(recordID string) ([]domain.Workflow, error) {
	workflows, err := s.ListWorkflows()
	if err != nil {
		return nil, err
	}
	result := make([]domain.Workflow, 0)
	for _, workflow := range workflows {
		if workflow.RecordID == recordID {
			result = append(result, workflow)
		}
	}
	return result, nil
}

func (s *Store) AttachmentsForRecord(recordID string) ([]domain.Attachment, error) {
	attachments, err := s.ListAttachments()
	if err != nil {
		return nil, err
	}
	result := make([]domain.Attachment, 0)
	for _, attachment := range attachments {
		if attachment.RecordID == recordID {
			result = append(result, attachment)
		}
	}
	return result, nil
}
