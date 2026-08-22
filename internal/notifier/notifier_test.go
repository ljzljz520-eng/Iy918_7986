package notifier

import (
	"context"
	"errors"
	"testing"
)

func TestNotifierSuccess(t *testing.T) {
	service, err := New(SenderFunc(func(context.Context, Request) error { return nil }))
	if err != nil {
		t.Fatal(err)
	}
	result := service.Notify(context.Background(), Request{RecordID: "r1", Message: "ok", MaxAttempts: 2})
	if result.Outcome != "sent" || result.Attempts != 1 {
		t.Fatalf("result=%+v", result)
	}
}

func TestNotifierInvalid(t *testing.T) {
	service, _ := New(SenderFunc(func(context.Context, Request) error { return errors.New("unused") }))
	result := service.Notify(context.Background(), Request{})
	if !errors.Is(result.Err, ErrInvalidRequest) {
		t.Fatalf("result=%+v", result)
	}
}
