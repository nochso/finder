package finder

import (
	"errors"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/gobwas/glob"
)

type itemType uint8

const (
	typeAll itemType = iota
	typeFile
	typeDir
)

// Finder contains options for a file/directory search.
type Finder struct {
	dirs        []string
	names       []Matcher
	paths       []Matcher
	notPaths    []Matcher
	notNames    []Matcher
	sizes       []Matcher
	userFilters []Matcher
	minDepth    int
	maxDepth    int
	setupErrors []error
	itype       itemType
}

// Matcher checks whether an Item matches.
type Matcher func(Item) bool

// New returns a new finder.
//
// By default it will search for both files and directories.
func New() *Finder {
	return &Finder{maxDepth: -1}
}

// In searches in the given list of directories.
func (f *Finder) In(directories ...string) *Finder {
	f.dirs = append(f.dirs, directories...)
	return f
}

// Path narrows down the folders to be searched using gobwas/glob
//
// p is matched against the item's RelPath()
// See https://github.com/gobwas/glob
func (f *Finder) Path(p string) *Finder {
	matcher := f.path(p)
	if matcher != nil {
		f.paths = append(f.paths, matcher)
	}
	return f
}

// NotPath excludes folders from the search using gobwas/glob
//
// p is matched against the item's RelPath()
// See https://github.com/gobwas/glob
func (f *Finder) NotPath(p string) *Finder {
	matcher := f.path(p)
	if matcher != nil {
		f.notPaths = append(f.notPaths, matcher)
	}
	return f
}

func (f *Finder) path(p string) Matcher {
	g, err := glob.Compile(p, os.PathSeparator)
	if err != nil {
		f.setupErrors = append(f.setupErrors, err)
		return nil
	}
	return func(i Item) bool {
		if i.IsDir() {
			return g.Match(i.RelPath())
		}
		return g.Match(filepath.Dir(i.RelPath()))
	}
}

// Name matches a file or directory name using gobwas/glob
//
// See https://github.com/gobwas/glob
func (f *Finder) Name(n string) *Finder {
	matcher := f.name(n)
	if matcher != nil {
		f.names = append(f.names, matcher)
	}
	return f
}

func (f *Finder) name(n string) Matcher {
	g, err := glob.Compile(n, os.PathSeparator)
	if err != nil {
		f.setupErrors = append(f.setupErrors, err)
		return nil
	}
	return func(i Item) bool {
		return g.Match(i.Name())
	}
}

// NameRegex matches a file or directory name using package regexp.
func (f *Finder) NameRegex(n string) *Finder {
	matcher := f.nameRegex(n)
	if matcher != nil {
		f.names = append(f.names, matcher)
	}
	return f
}

func (f *Finder) nameRegex(n string) Matcher {
	re, err := regexp.Compile(n)
	if err != nil {
		f.setupErrors = append(f.setupErrors, err)
		return nil
	}
	return func(i Item) bool {
		return re.MatchString(i.Name())
	}
}

// NotName excludes a file or directory name using gobwas/glob
//
// See https://github.com/gobwas/glob
func (f *Finder) NotName(n string) *Finder {
	matcher := f.name(n)
	if matcher != nil {
		f.notNames = append(f.notNames, matcher)
	}
	return f
}

// NotNameRegex excludes a file or directory name using package regexp.
func (f *Finder) NotNameRegex(n string) *Finder {
	matcher := f.nameRegex(n)
	if matcher != nil {
		f.notNames = append(f.notNames, matcher)
	}
	return f
}

func (f *Finder) IgnoreVCS() *Finder {
	excludes := []string{
		".svn", "_svn", "CVS", "_darcs", ".arch-params", ".monotone", ".bzr",
		".git", ".hg",
	}
	exNames := make(map[string]bool, len(excludes))
	for _, exName := range excludes {
		exNames[exName] = true
	}
	matcher := func(i Item) bool {
		_, exists := exNames[i.Name()]
		return exists
	}
	f.notPaths = append(f.notPaths, matcher)
	return f
}

func (f *Finder) IgnoreDots() *Finder {
	matcher := func(i Item) bool {
		return strings.HasPrefix(i.Name(), ".")
	}
	f.notPaths = append(f.notPaths, matcher)
	return f
}

// Files makes the finder return files only.
func (f *Finder) Files() *Finder {
	f.itype = typeFile
	return f
}

// Dirs makes the finder return directories only.
func (f *Finder) Dirs() *Finder {
	f.itype = typeDir
	return f
}

func (f *Finder) MaxDepth(max int) *Finder {
	f.maxDepth = max
	return f
}

func (f *Finder) MinDepth(min int) *Finder {
	f.minDepth = min
	return f
}

