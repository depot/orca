package main

import (
	"context"
	"fmt"
	"os"
	"path"
	"sync"

	"github.com/containerd/containerd"
	orca "github.com/depot/orca/client"
	"github.com/depot/orca/image"
	"github.com/depot/orca/snapshots"
)

func main() {
	client, err := orca.NewClient("")
	if err != nil {
		panic(err)
	}
	defer client.Close()

	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)

	events, errs := client.EventService().Subscribe(ctx)
	wg := &sync.WaitGroup{}
	wg.Add(1)
	go func() {
		defer wg.Done()
		for {
			select {
			case <-ctx.Done():
				return
			case event := <-events:
				fmt.Printf("Event: %v\n", event)
			case err := <-errs:
				fmt.Printf("Error: %+v\n", err)
			}
		}
	}()

	ctx = orca.WithNamespace(ctx, "simple")
	ref := "docker.io/library/alpine:latest"
	unpackedImage, err := image.Pull(ctx, client, ref)
	if err != nil {
		panic(err)
	}

	snapshotter := client.SnapshotService(containerd.DefaultSnapshotter)

	mountPoint := "/home/goller/src/experiments/2023-12-19-containerd/mnt/root"
	committer, err := snapshots.Prepare(ctx, snapshotter, unpackedImage.SnapshotKey, mountPoint)
	if err != nil {
		panic(err)
	}

	filename := path.Join(mountPoint, "howdy")
	h, err := os.OpenFile(filename, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		panic(err)
	}
	_, err = h.Write([]byte("hello world"))
	if err != nil {
		panic(err)
	}
	_ = h.Close()

	fmt.Printf("FILE CREATED %s\n", filename)

	err = committer.Unmount()
	if err != nil {
		panic(err)
	}

	committedKey := "howdy"
	err = committer.Commit(ctx, snapshotter, committedKey)
	if err != nil {
		panic(err)
	}

	descriptor, err := snapshots.Compress(ctx, client, snapshotter, unpackedImage.SnapshotKey, committedKey)
	if err != nil {
		panic(err)
	}

	fmt.Printf("compressed layer: %+v\n", descriptor)

	cancel()
	wg.Wait()
}
