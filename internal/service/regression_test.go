package service

import (
	"context"
	"errors"
	"path/filepath"
	"testing"

	"guardpanel.local/guardpanel/internal/domain"
	"guardpanel.local/guardpanel/internal/notifier"
	"guardpanel.local/guardpanel/internal/store"
)

func TestBusiness11Regression(t *testing.T) {
	st, err := store.Open(filepath.Join(t.TempDir(), "regression.db"))
	if err != nil {
		t.Fatal(err)
	}
	defer st.Close()
	var cancel context.CancelFunc
	notice, _ := notifier.New(notifier.SenderFunc(func(ctx context.Context, request notifier.Request) error {
		cancel()
		if ctx.Err() == nil {
			return nil
		}
		return errors.New("temporary channel failure")
	}))
	svc, err := New(st, notice, domain.FixedClock{Value: 10})
	if err != nil {
		t.Fatal(err)
	}
	if _, err = svc.CreateRecord("918-11", "vm918", "隐患升级", "operator", nil); err != nil {
		t.Fatal(err)
	}
	ctx, cancel := context.WithCancel(context.Background())
	result := svc.NotifyUpgrade(ctx, "918-11", "deadline notice", 10)
	if !notifier.IsTimeout(result) {
		t.Fatalf("expected explicit timeout, got %+v", result)
	}
}
