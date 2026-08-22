package importer

import (
	"errors"
	"fmt"
	"sort"
	"strings"

	"guardpanel.local/guardpanel/internal/domain"
)

type NormalizedRow struct {
	Record     domain.Record
	Workflow   domain.Workflow
	Attachment domain.Attachment
	Warnings   []string
}

func NormalizeRow(row Row) (NormalizedRow, error) {
	record := domain.NormalizeRecord(row.Record)
	workflow := row.Workflow
	attachment := row.Attachment
	warnings := make([]string, 0)
	if record.Status == "" {
		record.Status = domain.StatusDraft
		warnings = append(warnings, "status defaulted to draft")
	}
	if attachment.Checksum == "" {
		attachment.Checksum = domain.AttachmentChecksum(attachment.Name, attachment.Path)
		warnings = append(warnings, "checksum generated")
	}
	if attachment.Path == "" {
		attachment.Path = "attachments/" + attachment.Name
		warnings = append(warnings, "path generated")
	}
	if err := domain.ValidateEntitySet(record, workflow, attachment); err != nil {
		return NormalizedRow{}, err
	}
	return NormalizedRow{Record: record, Workflow: workflow, Attachment: attachment, Warnings: warnings}, nil
}

func NormalizeBatch(batch Batch) ([]NormalizedRow, error) {
	if len(batch.Rows) == 0 {
		return nil, errors.New("no rows to normalize")
	}
	result := make([]NormalizedRow, 0, len(batch.Rows))
	for _, row := range batch.Rows {
		normalized, err := NormalizeRow(row)
		if err != nil {
			return nil, err
		}
		result = append(result, normalized)
	}
	return result, nil
}

func SortRows(rows []NormalizedRow) []NormalizedRow {
	result := append([]NormalizedRow(nil), rows...)
	sort.Slice(result, func(i, j int) bool {
		return result[i].Record.ID < result[j].Record.ID
	})
	return result
}

func WarningSummary(rows []NormalizedRow) string {
	counts := map[string]int{}
	for _, row := range rows {
		for _, warning := range row.Warnings {
			counts[warning]++
		}
	}
	keys := make([]string, 0, len(counts))
	for key := range counts {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	parts := make([]string, 0, len(keys))
	for _, key := range keys {
		parts = append(parts, fmt.Sprintf("%s=%d", key, counts[key]))
	}
	return strings.Join(parts, ",")
}
