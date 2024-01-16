package filter

// MatchIndices is a slice of booleans that indicate which patterns matched the path.
type MatchedPatterns []bool

type FileFilter interface {
	// Matches returns true if the path matches the filter.
	// The patterns slice contains a boolean for each pattern indicating whether it matched.
	Matches(path string, parentPathMatches MatchedPatterns) (ok bool, patterns MatchedPatterns, err error)
	// MatchedPatterns returns the patterns that matched the path.
	MatchedPatterns(patterns MatchedPatterns) []string
}
