package filter_test

import (
	"fmt"
	"os"
	"path"

	"github.com/containerd/continuity"
	"github.com/depot/orca/filter"
)

func ExampleNewDockerIgnoreContext() {
	// Make some test data:
	root, _ := os.MkdirTemp("", "testwalk")
	f, err := os.Create(path.Join(root, "index.ts"))
	if err != nil {
		panic(err)
	}
	_ = f.Close()

	// This dockerignore filter will ignore all files in node_modules.
	ignoreNodeModules, err := filter.NewDockerIgnoreFilter([]string{"node_modules"})
	if err != nil {
		panic(err)
	}

	// Create a context that will walk the file tree and filter out node_modules.
	ctx, err := filter.NewDockerIgnoreContext(root, ignoreNodeModules, continuity.ContextOptions{
		// This digester will create chunked digests for each file.
		Digester: &filter.Digester{},
	})
	if err != nil {
		panic(err)
	}

	manifest, _ := continuity.BuildManifest(ctx)
	for _, r := range manifest.Resources {
		fmt.Printf("%v\n", r.Path())
	}
	// Output:
	// /index.ts
}
