package domain

import (
	"errors"
	"strings"
)

type Attachment struct {
	ID        string `json:"id"`
	RecordID  string `json:"record_id"`
	Name      string `json:"name"`
	Checksum  string `json:"checksum"`
	Path      string `json:"path"`
	CreatedAt int64  `json:"created_at"`
}

func NewAttachment(id, recordID, name, checksum, path string, createdAt int64) (Attachment, error) {
	if id == "" || recordID == "" || strings.TrimSpace(name) == "" {
		return Attachment{}, errors.New("attachment identity and name are required")
	}
	return Attachment{ID: id, RecordID: recordID, Name: name, Checksum: checksum, Path: path, CreatedAt: createdAt}, nil
}

func (a Attachment) Validate() error {
	if a.ID == "" || a.RecordID == "" || a.Name == "" {
		return errors.New("invalid attachment")
	}
	return nil
}

func (a Attachment) Extension() string {
	idx := strings.LastIndex(a.Name, ".")
	if idx < 0 {
		return ""
	}
	return strings.ToLower(a.Name[idx+1:])
}

func (a Attachment) IsSafeName() bool {
	return !strings.Contains(a.Name, "..") && !strings.ContainsAny(a.Name, "/\\")
}
