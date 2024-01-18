package filter

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"github.com/containerd/continuity"
)

var _ continuity.Context = (*DockerIgnoreContext)(nil)

// DockerIgnoreContext wrap a continuity.Context that filters resources based on a .dockerignore file.
type DockerIgnoreContext struct {
	root         string
	dockerignore *DockerIgnoreFilter
	inner        continuity.Context
}

// NewDockerIgnoreContext creates a new DockerIgnoreContext with a dockerignore filter.
func NewDockerIgnoreContext(root string, dockerignore *DockerIgnoreFilter, opts continuity.ContextOptions) (*DockerIgnoreContext, error) {
	inner, err := continuity.NewContextWithOptions(root, opts)
	if err != nil {
		return nil, err
	}

	return &DockerIgnoreContext{
		root:         root,
		dockerignore: dockerignore,
		inner:        inner,
	}, nil
}

// Apply passes through to the inner context.
func (c *DockerIgnoreContext) Apply(r continuity.Resource) error { return c.inner.Apply(r) }

// Verify passes through to the inner context.
func (c *DockerIgnoreContext) Verify(r continuity.Resource) error { return c.inner.Verify(r) }

// Resource passes through to the inner context.
func (c *DockerIgnoreContext) Resource(p string, fi os.FileInfo) (continuity.Resource, error) {
	return c.inner.Resource(p, fi)
}

// Walk applies a dockerignore to filter files from a directory tree.
// Files that are not filtered by the dockerignore are passed through to the inner context.
//
// This was transliterated from buildkit. It has *BIG* assumption that it is ok for the same
// directory to be visited multiple times.
func (c *DockerIgnoreContext) Walk(walkFn filepath.WalkFunc) error {
	type visitedDir struct {
		info             os.FileInfo
		path             string
		walkPath         string
		pathWithSep      string
		excludeMatchInfo []bool
		visited          bool
	}

	var (
		parentDirs []visitedDir
	)

	// This is transliterated from buildkit. It can generate more than one match per match
	// as parent directories are potentially revisited. This means it has an _assumption_
	// that any state that is tracked per match is inserted into maps to ensure uniqueness.
	walk := func(walkPath string, info fs.FileInfo, walkErr error) error {
		path, err := filepath.Rel(c.root, filepath.Join(c.root, walkPath))
		if err != nil {
			return err
		}
		// Skip root
		if path == "." {
			return nil
		}

		skip := false
		var parentExcludeMatchInfo []bool
		if len(parentDirs) != 0 {
			parentExcludeMatchInfo = parentDirs[len(parentDirs)-1].excludeMatchInfo
		}
		matches, matchInfo, err := c.dockerignore.Matches(path, parentExcludeMatchInfo)
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
		isDir := info != nil && info.IsDir()
		if isDir {
			dir = visitedDir{
				info:             info,
				path:             path,
				walkPath:         walkPath,
				pathWithSep:      path + string(filepath.Separator),
				excludeMatchInfo: matchInfo,
			}
		}

		if matches {
			if isDir && c.dockerignore.OnlySimplePatterns {
				// Optimization: we can skip walking this dir if no
				// exceptions to exclude patterns could match anything
				// inside it.
				if !c.dockerignore.HasExclusions {
					return filepath.SkipDir
				}

				dirSlash := path + string(filepath.Separator)
				for _, pat := range c.dockerignore.Patterns {
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

		// This revisits all parent directories just so that exclusions with ignored parent directories
		// have their parent directories included in the manifest.
		for i, parentDir := range parentDirs {
			if parentDir.visited {
				continue
			}

			err := walkFn(parentDir.walkPath, parentDir.info, nil)
			if err != nil {
				return err
			}

			parentDirs[i].visited = true
		}

		return walkFn(walkPath, info, nil)
	}

	return c.inner.Walk(walk)
}
