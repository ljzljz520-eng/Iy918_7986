package config

import (
	"errors"
	"net"
	"path/filepath"
	"strings"
)

func (c Config) Validate() error {
	if strings.TrimSpace(c.DBPath) == "" {
		return errors.New("database path is empty")
	}
	if filepath.Base(c.DBPath) == "." || filepath.Base(c.DBPath) == ".." {
		return errors.New("database path must name a file")
	}
	if _, _, err := net.SplitHostPort(c.Address); err != nil {
		return errors.New("address must include host and port")
	}
	return nil
}

func (c Config) DatabaseFile() string {
	return filepath.Clean(c.DBPath)
}

func (c Config) IsHelp() bool {
	return c.Command == "help" || c.Command == ""
}

func SupportedCommands() []string {
	return []string{"help", "serve"}
}
