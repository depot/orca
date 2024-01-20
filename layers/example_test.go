package layers_test

import (
	"fmt"
	"os"
	"path"

	"github.com/containerd/continuity"
	"github.com/depot/orca/filter"
	"github.com/depot/orca/layers"
)

func ExampleWalk() {
	// Make some test data:
	root, _ := os.MkdirTemp("", "testwalk")
	f, err := os.Create(path.Join(root, "index.ts"))
	if err != nil {
		panic(err)
	}
	f.Write([]byte(`console.log("hello world")`))
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

	layer, _ := layers.Walk(ctx)
	for _, entry := range layer.Entries {
		fmt.Printf("%v\n", entry.Path[0])
		fmt.Printf("%v\n", entry.Blocks[0])
	}
	// Output:
	// /index.ts
	// size_bytes:26 digest:{algorithm:ALGORITHM_XXH64 sum:15630286323523909447}
}
