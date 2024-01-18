package filter

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/moby/patternmatcher"
)

// MatchIndices is a slice of booleans that indicate which patterns matched the path.
type MatchedPatterns []bool

// DockerIgnoreFilter is a filter that uses dockerignore rules  to filter files from a directory tree.
//
// This is largely transliterated from buildkit.
type DockerIgnoreFilter struct {
	Patterns           []*patternmatcher.Pattern
	HasExclusions      bool
	OnlySimplePatterns bool
}

func NewDockerIgnoreFilter(patterns []string) (*DockerIgnoreFilter, error) {
	ps, err := patternmatcher.NewPatterns(patterns)
	if err != nil {
		return nil, err
	}

	var hasExclusions bool
	for _, p := range ps {
		if p.Exclusion {
			hasExclusions = true
			break
		}
	}

	onlySimplePatterns := len(ComplexPatterns(ps)) == 0

	return &DockerIgnoreFilter{
		Patterns:           ps,
		HasExclusions:      hasExclusions,
		OnlySimplePatterns: onlySimplePatterns,
	}, nil
}

// Matches checks if the path matches the dockerignore patterns.
func (f *DockerIgnoreFilter) Matches(path string, parentPathMatches MatchedPatterns) (ok bool, patterns MatchedPatterns, err error) {
	ok, patterns, err = patternmatcher.MatchesUsingParentResults(f.Patterns, path, parentPathMatches)
	return
}

// ComplexPatterns checks if all patterns are prefix patterns.
func ComplexPatterns(patterns []*patternmatcher.Pattern) []*patternmatcher.Pattern {
	patternChars := "*[]?^"
	if os.PathSeparator != '\\' {
		patternChars += `\`
	}

	var complexPatterns []*patternmatcher.Pattern

	for _, p := range patterns {
		if p.Exclusion && strings.ContainsAny(patternWithoutTrailingGlob(p), patternChars) {
			complexPatterns = append(complexPatterns, p)
		}
	}
	return complexPatterns
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
