package importer

import (
	"encoding/csv"
	"errors"
	"fmt"
	"io"
	"strconv"
	"strings"

	"guardpanel.local/guardpanel/internal/domain"
)

type Row struct {
	Record     domain.Record
	Workflow   domain.Workflow
	Attachment domain.Attachment
}
type Batch struct {
	Rows     []Row
	Rejected []string
}

func ParseCSV(input string) (Batch, error) {
	reader := csv.NewReader(strings.NewReader(input))
	reader.TrimLeadingSpace = true
	header, err := reader.Read()
	if err != nil {
		return Batch{}, err
	}
	if len(header) < 8 {
		return Batch{}, errors.New("csv requires eight columns")
	}
	batch := Batch{}
	for line := 2; ; line++ {
		fields, readErr := reader.Read()
		if readErr == io.EOF {
			break
		}
		if readErr != nil {
			return Batch{}, readErr
		}
		row, err := parseRow(fields, line)
		if err != nil {
			batch.Rejected = append(batch.Rejected, err.Error())
			continue
		}
		batch.Rows = append(batch.Rows, row)
	}
	if len(batch.Rows) == 0 && len(batch.Rejected) > 0 {
		return batch, errors.New("all import rows rejected")
	}
	return batch, nil
}

func parseRow(fields []string, line int) (Row, error) {
	if len(fields) < 8 {
		return Row{}, fmt.Errorf("line %d: expected eight columns", line)
	}
	deadline, err := strconv.ParseInt(fields[7], 10, 64)
	if err != nil {
		return Row{}, fmt.Errorf("line %d: invalid deadline", line)
	}
	record, err := domain.NewRecord(fields[0], fields[1], fields[2], fields[3], strings.Split(fields[4], "|"))
	if err != nil {
		return Row{}, fmt.Errorf("line %d: %w", line, err)
	}
	workflow, err := domain.NewWorkflow(fields[0]+"-wf", fields[0], fields[5], deadline)
	if err != nil {
		return Row{}, fmt.Errorf("line %d: %w", line, err)
	}
	attachment, err := domain.NewAttachment(fields[0]+"-att", fields[0], fields[6], fields[0], "attachments/"+fields[6], deadline)
	if err != nil {
		return Row{}, fmt.Errorf("line %d: %w", line, err)
	}
	return Row{Record: record, Workflow: workflow, Attachment: attachment}, nil
}

func ValidateBatch(batch Batch) error {
	if len(batch.Rows) == 0 {
		return errors.New("batch has no accepted rows")
	}
	seen := map[string]bool{}
	for _, row := range batch.Rows {
		if seen[row.Record.ID] {
			return fmt.Errorf("duplicate record %s", row.Record.ID)
		}
		seen[row.Record.ID] = true
		if err := row.Record.Validate(); err != nil {
			return err
		}
		if err := row.Workflow.Validate(); err != nil {
			return err
		}
		if err := row.Attachment.Validate(); err != nil {
			return err
		}
		if !row.Attachment.IsSafeName() {
			return fmt.Errorf("unsafe attachment name %s", row.Attachment.Name)
		}
	}
	return nil
}

func AcceptedIDs(batch Batch) []string {
	ids := make([]string, 0, len(batch.Rows))
	for _, row := range batch.Rows {
		ids = append(ids, row.Record.ID)
	}
	return ids
}
