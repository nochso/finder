package finder

import (
	"path/filepath"
	"sort"
)

// ItemSlice is a list of items.
type ItemSlice []Item

// ByName sorts by the base name in ascending order.
func ByName(i, j Item) bool {
	return i.Name() < j.Name()
}

// ByPath sorts by path in ascending order.
func ByPath(i, j Item) bool {
	return i.Path() < j.Path()
}

// ByModified sorts by modification time in ascending order.
func ByModified(i, j Item) bool {
	return i.ModTime().Before(j.ModTime())
}

// BySize sorts by file size in ascending order.
func BySize(i, j Item) bool {
	return i.Size() < j.Size()
}

// ByExtension sorts by file extension in ascending order.
func ByExtension(i, j Item) bool {
	return filepath.Ext(i.Name()) < filepath.Ext(j.Name())
}

// Sort the slice given a function that mimics `Less` of `sort.Interface`.
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

// ToStringSlice returns a string slice with each item's path.
func (is ItemSlice) ToStringSlice() []string {
	ss := make([]string, 0, len(is))
	for _, item := range is {
		ss = append(ss, item.String())
	}
	return ss
}
