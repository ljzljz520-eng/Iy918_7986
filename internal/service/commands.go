package service

import (
	"errors"
	"fmt"
	"strings"

	"guardpanel.local/guardpanel/internal/domain"
)

type Command struct {
	Name     string
	RecordID string
	Actor    string
	Value    string
	Deadline int64
}

type CommandResult struct {
	Record   domain.Record
	Workflow domain.Workflow
	Message  string
}

func ParseCommand(line string) (Command, error) {
	parts := strings.Fields(line)
	if len(parts) == 0 {
		return Command{}, errors.New("command is empty")
	}
	command := Command{Name: parts[0], Actor: "operator"}
	if len(parts) > 1 {
		command.RecordID = parts[1]
	}
	if len(parts) > 2 {
		command.Value = strings.Join(parts[2:], " ")
	}
	return command, nil
}

func (s *Service) Execute(command Command) (CommandResult, error) {
	switch command.Name {
	case "create":
		record, err := s.CreateRecord(command.RecordID, "machine", command.Value, command.Actor, nil)
		return CommandResult{Record: record, Message: "created"}, err
	case "review":
		workflow, err := s.SubmitReview(command.RecordID, command.Deadline)
		return CommandResult{Workflow: workflow, Message: "review started"}, err
	case "approve":
		err := s.Approve(command.RecordID)
		return CommandResult{Message: "approved"}, err
	case "archive":
		err := s.Archive(command.RecordID)
		return CommandResult{Message: "archived"}, err
	case "publish":
		err := s.Publish(command.RecordID)
		return CommandResult{Message: "published"}, err
	default:
		return CommandResult{}, fmt.Errorf("unknown command %s", command.Name)
	}
}

func (s *Service) ValidateReady(id string) error {
	record, err := s.store.RequireRecord(id)
	if err != nil {
		return err
	}
	if record.Status == domain.StatusArchived {
		return errors.New("record is already archived")
	}
	return nil
}

func (s *Service) ResolveWorkflow(recordID, kind string) (domain.Workflow, error) {
	workflows, err := s.store.WorkflowsForRecord(recordID)
	if err != nil {
		return domain.Workflow{}, err
	}
	for _, workflow := range workflows {
		if workflow.Kind == kind && workflow.State == domain.WorkflowPending {
			return workflow, nil
		}
	}
	return domain.Workflow{}, fmt.Errorf("pending %s workflow not found", kind)
}
