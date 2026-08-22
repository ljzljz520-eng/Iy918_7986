package importer

import (
	"fmt"
	"sort"
)

type Report struct {
	Accepted int
	Rejected int
	IDs      []string
	Message  string
}

func BuildReport(batch Batch) Report {
	ids := AcceptedIDs(batch)
	sort.Strings(ids)
	message := fmt.Sprintf("accepted=%d rejected=%d", len(ids), len(batch.Rejected))
	return Report{Accepted: len(ids), Rejected: len(batch.Rejected), IDs: ids, Message: message}
}

func MergeReports(reports ...Report) Report {
	merged := Report{}
	for _, report := range reports {
		merged.Accepted += report.Accepted
		merged.Rejected += report.Rejected
		merged.IDs = append(merged.IDs, report.IDs...)
	}
	sort.Strings(merged.IDs)
	merged.Message = fmt.Sprintf("accepted=%d rejected=%d", merged.Accepted, merged.Rejected)
	return merged
}
