package service

import (
	"sort"

	"guardpanel.local/guardpanel/internal/domain"
)

type Metrics struct {
	Records     int
	Events      int
	Workflows   int
	Attachments int
	Statuses    map[string]int
}

func (s *Service) Metrics() (Metrics, error) {
	records, err := s.store.ListRecords()
	if err != nil {
		return Metrics{}, err
	}
	events, err := s.store.ListEvents()
	if err != nil {
		return Metrics{}, err
	}
	workflows, err := s.store.ListWorkflows()
	if err != nil {
		return Metrics{}, err
	}
	attachments, err := s.store.ListAttachments()
	if err != nil {
		return Metrics{}, err
	}
	statuses := map[string]int{}
	for _, record := range records {
		statuses[record.Status]++
	}
	return Metrics{Records: len(records), Events: len(events), Workflows: len(workflows), Attachments: len(attachments), Statuses: statuses}, nil
}

func SortEvents(events []domain.AuditEvent) []domain.AuditEvent {
	result := append([]domain.AuditEvent(nil), events...)
	sort.Slice(result, func(i, j int) bool {
		if result[i].CreatedAt == result[j].CreatedAt {
			return result[i].ID < result[j].ID
		}
		return result[i].CreatedAt < result[j].CreatedAt
	})
	return result
}
