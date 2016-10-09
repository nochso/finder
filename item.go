package finder

import (
	"os"
	"path"
	"strings"
)

// ItemSlice is a list of items.
type ItemSlice []Item

// Size returns the sum of all file sizes.
func (il ItemSlice) Size() uint64 {
	var size uint64
	for _, i := range il {
		if !i.IsDir() {
			size += uint64(i.Size())
		}
	}
	return size
}

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
func (f Item) Depth() int {
	rp := f.RelPath()
	return strings.Count(rp, string(os.PathSeparator)) + 1
}

// Path returns the path including the folder it was found in.
func (f Item) Path() string {
	return path.Join(f.base, f.path)
}

// RelPath returns the path relative to the searched directory.
func (f Item) RelPath() string {
	return f.path
}

// String returns the path
func (f Item) String() string {
	return f.Path()
}
