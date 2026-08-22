package httpapi

import (
	"context"
	"net/http/httptest"
	"path/filepath"
	"testing"

	"guardpanel.local/guardpanel/internal/domain"
	"guardpanel.local/guardpanel/internal/notifier"
	"guardpanel.local/guardpanel/internal/service"
	"guardpanel.local/guardpanel/internal/store"
)

func TestHTTPHealth(t *testing.T) {
	st, err := store.Open(filepath.Join(t.TempDir(), "http.db"))
	if err != nil {
		t.Fatal(err)
	}
	defer st.Close()
	notice, _ := notifier.New(notifier.SenderFunc(func(context.Context, notifier.Request) error { return nil }))
	svc, _ := service.New(st, notice, domain.FixedClock{Value: 1})
	server := New(svc)
	request := httptest.NewRequest("GET", "/health", nil)
	response := httptest.NewRecorder()
	server.Handler().ServeHTTP(response, request)
	if response.Code != 200 {
		t.Fatalf("code=%d", response.Code)
	}
}
