package filter

import "github.com/moby/patternmatcher"

var _ FileFilter = (*DockerIgnoreFilter)(nil)

type DockerIgnoreFilter struct {
	Patterns []*patternmatcher.Pattern
}

func NewDockerIgnoreFilter(patterns []string) (*DockerIgnoreFilter, error) {
	ps, err := patternmatcher.NewPatterns(patterns)
	if err != nil {
		return nil, err
	}

	return &DockerIgnoreFilter{
		Patterns: ps,
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
