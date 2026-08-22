package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"

	"guardpanel.local/guardpanel/internal/config"
	"guardpanel.local/guardpanel/internal/domain"
	"guardpanel.local/guardpanel/internal/httpapi"
	"guardpanel.local/guardpanel/internal/notifier"
	"guardpanel.local/guardpanel/internal/service"
	"guardpanel.local/guardpanel/internal/store"
)

func main() {
	cfg, err := config.Parse(os.Args[1:])
	if err != nil {
		log.Fatal(err)
	}
	st, err := store.Open(cfg.DBPath)
	if err != nil {
		log.Fatal(err)
	}
	defer st.Close()
	notice, err := notifier.New(notifier.SenderFunc(func(context.Context, notifier.Request) error { return nil }))
	if err != nil {
		log.Fatal(err)
	}
	svc, err := service.New(st, notice, domain.FixedClock{Value: 1})
	if err != nil {
		log.Fatal(err)
	}
	if cfg.IsServe() {
		server := httpapi.New(svc)
		log.Printf("guardpanel listening on %s", cfg.Address)
		log.Fatal(http.ListenAndServe(cfg.Address, server.Handler()))
	}
	fmt.Println("无人售货机本地守护面板")
	fmt.Println(cfg.String())
}
