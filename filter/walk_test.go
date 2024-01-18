package filter_test

import (
	"math/rand"
	"os"
	"path"
	"testing"

	"github.com/containerd/continuity"
	"github.com/depot/orca/filter"
)

func TestWalk(t *testing.T) {
	root, err := os.MkdirTemp("", "testwalk")
	if err != nil {
		t.Fatal(err)
	}
	defer func() {
		_ = os.RemoveAll(root)
	}()

	/*
		Creates a file tree like this:
		howdy
		└── doody
		    └── main.go

		2 directories, 1 file
	*/
	err = os.MkdirAll(path.Join(root, "howdy", "doody"), 0700)
	if err != nil {
		t.Fatal(err)
	}

	f, err := os.Create(path.Join(root, "howdy", "doody", "main.go"))
	if err != nil {
		t.Fatal(err)
	}
	_ = f.Close()

	ignorehowdy, err := filter.NewDockerIgnoreFilter([]string{"howdy"})
	if err != nil {
		t.Fatal(err)
	}

	ctx, err := filter.NewDockerIgnoreContext(root, ignorehowdy, continuity.ContextOptions{})
	if err != nil {
		t.Fatal(err)
	}

	manifest, err := continuity.BuildManifest(ctx)
	if err != nil {
		t.Fatal(err)
	}

	if len(manifest.Resources) != 0 {
		t.Fatalf("expected no resources, got %d", len(manifest.Resources))
	}
}

func TestWalkNoneFiltered(t *testing.T) {
	root, err := os.MkdirTemp("", "testwalk")
	if err != nil {
		t.Fatal(err)
	}
	defer func() {
		_ = os.RemoveAll(root)
	}()

	/*
		Creates a file tree like this:
		howdy
		└── doody
		    └── main.go

		2 directories, 1 file
	*/
	err = os.MkdirAll(path.Join(root, "howdy", "doody"), 0700)
	if err != nil {
		t.Fatal(err)
	}

	f, err := os.Create(path.Join(root, "howdy", "doody", "main.go"))
	if err != nil {
		t.Fatal(err)
	}
	_ = f.Close()

	filterNothing, err := filter.NewDockerIgnoreFilter([]string{"doesntexist"})
	if err != nil {
		t.Fatal(err)
	}

	ctx, err := filter.NewDockerIgnoreContext(root, filterNothing, continuity.ContextOptions{})
	if err != nil {
		t.Fatal(err)
	}

	manifest, err := continuity.BuildManifest(ctx)
	if err != nil {
		t.Fatal(err)
	}

	if len(manifest.Resources) != 3 {
		t.Fatalf("expected three resources, got %d", len(manifest.Resources))
	}
}

func TestWalkExclusions(t *testing.T) {
	root, err := os.MkdirTemp("", "testwalk")
	if err != nil {
		t.Fatal(err)
	}
	defer func() {
		_ = os.RemoveAll(root)
	}()

	/*
		Creates a file tree like this:
		howdy
		└── doody
		    └── main.go

		2 directories, 1 file
	*/
	err = os.MkdirAll(path.Join(root, "howdy", "doody"), 0700)
	if err != nil {
		t.Fatal(err)
	}

	f, err := os.Create(path.Join(root, "howdy", "doody", "main.go"))
	if err != nil {
		t.Fatal(err)
	}
	_ = f.Close()

	exclusion, err := filter.NewDockerIgnoreFilter([]string{"howdy", "!howdy/doody/**"})
	if err != nil {
		t.Fatal(err)
	}

	ctx, err := filter.NewDockerIgnoreContext(root, exclusion, continuity.ContextOptions{})
	if err != nil {
		t.Fatal(err)
	}

	manifest, err := continuity.BuildManifest(ctx)
	if err != nil {
		t.Fatal(err)
	}

	if len(manifest.Resources) != 3 {
		t.Fatalf("expected three resources, got %d", len(manifest.Resources))
	}
}

// TestWalkSkipExclusions tests that exclusions are skipped if the excluded directory
// is not explicitly prefixed with filtered directories.
func TestWalkSkipExclusions(t *testing.T) {
	root, err := os.MkdirTemp("", "testwalk")
	if err != nil {
		t.Fatal(err)
	}
	defer func() {
		_ = os.RemoveAll(root)
	}()

	/*
		Creates a file tree like this:
		howdy
		└── doody
		    └── main.go

		2 directories, 1 file
	*/
	err = os.MkdirAll(path.Join(root, "howdy", "doody"), 0700)
	if err != nil {
		t.Fatal(err)
	}

	f, err := os.Create(path.Join(root, "howdy", "doody", "main.go"))
	if err != nil {
		t.Fatal(err)
	}
	_ = f.Close()

	exclusion, err := filter.NewDockerIgnoreFilter([]string{"howdy", "!doody/**"})
	if err != nil {
		t.Fatal(err)
	}

	ctx, err := filter.NewDockerIgnoreContext(root, exclusion, continuity.ContextOptions{})
	if err != nil {
		t.Fatal(err)
	}

	manifest, err := continuity.BuildManifest(ctx)
	if err != nil {
		t.Fatal(err)
	}

	if len(manifest.Resources) != 0 {
		t.Fatalf("expected no resources, got %d", len(manifest.Resources))
	}
}

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

	ctx, err := filter.NewDockerIgnoreContext(root, dockerignore, continuity.ContextOptions{})
	if err != nil {
		b.Fatal(err)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := continuity.BuildManifest(ctx)
		if err != nil {
			b.Fatal(err)
		}
	}
}
