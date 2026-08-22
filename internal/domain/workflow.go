package domain

import "errors"

type Workflow struct {
	ID        string `json:"id"`
	RecordID  string `json:"record_id"`
	Kind      string `json:"kind"`
	State     string `json:"state"`
	Deadline  int64  `json:"deadline"`
	UpdatedAt int64  `json:"updated_at"`
}

const (
	WorkflowReview  = "review"
	WorkflowImport  = "import"
	WorkflowArchive = "archive"
	WorkflowPending = "pending"
	WorkflowDone    = "done"
	WorkflowExpired = "expired"
)

func NewWorkflow(id, recordID, kind string, deadline int64) (Workflow, error) {
	if id == "" || recordID == "" || kind == "" {
		return Workflow{}, errors.New("workflow identity is required")
	}
	return Workflow{ID: id, RecordID: recordID, Kind: kind, State: WorkflowPending, Deadline: deadline}, nil
}

func (w *Workflow) Complete() error {
	if w.State != WorkflowPending {
		return errors.New("workflow is not pending")
	}
	w.State = WorkflowDone
	return nil
}

func (w *Workflow) Expire() error {
	if w.State == WorkflowDone {
		return errors.New("completed workflow cannot expire")
	}
	w.State = WorkflowExpired
	return nil
}

func (w Workflow) IsDue(now int64) bool {
	return w.Deadline > 0 && now >= w.Deadline && w.State == WorkflowPending
}

func (w Workflow) Validate() error {
	if w.ID == "" || w.RecordID == "" || w.Kind == "" {
		return errors.New("workflow fields are required")
	}
	if w.State != WorkflowPending && w.State != WorkflowDone && w.State != WorkflowExpired {
		return errors.New("invalid workflow state")
	}
	return nil
}
