package filter

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"syscall"

	continuityapi "github.com/containerd/continuity/proto"
	"github.com/moby/patternmatcher"
	"golang.org/x/sys/unix"
)

func Walk(root string, dockerignore *DockerIgnoreFilter) (*continuityapi.Manifest, error) {
	var (
		parentDirs  []visitedDir
		pathMatches []match
	)

	filepath.WalkDir(root, func(path string, entry fs.DirEntry, walkErr error) error {
		skip := false
		var parentExcludeMatchInfo []bool
		if len(parentDirs) != 0 {
			parentExcludeMatchInfo = parentDirs[len(parentDirs)-1].excludeMatchInfo
		}
		matches, matchInfo, err := dockerignore.Matches(path, parentExcludeMatchInfo)
		if err != nil {
			return fmt.Errorf("failed to match patterns: %w", err)
		}

		for len(parentDirs) != 0 {
			lastParentDir := parentDirs[len(parentDirs)-1].pathWithSep
			if strings.HasPrefix(path, lastParentDir) {
				break
			}
			parentDirs = parentDirs[:len(parentDirs)-1]
		}

		var dir visitedDir
		isDir := entry != nil && entry.IsDir()
		if isDir {
			dir = visitedDir{
				entry:            entry,
				path:             path,
				pathWithSep:      path + string(filepath.Separator),
				excludeMatchInfo: matchInfo,
			}
		}

		if matches {
			if isDir && dockerignore.OnlySimplePatterns {
				// Optimization: we can skip walking this dir if no
				// exceptions to exclude patterns could match anything
				// inside it.
				if !dockerignore.HasExclusions {
					return filepath.SkipDir
				}

				dirSlash := path + string(filepath.Separator)
				for _, pat := range dockerignore.Patterns {
					if !pat.Exclusion {
						continue
					}
					patStr := patternWithoutTrailingGlob(pat) + string(filepath.Separator)
					if strings.HasPrefix(patStr, dirSlash) {
						goto passedExcludeFilter
					}
				}
				return filepath.SkipDir
			}
		passedExcludeFilter:
			skip = true
		}

		if walkErr != nil {
			if skip && errors.Is(walkErr, os.ErrPermission) {
				return nil
			}
			return walkErr
		}

		if isDir {
			parentDirs = append(parentDirs, dir)
		}

		if skip {
			return nil
		}

		for i, parentDir := range parentDirs {
			if parentDir.visited {
				continue
			}

			pathMatches = append(pathMatches, match{
				Path:  parentDir.path,
				Entry: parentDir.entry,
			})

			parentDirs[i].visited = true
		}

		pathMatches = append(pathMatches, match{
			Path:  path,
			Entry: entry,
		})

		return nil
	})

	resources := make([]*continuityapi.Resource, 0, len(pathMatches))
	for _, match := range pathMatches {
		info, err := match.Entry.Info()
		if err != nil {
			continue
		}

		resource := &continuityapi.Resource{
			Path: []string{match.Path},

			Size: uint64(info.Size()),
			// TODO: Handle symlinks
			// Target:

			Mode: uint32(info.Mode()),
		}

		if s, ok := info.Sys().(*syscall.Stat_t); ok {
			resource.Uid = int64(s.Uid)
			resource.Gid = int64(s.Gid)
			resource.Major = uint64(unix.Major(uint64(s.Rdev)))
			resource.Minor = uint64(unix.Minor(uint64(s.Rdev)))
		}

		// TODO: xattrs
		// TODO: windows ads

		resources = append(resources, resource)
	}

	return &continuityapi.Manifest{
		Resource: resources,
	}, nil
}

func patternWithoutTrailingGlob(pattern *patternmatcher.Pattern) string {
	patStr := pattern.CleanedPattern
	// We use filepath.Separator here because patternmatcher.Pattern patterns
	// get transformed to use the native path separator:
	// https://github.com/moby/patternmatcher/blob/130b41bafc16209dc1b52a103fdac1decad04f1a/patternmatcher.go#L52
	patStr = strings.TrimSuffix(patStr, string(filepath.Separator)+"**")
	patStr = strings.TrimSuffix(patStr, string(filepath.Separator)+"*")
	return patStr
}

type match struct {
	Path  string
	Entry fs.DirEntry
}

type visitedDir struct {
	entry            fs.DirEntry
	path             string
	pathWithSep      string
	excludeMatchInfo []bool
	visited          bool
}
