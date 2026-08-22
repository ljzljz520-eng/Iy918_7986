package report

import (
	"encoding/json"
	"errors"
	"strings"

	"guardpanel.local/guardpanel/internal/domain"
)

type Export struct {
	Records   []domain.Record     `json:"records"`
	Events    []domain.AuditEvent `json:"events"`
	Dashboard Dashboard           `json:"dashboard"`
}

func NewExport(records []domain.Record, events []domain.AuditEvent) Export {
	return Export{Records: append([]domain.Record(nil), records...), Events: append([]domain.AuditEvent(nil), events...), Dashboard: BuildDashboard(records)}
}

func Encode(export Export) ([]byte, error) {
	return json.MarshalIndent(export, "", "  ")
}

func Decode(data []byte) (Export, error) {
	if len(strings.TrimSpace(string(data))) == 0 {
		return Export{}, errors.New("empty export")
	}
	var export Export
	if err := json.Unmarshal(data, &export); err != nil {
		return Export{}, err
	}
	return export, nil
}

func FilterByTag(records []domain.Record, tag string) []domain.Record {
	tag = strings.ToLower(strings.TrimSpace(tag))
	result := make([]domain.Record, 0)
	for _, record := range records {
		for _, candidate := range record.Tags {
			if candidate == tag {
				result = append(result, record)
				break
			}
		}
	}
	return result
}
