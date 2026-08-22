package service

import (
	"errors"
	"fmt"

	"guardpanel.local/guardpanel/internal/domain"
	"guardpanel.local/guardpanel/internal/importer"
	"guardpanel.local/guardpanel/internal/store"
)

type ImportSummary struct {
	Report      importer.Report
	Records     int
	Workflows   int
	Attachments int
}

func (s *Service) ImportSnapshot(snapshot store.Snapshot) (ImportSummary, error) {
	if len(snapshot.Records) == 0 {
		return ImportSummary{}, errors.New("snapshot has no records")
	}
	for index := range snapshot.Records {
		record := domain.NormalizeRecord(snapshot.Records[index])
		if err := domain.RecordValidationError(record); err != nil {
			return ImportSummary{}, fmt.Errorf("record %s: %w", record.ID, err)
		}
		snapshot.Records[index] = record
	}
	if err := s.store.Replace(snapshot); err != nil {
		return ImportSummary{}, err
	}
	return ImportSummary{Records: len(snapshot.Records), Workflows: len(snapshot.Workflows), Attachments: len(snapshot.Attachments), Report: importer.Report{Accepted: len(snapshot.Records), Message: "snapshot imported"}}, nil
}

func (s *Service) ValidateSnapshot(snapshot store.Snapshot) error {
	for _, record := range snapshot.Records {
		if err := record.Validate(); err != nil {
			return err
		}
	}
	for _, workflow := range snapshot.Workflows {
		if err := workflow.Validate(); err != nil {
			return err
		}
	}
	for _, attachment := range snapshot.Attachments {
		if err := attachment.Validate(); err != nil {
			return err
		}
	}
	return nil
}

func (s *Service) ExportJSON() ([]byte, error) {
	return s.store.ExportJSON()
}