// Filter items using a custom Matcher.
//	func(Item) bool
func (f *Finder) Filter(m Matcher) *Finder {
	f.userFilters = append(f.userFilters, m)
	return f
}

// Size filters by minimum and maximum file size.
// Max is ignored if it's lower than min.
//
//	Size(0, 1024)    // <=1kB
//	Size(1024, 1024) // ==1kB
//	Size(1024, -1)   // >=1kB
func (f *Finder) Size(min, max int64) *Finder {
	m := func(i Item) bool {
		return i.Size() >= min && (max >= min && i.Size() <= max)
	}
	f.sizes = append(f.sizes, m)
	return f
}

var IsMatch error = nil
var ErrNoMatch = errors.New("Item did not match")
var ErrSkipDir = filepath.SkipDir

// fast excludes first, followed by more expensive path matching
func (f *Finder) match(i Item) error {
	if (f.itype == typeDir && !i.IsDir()) || (f.itype == typeFile && i.IsDir()) {
		return ErrNoMatch
	}
	match := f.matchDepth(i)
	if match != IsMatch {
		return match
	}
	match = f.matchSize(i)
	if match != IsMatch {
		return match
	}
	match = f.matchUserFilters(i)
	if match != IsMatch {
		return match
	}
	match = f.matchPaths(i)
	if match != IsMatch {
		return match
	}
	match = f.matchNotPaths(i)
	if match != IsMatch {
		return match
	}
	match = f.matchNames(i)
	if match != IsMatch {
		return match
	}
	return f.matchNotNames(i)
}

func (f *Finder) matchSize(i Item) error {
	if len(f.sizes) == 0 {
		return IsMatch
	}
	for _, s := range f.sizes {
		if s(i) {
			return IsMatch
		}
	}
	return ErrNoMatch
}

func (f *Finder) matchUserFilters(i Item) error {
	if len(f.userFilters) == 0 {
		return IsMatch
	}
	for _, f := range f.userFilters {
		if f(i) {
			return IsMatch
		}
	}
	return ErrNoMatch
}

func (f *Finder) matchDepth(i Item) error {
	if f.maxDepth == -1 && f.minDepth == 0 {
		return IsMatch
	}
	depth := i.Depth()
	if f.minDepth != 0 && depth < f.minDepth {
		return ErrNoMatch
	}
	if f.maxDepth != -1 && depth > f.maxDepth {
		return ErrSkipDir
	}
	return IsMatch
}

func (f *Finder) matchPaths(i Item) error {
	if len(f.paths) == 0 {
		return IsMatch
	}
	if i.IsDir() {
		for _, p := range f.paths {
			if !p(i) {
				return ErrSkipDir
			}
		}
		return IsMatch
	}
	match := ErrNoMatch
	for _, p := range f.paths {
		if p(i) {
			match = IsMatch
		}
	}
	return match
}

func (f *Finder) matchNotPaths(i Item) error {
	if len(f.notPaths) == 0 {
		return IsMatch
	}
	if i.IsDir() {
		for _, p := range f.notPaths {
			if p(i) {
				return ErrSkipDir
			}
		}
		return IsMatch
	}
	for _, p := range f.notPaths {
		if p(i) {
			return ErrNoMatch
		}
	}
	return IsMatch
}

func (f *Finder) matchNames(i Item) error {
	if len(f.names) == 0 {
		return IsMatch
	}
	for _, n := range f.names {
		if n(i) {
			return IsMatch
		}
	}
	return ErrNoMatch
}

func (f *Finder) matchNotNames(i Item) error {
	if len(f.notNames) == 0 {
		return IsMatch
	}
	for _, n := range f.notNames {
		if n(i) {
			return ErrNoMatch
		}
	}
	return IsMatch
}

// Each calls func fn with each found item.
func (f *Finder) Each(fn func(Item)) []error {
	var errs []error
	var dir string
	walker := func(path string, info os.FileInfo, err error) error {
		if dir == path {
			return nil
		}
		if err != nil {
			return err
		}
		relDir, err := filepath.Rel(dir, path)
		if err != nil {
			return err
		}
		item := newItem(info, dir, relDir)
		match := f.match(item)
		switch match {
		case IsMatch:
			fn(item)
			return nil
		case ErrNoMatch:
			return nil
		case ErrSkipDir:
			return ErrSkipDir
		}
		return nil
	}
	for _, dir = range f.dirs {
		err := filepath.Walk(dir, walker)
		if err != nil && err != ErrSkipDir {
			errs = append(errs, err)
		}
	}
	return append(f.setupErrors, errs...)
}

// ToSlice returns a slice of all found items.
func (f *Finder) ToSlice() (ItemSlice, []error) {
	var l []Item
	errs := f.Each(func(file Item) {
		l = append(l, file)
	})
	return l, errs
}
