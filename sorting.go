package deepunique

import (
	"sort"
)

type indSort struct {
	items [][2]any
	index []string
}

func (sbo indSort) Len() int {
	return len(sbo.index)
}

func (sbo indSort) Swap(i, j int) {
	sbo.items[i], sbo.items[j] = sbo.items[j], sbo.items[i]
	sbo.index[i], sbo.index[j] = sbo.index[j], sbo.index[i]
}

func (sbo indSort) Less(i, j int) bool {
	return sbo.index[i] < sbo.index[j]
}

func SortMapTuples(items [][2]any, index []string) {
	ts := indSort{items, index}
	sort.Sort(ts)
}
