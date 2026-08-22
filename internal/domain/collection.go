package domain

import "sort"

type RecordCollection struct {
	Items []Record
}

func NewRecordCollection(items []Record) RecordCollection {
	return RecordCollection{Items: append([]Record(nil), items...)}
}

func (c RecordCollection) Len() int {
	return len(c.Items)
}

func (c RecordCollection) IDs() []string {
	ids := make([]string, 0, len(c.Items))
	for _, item := range c.Items {
		ids = append(ids, item.ID)
	}
	sort.Strings(ids)
	return ids
}

func (c RecordCollection) ByStatus(status string) RecordCollection {
	items := make([]Record, 0)
	for _, item := range c.Items {
		if item.Status == status {
			items = append(items, item)
		}
	}
	return NewRecordCollection(items)
}

func (c RecordCollection) ByOwner(owner string) RecordCollection {
	items := make([]Record, 0)
	for _, item := range c.Items {
		if item.Owner == owner {
			items = append(items, item)
		}
	}
	return NewRecordCollection(items)
}

func (c RecordCollection) SortByUpdated(descending bool) RecordCollection {
	items := append([]Record(nil), c.Items...)
	sort.Slice(items, func(i, j int) bool {
		if descending {
			return items[i].UpdatedAt > items[j].UpdatedAt
		}
		return items[i].UpdatedAt < items[j].UpdatedAt
	})
	return NewRecordCollection(items)
}

func (c RecordCollection) CountByStatus() map[string]int {
	counts := make(map[string]int)
	for _, item := range c.Items {
		counts[item.Status]++
	}
	return counts
}
