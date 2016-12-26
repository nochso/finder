package finder

import (
	"os"
	"path"
	"strings"
)

// ItemSlice is a list of items.
type ItemSlice []Item

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
	return path.Join(i.base, i.path)
}

// RelPath returns the path relative to the searched directory.
func (i Item) RelPath() string {
	return i.path
}

// String returns the path
func (i Item) String() string {
	return i.Path()
}
