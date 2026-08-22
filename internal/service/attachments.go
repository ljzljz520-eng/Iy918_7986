package service

import (
	"errors"
	"strings"

	"guardpanel.local/guardpanel/internal/domain"
)

func (s *Service) AddAttachment(recordID, name, path string) (domain.Attachment, error) {
	if _, err := s.store.RequireRecord(recordID); err != nil {
		return domain.Attachment{}, err
	}
	if strings.TrimSpace(name) == "" {
		return domain.Attachment{}, errors.New("attachment name is required")
	}
	attachment, err := domain.NewAttachment(
		domain.StableID(recordID, name),
		recordID,
		name,
		domain.AttachmentChecksum(name, path),
		path,
		s.clock.Now(),
	)
	if err != nil {
		return domain.Attachment{}, err
	}
	if !attachment.IsSafeName() {
		return domain.Attachment{}, errors.New("unsafe attachment name")
	}
	if err = s.store.PutAttachment(attachment); err != nil {
		return domain.Attachment{}, err
	}
	if err = s.recordEvent(recordID, "operator", "attach", "ok"); err != nil {
		return domain.Attachment{}, err
	}
	return attachment, nil
}

func (s *Service) RemoveAttachment(recordID, attachmentID string) error {
	attachments, err := s.store.AttachmentsForRecord(recordID)
	if err != nil {
		return err
	}
	for _, attachment := range attachments {
		if attachment.ID == attachmentID {
			if err = s.store.DeleteAttachment(attachmentID); err != nil {
				return err
			}
			return s.recordEvent(recordID, "operator", "detach", "ok")
		}
	}
	return errors.New("attachment not found")
}

func (s *Service) ListAttachments(recordID string) ([]domain.Attachment, error) {
	return s.store.AttachmentsForRecord(recordID)
}

func (s *Service) RecordBundle(recordID string) (domain.Record, []domain.Workflow, []domain.Attachment, error) {
	record, err := s.store.RequireRecord(recordID)
	if err != nil {
		return domain.Record{}, nil, nil, err
	}
	workflows, err := s.store.WorkflowsForRecord(recordID)
	if err != nil {
		return domain.Record{}, nil, nil, err
	}
	attachments, err := s.store.AttachmentsForRecord(recordID)
	if err != nil {
		return domain.Record{}, nil, nil, err
	}
	return record, workflows, attachments, nil
}
