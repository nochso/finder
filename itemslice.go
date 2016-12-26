package finder

import "sort"

// ItemSlice is a list of items.
type ItemSlice []Item

// Sort the item slice using the given `Less` implementation of the sort interface.
func (is ItemSlice) Sort(lessFn func(i, j Item) bool) {
	sort.Sort(sortableItemSlice{is, lessFn})
}

type sortableItemSlice struct {
	items ItemSlice
	less  func(i, j Item) bool
}

func (sil sortableItemSlice) Len() int {
	return len(sil.items)
}

func (sil sortableItemSlice) Swap(i, j int) {
	sil.items[i], sil.items[j] = sil.items[j], sil.items[i]
}

func (sil sortableItemSlice) Less(i, j int) bool {
	return sil.less(sil.items[i], sil.items[j])
}

// Size returns the sum of all file sizes.
func (is ItemSlice) Size() uint64 {
	var size uint64
	for _, i := range is {
		if !i.IsDir() {
			size += uint64(i.Size())
		}
	}
	return size
}
