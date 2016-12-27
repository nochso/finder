Change Log
==========

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](http://keepachangelog.com/) 
and this project adheres to [Semantic Versioning](http://semver.org/).


[Unreleased]
------------

### Added

- New method `ItemSlice.Sort(func (i, j Item) bool)` sorts the slice given a function that mimics `Less` of `sort.Interface`.
- Common functions for `Sort()`: `ByName`, `ByPath`, `ByModified`, `BySize` and `ByExtension`.
- New method `ItemSlice.ToStringSlice()` returns a string slice of paths.
- New method `Finder.Error()` returns a `MultiErr` containing all errors that
  occured during setup or walking.


### Changed

- Completely skip folders that are deeper than the maximum depth. Restores
  behaviour of `<=0.1.0`, reducing time for shallow scans in nested directories.
- `Each` and `ToSlice` no longer return an error slice. Instead a `MultiErr` is
  available via method `Finder.Error()`


### Fixed

- `IgnoreVCS` happened to ignore files. Now only directories are ignored as the
  documentation already claimed.


[0.2.0] - 2016-12-26
--------------------

### Added

- This changelog.
- Method `Filter(func(Item) bool)` to filter using a custom matcher.
- Method `Depth(min, max int)` to filter based on directory depth.
  Replaces the removed methods `MinDepth` and `MaxDepth`.


### Changed

- ***BREAKING***: Change parameters of `Finder.Size()` from `func(size int64) bool` to `min, max int64`.


### Removed

- ***BREAKING*** Methods `MaxDepth` and `MinDepth`. Use `Depth` instead.


[0.1.0] - 2016-10-10
--------------------

### Added

- Initial public release


[Unreleased]: https://github.com/nochso/finder/compare/0.2.0...HEAD
[0.2.0]: https://github.com/nochso/finder/compare/0.1.0...0.2.0
[0.1.0]: https://github.com/nochso/finder/compare/a71aecf5b715e482a6b29121a271936f92aeea51...0.1.0