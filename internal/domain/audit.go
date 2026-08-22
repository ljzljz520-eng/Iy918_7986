package domain

import "errors"

type AuditEvent struct {
	ID        string `json:"id"`
	RecordID  string `json:"record_id"`
	Actor     string `json:"actor"`
	Action    string `json:"action"`
	Outcome   string `json:"outcome"`
	CreatedAt int64  `json:"created_at"`
}

func NewAuditEvent(id, recordID, actor, action, outcome string, createdAt int64) (AuditEvent, error) {
	if id == "" || recordID == "" || actor == "" || action == "" {
		return AuditEvent{}, errors.New("audit event fields are required")
	}
	return AuditEvent{ID: id, RecordID: recordID, Actor: actor, Action: action, Outcome: outcome, CreatedAt: createdAt}, nil
}

func (e AuditEvent) Validate() error {
	if e.ID == "" || e.RecordID == "" || e.Actor == "" || e.Action == "" {
		return errors.New("invalid audit event")
	}
	return nil
}

func (e AuditEvent) IsSuccessful() bool { return e.Outcome == "ok" || e.Outcome == "approved" }

func AuditSummary(events []AuditEvent) map[string]int {
	result := map[string]int{"ok": 0, "failed": 0}
	for _, event := range events {
		if event.IsSuccessful() {
			result["ok"]++
		} else {
			result["failed"]++
		}
	}
	return result
}
