package filter

import (
	"os"
	"strings"

	"github.com/moby/patternmatcher"
)

var _ FileFilter = (*DockerIgnoreFilter)(nil)

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

func (f *DockerIgnoreFilter) Matches(path string, parentPathMatches MatchedPatterns) (ok bool, patterns MatchedPatterns, err error) {
	ok, patterns, err = patternmatcher.MatchesUsingParentResults(f.Patterns, path, parentPathMatches)
	return
}

func (f *DockerIgnoreFilter) MatchedPatterns(patterns MatchedPatterns) []string {
	matched := make([]string, 0, len(patterns))
	for i, p := range patterns {
		if p {
			matched = append(matched, f.Patterns[i].CleanedPattern)
		}
	}
	return matched
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
