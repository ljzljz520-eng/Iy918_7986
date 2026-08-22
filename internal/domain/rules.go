package domain

import (
	"errors"
	"fmt"
	"sort"
	"strings"
)

type ValidationIssue struct {
	Field   string
	Message string
}

func (i ValidationIssue) Error() string {
	return i.Field + ": " + i.Message
}

func ValidateRecordDetailed(record Record) []ValidationIssue {
	issues := make([]ValidationIssue, 0)
	if strings.TrimSpace(record.ID) == "" {
		issues = append(issues, ValidationIssue{Field: "id", Message: "required"})
	}
	if strings.TrimSpace(record.MachineID) == "" {
		issues = append(issues, ValidationIssue{Field: "machine_id", Message: "required"})
	}
	if strings.TrimSpace(record.Title) == "" {
		issues = append(issues, ValidationIssue{Field: "title", Message: "required"})
	}
	if len(record.Title) > 160 {
		issues = append(issues, ValidationIssue{Field: "title", Message: "must be 160 characters or fewer"})
	}
	if strings.TrimSpace(record.Owner) == "" {
		issues = append(issues, ValidationIssue{Field: "owner", Message: "required"})
	}
	if !ValidStatus(record.Status) {
		issues = append(issues, ValidationIssue{Field: "status", Message: "unknown value"})
	}
	if record.CreatedAt < 0 || record.UpdatedAt < 0 {
		issues = append(issues, ValidationIssue{Field: "timestamps", Message: "cannot be negative"})
	}
	if record.UpdatedAt < record.CreatedAt {
		issues = append(issues, ValidationIssue{Field: "updated_at", Message: "cannot precede created_at"})
	}
	return issues
}

func RecordValidationError(record Record) error {
	issues := ValidateRecordDetailed(record)
	if len(issues) == 0 {
		return nil
	}
	parts := make([]string, 0, len(issues))
	for _, issue := range issues {
		parts = append(parts, issue.Error())
	}
	return errors.New(strings.Join(parts, "; "))
}

func NormalizeRecord(record Record) Record {
	record.ID = strings.TrimSpace(record.ID)
	record.MachineID = strings.TrimSpace(record.MachineID)
	record.Title = strings.TrimSpace(record.Title)
	record.Owner = strings.TrimSpace(record.Owner)
	record.Tags = NormalizeTags(record.Tags)
	return record
}

func TransitionPath(from, to string) ([]string, error) {
	if from == to {
		return nil, errors.New("transition does not change state")
	}
	paths := map[string][]string{
		StatusDraft + ":" + StatusReview:      {StatusDraft, StatusReview},
		StatusReview + ":" + StatusApproved:   {StatusReview, StatusApproved},
		StatusReview + ":" + StatusRejected:   {StatusReview, StatusRejected},
		StatusRejected + ":" + StatusReview:   {StatusRejected, StatusReview},
		StatusApproved + ":" + StatusArchived: {StatusApproved, StatusArchived},
	}
	path, ok := paths[from+":"+to]
	if !ok {
		return nil, fmt.Errorf("transition %s to %s is not allowed", from, to)
	}
	return append([]string(nil), path...), nil
}

func AllowedNextStates(status string) []string {
	var values []string
	switch status {
	case StatusDraft:
		values = []string{StatusReview}
	case StatusReview:
		values = []string{StatusApproved, StatusRejected}
	case StatusRejected:
		values = []string{StatusReview}
	case StatusApproved:
		values = []string{StatusArchived}
	default:
		values = []string{}
	}
	sort.Strings(values)
	return values
}

func StatusLabel(status string) string {
	labels := map[string]string{
		StatusDraft:    "草稿",
		StatusReview:   "审核中",
		StatusApproved: "已批准",
		StatusArchived: "已归档",
		StatusRejected: "已驳回",
	}
	if label, ok := labels[status]; ok {
		return label
	}
	return "未知"
}

func SortTags(tags []string) []string {
	result := NormalizeTags(tags)
	sort.Strings(result)
	return result
}

func ValidateEntitySet(record Record, workflow Workflow, attachment Attachment) error {
	if err := RecordValidationError(record); err != nil {
		return err
	}
	if err := workflow.Validate(); err != nil {
		return err
	}
	if err := attachment.Validate(); err != nil {
		return err
	}
	if workflow.RecordID != record.ID || attachment.RecordID != record.ID {
		return errors.New("entity relationship mismatch")
	}
	return nil
}

func WorkflowLabel(state string) string {
	switch state {
	case WorkflowPending:
		return "待处理"
	case WorkflowDone:
		return "已完成"
	case WorkflowExpired:
		return "已过期"
	default:
		return "未知"
	}
}

func EventKey(event AuditEvent) string {
	return fmt.Sprintf("%d:%s:%s", event.CreatedAt, event.RecordID, event.Action)
}
