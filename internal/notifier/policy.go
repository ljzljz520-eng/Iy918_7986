package notifier

import "context"

type Policy struct {
	MaxAttempts int
	Retryable   map[error]bool
}

func DefaultPolicy() Policy {
	return Policy{MaxAttempts: 4, Retryable: map[error]bool{ErrAttemptsExceeded: true}}
}

func (p Policy) Attempts() int {
	if p.MaxAttempts < 1 {
		return 1
	}
	return p.MaxAttempts
}

func (p Policy) ShouldRetry(ctx context.Context, err error, attempt int) bool {
	if ctx.Err() != nil || err == nil {
		return false
	}
	if attempt >= p.Attempts() {
		return false
	}
	return p.Retryable[err]
}

func (p Policy) Validate() bool { return p.MaxAttempts > 0 && p.Retryable != nil }
