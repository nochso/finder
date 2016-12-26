package finder

import (
	"os"
	"path/filepath"
	"sort"
	"strings"
)

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

// Item represents a single file or directory.
type Item struct {
	os.FileInfo
	base string
	path string
}

func newItem(info os.FileInfo, base, path string) Item {
	return Item{
		FileInfo: info,
		base:     base,
		path:     path,
	}
}

// Depth returns the depth of the item based on RelPath.
// It starts at 1 (one).
func (i Item) Depth() int {
	rp := i.RelPath()
	return strings.Count(rp, string(os.PathSeparator)) + 1
}

// Path returns the path including the folder it was found in.
func (i Item) Path() string {
	return filepath.Join(i.base, i.path)
}

// RelPath returns the path relative to the searched directory.
func (i Item) RelPath() string {
	return i.path
}

// String returns the path
func (i Item) String() string {
	return i.Path()
}
