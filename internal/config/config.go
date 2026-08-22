package config

import (
	"errors"
	"flag"
	"fmt"
)

type Config struct {
	DBPath  string
	Address string
	Command string
}

func Default() Config {
	return Config{DBPath: "guardpanel.db", Address: "127.0.0.1:8080", Command: "help"}
}

func Parse(args []string) (Config, error) {
	cfg := Default()
	flags := flag.NewFlagSet("guardpanel", flag.ContinueOnError)
	flags.StringVar(&cfg.DBPath, "db", cfg.DBPath, "database path")
	flags.StringVar(&cfg.Address, "addr", cfg.Address, "listen address")
	if len(args) > 0 {
		cfg.Command = args[0]
		args = args[1:]
	}
	if err := flags.Parse(args); err != nil {
		return cfg, err
	}
	if cfg.DBPath == "" || cfg.Address == "" {
		return cfg, errors.New("db and address are required")
	}
	return cfg, nil
}

func (c Config) String() string {
	return fmt.Sprintf("command=%s db=%s addr=%s", c.Command, c.DBPath, c.Address)
}

func (c Config) IsServe() bool { return c.Command == "serve" }
