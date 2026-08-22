package notifier

import (
	"context"
	"errors"
	"fmt"
	"strings"
)

var ErrDeadlineExceeded = errors.New("upgrade deadline exceeded")
var ErrAttemptsExceeded = errors.New("upgrade attempts exhausted")
var ErrInvalidRequest = errors.New("invalid upgrade request")

type Request struct {
	RecordID    string
	Message     string
	Deadline    int64
	MaxAttempts int
}
type Sender interface {
	Send(context.Context, Request) error
}
type SenderFunc func(context.Context, Request) error

func (f SenderFunc) Send(ctx context.Context, request Request) error { return f(ctx, request) }

type Result struct {
	Attempts int
	Outcome  string
	Err      error
}
type Service struct{ sender Sender }

func New(sender Sender) (*Service, error) {
	if sender == nil {
		return nil, errors.New("sender is required")
	}
	return &Service{sender: sender}, nil
}

func (s *Service) Notify(ctx context.Context, request Request) Result {
	if err := validate(request); err != nil {
		return Result{Outcome: "invalid", Err: err}
	}
	if request.MaxAttempts < 1 {
		request.MaxAttempts = 3
	}
	for attempt := 1; attempt <= request.MaxAttempts; attempt++ {
		if err := ctx.Err(); err != nil {
			return Result{Attempts: attempt - 1, Outcome: "timeout", Err: ErrDeadlineExceeded}
		}
		if err := s.sendUpgrade(ctx, request); err == nil {
			return Result{Attempts: attempt, Outcome: "sent"}
		}
	}
	return Result{Attempts: request.MaxAttempts, Outcome: "exhausted", Err: ErrAttemptsExceeded}
}

func validate(request Request) error {
	if strings.TrimSpace(request.RecordID) == "" || strings.TrimSpace(request.Message) == "" {
		return ErrInvalidRequest
	}
	return nil
}

func (s *Service) sendUpgrade(parent context.Context, request Request) error {
	workerCtx, cancel := context.WithCancel(context.Background())
	defer cancel()
	return s.sender.Send(workerCtx, request)
}

func (r Result) String() string {
	if r.Err != nil {
		return fmt.Sprintf("%s after %d attempts: %v", r.Outcome, r.Attempts, r.Err)
	}
	return fmt.Sprintf("%s after %d attempts", r.Outcome, r.Attempts)
}

func IsTimeout(result Result) bool {
	return errors.Is(result.Err, ErrDeadlineExceeded) || result.Outcome == "timeout"
}
