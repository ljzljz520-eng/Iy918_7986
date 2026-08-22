package store

import (
	"encoding/json"
	"errors"
	"sort"

	"guardpanel.local/guardpanel/internal/domain"
)

type Snapshot struct {
	Records     []domain.Record     `json:"records"`
	Events      []domain.AuditEvent `json:"events"`
	Workflows   []domain.Workflow   `json:"workflows"`
	Attachments []domain.Attachment `json:"attachments"`
}

func (s *Store) Export() (Snapshot, error) {
	records, err := s.ListRecords()
	if err != nil {
		return Snapshot{}, err
	}
	events, err := s.ListEvents()
	if err != nil {
		return Snapshot{}, err
	}
	workflows, err := s.ListWorkflows()
	if err != nil {
		return Snapshot{}, err
	}
	attachments, err := s.ListAttachments()
	if err != nil {
		return Snapshot{}, err
	}
	return Snapshot{Records: records, Events: events, Workflows: workflows, Attachments: attachments}, nil
}

func (s *Store) Import(snapshot Snapshot) error {
	for _, record := range snapshot.Records {
		if err := s.PutRecord(record); err != nil {
			return err
		}
	}
	for _, event := range snapshot.Events {
		if err := s.PutEvent(event); err != nil {
			return err
		}
	}
	for _, workflow := range snapshot.Workflows {
		if err := s.PutWorkflow(workflow); err != nil {
			return err
		}
	}
	for _, attachment := range snapshot.Attachments {
		if err := s.PutAttachment(attachment); err != nil {
			return err
		}
	}
	return nil
}

func (s *Store) ExportJSON() ([]byte, error) {
	snapshot, err := s.Export()
	if err != nil {
		return nil, err
	}
	return json.MarshalIndent(snapshot, "", "  ")
}

func DecodeSnapshot(data []byte) (Snapshot, error) {
	var snapshot Snapshot
	if len(data) == 0 {
		return snapshot, errors.New("empty snapshot")
	}
	err := json.Unmarshal(data, &snapshot)
	return snapshot, err
}

func (s *Store) Replace(snapshot Snapshot) error {
	if _, err := s.Count(); err != nil {
		return err
	}
	for _, record := range snapshot.Records {
		if err := record.Validate(); err != nil {
			return err
		}
	}
	if len(snapshot.Records) > 1 {
		sort.Slice(snapshot.Records, func(i, j int) bool { return snapshot.Records[i].ID < snapshot.Records[j].ID })
	}
	return s.Import(snapshot)
}
