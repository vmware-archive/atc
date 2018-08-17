package teamserver

import (
	"sort"
)

type extraBy func(item1, item2 *interface{}) bool

type Sorter struct {
	items   []interface{}interface
	extraBy func(item1, item2 *{}) bool
}

func (b extraBy) Sort(items []interface{}) {
	s := &Sorter{
		items:   items,
		extraBy: b,
	}
	sort.Sort(s)
}

func (s *Sorter) Len() int {
	return len(s.items)
}

func (s *Sorter) Swap(i, j int) {
	s.items[i], s.items[j] = s.items[j], s.items[i]
}

func (s *Sorter) Less(i, j int) bool {
	return s.extraBy(&s.items[i], &s.items[j])
}

func (s Sorter) GenericSort(sorter func(item1, item2 *interface{}) bool) interface{} {
	extraBy(sorter).Sort(s.items)

	return s.items
}
