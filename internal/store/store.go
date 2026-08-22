package store

import (
	"encoding/json"
	"errors"
	"fmt"
	"path/filepath"
	"sort"
	"sync"

	"go.etcd.io/bbolt"
	"guardpanel.local/guardpanel/internal/domain"
)

var (
	bucketRecords     = []byte("records")
	bucketEvents      = []byte("audit_events")
	bucketWorkflows   = []byte("workflows")
	bucketAttachments = []byte("attachments")
)

type Store struct {
	db *bbolt.DB
	mu sync.RWMutex
}

func Open(path string) (*Store, error) {
	if path == "" {
		return nil, errors.New("store path is required")
	}
	db, err := bbolt.Open(filepath.Clean(path), 0600, &bbolt.Options{NoSync: true})
	if err != nil {
		return nil, err
	}
	s := &Store{db: db}
	if err = s.initialize(); err != nil {
		_ = db.Close()
		return nil, err
	}
	return s, nil
}

func (s *Store) initialize() error {
	return s.db.Update(func(tx *bbolt.Tx) error {
		for _, bucket := range [][]byte{bucketRecords, bucketEvents, bucketWorkflows, bucketAttachments} {
			if _, err := tx.CreateBucketIfNotExists(bucket); err != nil {
				return err
			}
		}
		return nil
	})
}

func (s *Store) Close() error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.db == nil {
		return nil
	}
	err := s.db.Close()
	s.db = nil
	return err
}

func (s *Store) PutRecord(record domain.Record) error {
	return put(s, bucketRecords, record.ID, record)
}
func (s *Store) PutEvent(event domain.AuditEvent) error { return put(s, bucketEvents, event.ID, event) }
func (s *Store) PutWorkflow(workflow domain.Workflow) error {
	return put(s, bucketWorkflows, workflow.ID, workflow)
}
func (s *Store) PutAttachment(attachment domain.Attachment) error {
	return put(s, bucketAttachments, attachment.ID, attachment)
}

func put[T any](s *Store, bucket []byte, key string, value T) error {
	s.mu.RLock()
	defer s.mu.RUnlock()
	if s.db == nil {
		return errors.New("store is closed")
	}
	if key == "" {
		return errors.New("empty key")
	}
	data, err := json.Marshal(value)
	if err != nil {
		return err
	}
	return s.db.Update(func(tx *bbolt.Tx) error { return tx.Bucket(bucket).Put([]byte(key), data) })
}

func get[T any](s *Store, bucket []byte, key string) (T, error) {
	var zero T
	s.mu.RLock()
	defer s.mu.RUnlock()
	if s.db == nil {
		return zero, errors.New("store is closed")
	}
	var value T
	err := s.db.View(func(tx *bbolt.Tx) error {
		data := tx.Bucket(bucket).Get([]byte(key))
		if data == nil {
			return fmt.Errorf("%s %q not found", string(bucket), key)
		}
		return json.Unmarshal(data, &value)
	})
	return value, err
}

func (s *Store) GetRecord(key string) (domain.Record, error) {
	return get[domain.Record](s, bucketRecords, key)
}
func (s *Store) GetEvent(key string) (domain.AuditEvent, error) {
	return get[domain.AuditEvent](s, bucketEvents, key)
}
func (s *Store) GetWorkflow(key string) (domain.Workflow, error) {
	return get[domain.Workflow](s, bucketWorkflows, key)
}
func (s *Store) GetAttachment(key string) (domain.Attachment, error) {
	return get[domain.Attachment](s, bucketAttachments, key)
}

func list[T any](s *Store, bucket []byte) ([]T, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	if s.db == nil {
		return nil, errors.New("store is closed")
	}
	values := make([]T, 0)
	err := s.db.View(func(tx *bbolt.Tx) error {
		return tx.Bucket(bucket).ForEach(func(_, data []byte) error {
			var value T
			if err := json.Unmarshal(data, &value); err != nil {
				return err
			}
			values = append(values, value)
			return nil
		})
	})
	return values, err
}

func (s *Store) ListRecords() ([]domain.Record, error) {
	values, err := list[domain.Record](s, bucketRecords)
	sort.Slice(values, func(i, j int) bool { return values[i].ID < values[j].ID })
	return values, err
}
func (s *Store) ListEvents() ([]domain.AuditEvent, error) {
	values, err := list[domain.AuditEvent](s, bucketEvents)
	sort.Slice(values, func(i, j int) bool { return values[i].ID < values[j].ID })
	return values, err
}
func (s *Store) ListWorkflows() ([]domain.Workflow, error) {
	values, err := list[domain.Workflow](s, bucketWorkflows)
	sort.Slice(values, func(i, j int) bool { return values[i].ID < values[j].ID })
	return values, err
}
func (s *Store) ListAttachments() ([]domain.Attachment, error) {
	values, err := list[domain.Attachment](s, bucketAttachments)
	sort.Slice(values, func(i, j int) bool { return values[i].ID < values[j].ID })
	return values, err
}

func (s *Store) DeleteRecord(id string) error     { return remove(s, bucketRecords, id) }
func (s *Store) DeleteEvent(id string) error      { return remove(s, bucketEvents, id) }
func (s *Store) DeleteWorkflow(id string) error   { return remove(s, bucketWorkflows, id) }
func (s *Store) DeleteAttachment(id string) error { return remove(s, bucketAttachments, id) }

func remove(s *Store, bucket []byte, key string) error {
	s.mu.RLock()
	defer s.mu.RUnlock()
	if s.db == nil {
		return errors.New("store is closed")
	}
	return s.db.Update(func(tx *bbolt.Tx) error { return tx.Bucket(bucket).Delete([]byte(key)) })
}

func (s *Store) Count() (map[string]int, error) {
	result := map[string]int{}
	records, err := s.ListRecords()
	if err != nil {
		return nil, err
	}
	result["records"] = len(records)
	events, err := s.ListEvents()
	if err != nil {
		return nil, err
	}
	result["events"] = len(events)
	workflows, err := s.ListWorkflows()
	if err != nil {
		return nil, err
	}
	result["workflows"] = len(workflows)
	attachments, err := s.ListAttachments()
	if err != nil {
		return nil, err
	}
	result["attachments"] = len(attachments)
	return result, nil
}
