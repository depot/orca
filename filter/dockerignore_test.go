package filter_test

import (
	"fmt"
	"math/rand"
	"testing"

	"github.com/depot/orca/filter"
	"github.com/moby/patternmatcher"
)

func TestComplexPatterns(t *testing.T) {
	complex, _ := patternmatcher.NewPatterns([]string{"!README.*"})
	simple, _ := patternmatcher.NewPatterns([]string{"node_modules", "*.ts", "!LICENSE", "!src/*", "!src/**"})
	dockerignore, _ := patternmatcher.NewPatterns(patterns())
	node, _ := patternmatcher.NewPatterns(node_patterns())

	if len(filter.ComplexPatterns(complex)) != 1 {
		t.Fatal("expected complex patterns to be complex")
	}
	if len(filter.ComplexPatterns(simple)) != 0 {
		t.Fatal("expected simple patterns to be simple")
	}
	if len(filter.ComplexPatterns(dockerignore)) != 0 {
		t.Fatal("expected test dockerignore patterns to be simple")
	}
	if len(filter.ComplexPatterns(node)) != 0 {
		t.Fatal("expected test node patterns to be simple")
	}
}

// Conclusion: The allocations of strings.Split and strings.Join are the largest problems.
func BenchmarkMatched50(b *testing.B)    { benchmarkMatched(50, patterns(), b) }
func BenchmarkMatched500(b *testing.B)   { benchmarkMatched(500, patterns(), b) }
func BenchmarkMatched5000(b *testing.B)  { benchmarkMatched(5000, patterns(), b) }
func BenchmarkMatched50000(b *testing.B) { benchmarkMatched(50000, patterns(), b) }

func BenchmarkNodeMatched50(b *testing.B)    { benchmarkMatched(50, node_patterns(), b) }
func BenchmarkNodeMatched500(b *testing.B)   { benchmarkMatched(500, node_patterns(), b) }
func BenchmarkNodeMatched5000(b *testing.B)  { benchmarkMatched(5000, node_patterns(), b) }
func BenchmarkNodeMatched50000(b *testing.B) { benchmarkMatched(50000, node_patterns(), b) }

var randomSource = rand.New(rand.NewSource(42))

func benchmarkMatched(numEntries int, patterns []string, b *testing.B) {
	const (
		maxDepth        = 5
		randomStrLength = 12
	)

	paths := make([]string, numEntries)
	for i := 0; i < numEntries; i++ {
		depth := rand.Intn(maxDepth) + 1
		fileNameLength := rand.Intn(randomStrLength) + 1

		paths[i] = generateRandomName(depth, fileNameLength)
	}

	filter, err := filter.NewDockerIgnoreFilter(patterns)
	if err != nil {
		b.Fatal(err)
	}

	b.ResetTimer()
	var ok bool
	for i := 0; i < b.N; i++ {
		for _, path := range paths {
			ok, _, _ = filter.Matches(path, nil)
		}
	}
	result = ok
}

var result bool

func generateRandomName(depth, fileNameLength int) string {
	var (
		fileExtensions = []string{"ts", "js", "json", "md", "css", "html", "txt", "log"}
		directoryNames = []string{"src", "lib", "public", "test", "config", "node_modules", "scripts"}
	)

	path := ""
	for i := 0; i < depth; i++ {
		path += randomElement(directoryNames) + "/"
	}
	randomStr := generateRandomString(fileNameLength)
	return path + fmt.Sprintf("%s.%s", randomStr, randomElement(fileExtensions))
}

func generateRandomString(length int) string {
	const charset = "abcdefghijklmnopqrstuvwxyz"
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[randomSource.Intn(len(charset))]
	}
	return string(b)
}

func randomElement(slice []string) string {
	return slice[randomSource.Intn(len(slice))]
}

// Taken from patternmatcher_test.go.
func patterns() []string {
	return []string{
		"abc",
		"*c",
		"a*",
		"a*",
		"a*",
		"a*/b",
		"a*/b",
		"a*b*c*d*e*/f",
		"a*b*c*d*e*/f",
		"a*b*c*d*e*/f",
		"a*b*c*d*e*/f",
		"a*b?c*x",
		"a*b?c*x",
		"ab[c]",
		"ab[b-d]",
		"ab[e-g]",
		"ab[^c]",
		"ab[^b-d]",
		"ab[^e-g]",
		"a\\*b",
		"a\\*b",
		"a?b",
		"a[^a]b",
		"a???b",
		"a[^a][^a][^a]b",
		"[a-ζ]*",
		"*[a-ζ]",
		"a?b",
		"a*b",
		"[\\]a]",
		"[\\-]",
		"[x\\-]",
		"[x\\-]",
		"[x\\-]",
		"[\\-x]",
		"[\\-x]",
		"[\\-x]",
	}
}

func node_patterns() []string {
	return []string{
		"!LICENSE",
		"!scripts",
		"**/tmp",
		"**/node_modules",
		".cache",
		".git",
		".github",
		".moon/cache",
		".vscode",
		"packages/db/client",
		"public/build",
		".depcheckrc.yml",
		".dockerignore",
		".eslintrc.js",
		".eslintrc.json",
		".gitignore",
		".prettier*",
		"*.env*",
		"*.log",
		"*.tsbuildinfo",
		"docker-compose.yml",
		"Dockerfile",
		"fly.toml",
		"fly.*.toml",
		"gha-creds-*.json",
		"Makefile",
		"README.md",
		"sentry.properties",
	}
}
