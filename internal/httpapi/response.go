package httpapi

import (
	"encoding/json"
	"net/http"
	"strings"

	"guardpanel.local/guardpanel/internal/domain"
	"guardpanel.local/guardpanel/internal/report"
)

type RecordView struct {
	ID          string   `json:"id"`
	MachineID   string   `json:"machine_id"`
	Title       string   `json:"title"`
	Status      string   `json:"status"`
	StatusLabel string   `json:"status_label"`
	Owner       string   `json:"owner"`
	Tags        []string `json:"tags"`
}

func NewRecordView(record domain.Record) RecordView {
	return RecordView{
		ID:          record.ID,
		MachineID:   record.MachineID,
		Title:       record.Title,
		Status:      record.Status,
		StatusLabel: domain.StatusLabel(record.Status),
		Owner:       record.Owner,
		Tags:        append([]string(nil), record.Tags...),
	}
}

func NewRecordViews(records []domain.Record) []RecordView {
	views := make([]RecordView, 0, len(records))
	for _, record := range records {
		views = append(views, NewRecordView(record))
	}
	return views
}

func writeRecordList(w http.ResponseWriter, records []domain.Record) {
	views := NewRecordViews(records)
	writeJSON(w, http.StatusOK, views)
}

func writeReport(w http.ResponseWriter, result report.RecordReport) {
	writeJSON(w, http.StatusOK, map[string]any{
		"total":     result.Total,
		"by_status": result.ByStatus,
		"lines":     result.Lines,
	})
}

func parseAction(path string) string {
	parts := strings.Split(strings.Trim(path, "/"), "/")
	if len(parts) < 3 {
		return ""
	}
	return parts[2]
}

func decodeMap(r *http.Request) (map[string]string, error) {
	var input map[string]string
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		return nil, err
	}
	return input, nil
}
