package finder

import "os"

type Item struct {
	os.FileInfo
	path string
}

func newItem(info os.FileInfo, path string) Item {
	return Item{
		FileInfo: info,
		path:     path,
	}
}

// Path returns the path based on the folder it was found in.
func (f Item) Path() string {
	return f.path
}

func (f Item) String() string {
	return f.path
}
