package filter_test

import (
	"math/rand"
	"os"
	"path"
	"testing"

	"github.com/depot/orca/filter"
)

func BenchmarkWalk50(b *testing.B)        { benchmarkWalk(50, patterns(), b) }
func BenchmarkWalk500(b *testing.B)       { benchmarkWalk(500, patterns(), b) }
func BenchmarkWalk5000(b *testing.B)      { benchmarkWalk(5000, patterns(), b) }
func BenchmarkWalk50000(b *testing.B)     { benchmarkWalk(50000, patterns(), b) }
func BenchmarkNodeWalk50(b *testing.B)    { benchmarkWalk(50, node_patterns(), b) }
func BenchmarkNodeWalk500(b *testing.B)   { benchmarkWalk(500, node_patterns(), b) }
func BenchmarkNodeWalk5000(b *testing.B)  { benchmarkWalk(5000, node_patterns(), b) }
func BenchmarkNodeWalk50000(b *testing.B) { benchmarkWalk(50000, node_patterns(), b) }

func benchmarkWalk(numEntries int, patterns []string, b *testing.B) {
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

	root, err := os.MkdirTemp("", "walkbench")
	if err != nil {
		b.Fatal(err)
	}

	defer func() {
		_ = os.RemoveAll(root)
	}()

	for _, p := range paths {
		p = path.Join(root, p)
		dir := path.Dir(p)
		err := os.MkdirAll(dir, 0700)
		if err != nil {
			b.Fatal(err)
		}

		f, err := os.Create(p)
		if err != nil {
			b.Fatal(err)
		}
		_ = f.Close()
	}

	dockerignore, err := filter.NewDockerIgnoreFilter(patterns)
	if err != nil {
		b.Fatal(err)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err = filter.Walk(root, dockerignore)
		if err != nil {
			b.Fatal(err)
		}
	}
}
