package store

import (
	"encoding/json"
	"errors"

	"go.etcd.io/bbolt"
	"guardpanel.local/guardpanel/internal/domain"
)

type Transaction struct {
	store       *Store
	records     []domain.Record
	events      []domain.AuditEvent
	workflows   []domain.Workflow
	attachments []domain.Attachment
}

func (s *Store) Begin() (*Transaction, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	if s.db == nil {
		return nil, errors.New("store is closed")
	}
	return &Transaction{store: s}, nil
}

func (t *Transaction) AddRecord(record domain.Record) error {
	if err := record.Validate(); err != nil {
		return err
	}
	t.records = append(t.records, record)
	return nil
}

func (t *Transaction) AddEvent(event domain.AuditEvent) error {
	if err := event.Validate(); err != nil {
		return err
	}
	t.events = append(t.events, event)
	return nil
}

func (t *Transaction) AddWorkflow(workflow domain.Workflow) error {
	if err := workflow.Validate(); err != nil {
		return err
	}
	t.workflows = append(t.workflows, workflow)
	return nil
}

func (t *Transaction) AddAttachment(attachment domain.Attachment) error {
	if err := attachment.Validate(); err != nil {
		return err
	}
	t.attachments = append(t.attachments, attachment)
	return nil
}

func (t *Transaction) Commit() error {
	if t.store == nil {
		return errors.New("transaction has no store")
	}
	t.store.mu.RLock()
	defer t.store.mu.RUnlock()
	if t.store.db == nil {
		return errors.New("store is closed")
	}
	return t.store.db.Update(func(tx *bbolt.Tx) error {
		for _, item := range t.records {
			data, err := json.Marshal(item)
			if err != nil {
				return err
			}
			if err = tx.Bucket(bucketRecords).Put([]byte(item.ID), data); err != nil {
				return err
			}
		}
		for _, item := range t.events {
			data, err := json.Marshal(item)
			if err != nil {
				return err
			}
			if err = tx.Bucket(bucketEvents).Put([]byte(item.ID), data); err != nil {
				return err
			}
		}
		for _, item := range t.workflows {
			data, err := json.Marshal(item)
			if err != nil {
				return err
			}
			if err = tx.Bucket(bucketWorkflows).Put([]byte(item.ID), data); err != nil {
				return err
			}
		}
		for _, item := range t.attachments {
			data, err := json.Marshal(item)
			if err != nil {
				return err
			}
			if err = tx.Bucket(bucketAttachments).Put([]byte(item.ID), data); err != nil {
				return err
			}
		}
		return nil
	})
}

func (t *Transaction) Size() int {
	return len(t.records) + len(t.events) + len(t.workflows) + len(t.attachments)
}
