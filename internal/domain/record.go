package domain

import (
	"errors"
	"sort"
	"strings"
)

const (
	StatusDraft    = "draft"
	StatusReview   = "review"
	StatusApproved = "approved"
	StatusArchived = "archived"
	StatusRejected = "rejected"
)

type Record struct {
	ID        string   `json:"id"`
	MachineID string   `json:"machine_id"`
	Title     string   `json:"title"`
	Status    string   `json:"status"`
	Owner     string   `json:"owner"`
	Tags      []string `json:"tags"`
	CreatedAt int64    `json:"created_at"`
	UpdatedAt int64    `json:"updated_at"`
}

func NewRecord(id, machineID, title, owner string, tags []string) (Record, error) {
	if strings.TrimSpace(id) == "" || strings.TrimSpace(machineID) == "" {
		return Record{}, errors.New("record id and machine id are required")
	}
	if strings.TrimSpace(title) == "" || strings.TrimSpace(owner) == "" {
		return Record{}, errors.New("title and owner are required")
	}
	return Record{ID: id, MachineID: machineID, Title: title, Owner: owner, Tags: NormalizeTags(tags), Status: StatusDraft}, nil
}

func NormalizeTags(tags []string) []string {
	seen := make(map[string]bool)
	out := make([]string, 0, len(tags))
	for _, tag := range tags {
		tag = strings.TrimSpace(strings.ToLower(tag))
		if tag != "" && !seen[tag] {
			seen[tag] = true
			out = append(out, tag)
		}
	}
	sort.Strings(out)
	return out
}

func (r Record) Validate() error {
	if r.ID == "" || r.MachineID == "" || r.Title == "" || r.Owner == "" {
		return errors.New("record has missing required fields")
	}
	if !ValidStatus(r.Status) {
		return errors.New("invalid record status")
	}
	return nil
}

func ValidStatus(status string) bool {
	switch status {
	case StatusDraft, StatusReview, StatusApproved, StatusArchived, StatusRejected:
		return true
	default:
		return false
	}
}

func (r *Record) SubmitReview() error {
	if r.Status != StatusDraft && r.Status != StatusRejected {
		return errors.New("record cannot enter review from current status")
	}
	r.Status = StatusReview
	return nil
}

func (r *Record) Approve() error {
	if r.Status != StatusReview {
		return errors.New("only review records can be approved")
	}
	r.Status = StatusApproved
	return nil
}

func (r *Record) Reject() error {
	if r.Status != StatusReview {
		return errors.New("only review records can be rejected")
	}
	r.Status = StatusRejected
	return nil
}

func (r *Record) Archive() error {
	if r.Status != StatusApproved {
		return errors.New("only approved records can be archived")
	}
	r.Status = StatusArchived
	return nil
}

func (r Record) Matches(machineID, status, query string) bool {
	if machineID != "" && r.MachineID != machineID {
		return false
	}
	if status != "" && r.Status != status {
		return false
	}
	if query != "" && !strings.Contains(strings.ToLower(r.Title), strings.ToLower(query)) {
		return false
	}
	return true
}
