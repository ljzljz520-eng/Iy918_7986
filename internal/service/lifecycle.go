package service

import (
	"errors"
	"fmt"

	"guardpanel.local/guardpanel/internal/domain"
)

type Lifecycle struct {
	States  []string
	Current string
}

func NewLifecycle() Lifecycle {
	return Lifecycle{States: []string{domain.StatusDraft, domain.StatusReview, domain.StatusApproved, domain.StatusArchived}, Current: domain.StatusDraft}
}

func (l *Lifecycle) Advance(next string) error {
	if !domain.ValidStatus(next) {
		return errors.New("invalid lifecycle state")
	}
	if l.Current == domain.StatusArchived {
		return errors.New("lifecycle already archived")
	}
	l.Current = next
	l.States = append(l.States, next)
	return nil
}

func (l Lifecycle) CanEdit() bool {
	return l.Current == domain.StatusDraft || l.Current == domain.StatusReview || l.Current == domain.StatusRejected
}
func (l Lifecycle) CanArchive() bool { return l.Current == domain.StatusApproved }
func (l Lifecycle) Summary() string  { return fmt.Sprintf("%s:%d", l.Current, len(l.States)) }

func ValidateTransition(from, to string) error {
	if from == to {
		return errors.New("transition must change state")
	}
	if to == domain.StatusReview && from != domain.StatusDraft && from != domain.StatusRejected {
		return errors.New("invalid review transition")
	}
	if to == domain.StatusApproved && from != domain.StatusReview {
		return errors.New("invalid approval transition")
	}
	if to == domain.StatusArchived && from != domain.StatusApproved {
		return errors.New("invalid archive transition")
	}
	return nil
}
