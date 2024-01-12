package snapshots

import (
	"context"
	"fmt"

	"github.com/containerd/containerd"
	"github.com/containerd/containerd/diff"
	"github.com/containerd/containerd/snapshots"
	ocispec "github.com/opencontainers/image-spec/specs-go/v1"
)

func Compress(ctx context.Context, client *containerd.Client, snapshotter snapshots.Snapshotter, parentKey, committedKey string) (ocispec.Descriptor, error) {
	parentViewKey := fmt.Sprintf("%s-view-key", parentKey)
	parentMounts, err := snapshotter.View(ctx, parentViewKey, parentKey)
	if err != nil {
		return ocispec.Descriptor{}, err
	}

	defer snapshotter.Remove(ctx, parentViewKey)

	committedViewKey := fmt.Sprintf("%s-view-key", committedKey)
	committedMounts, err := snapshotter.View(ctx, committedViewKey, committedKey)
	if err != nil {
		return ocispec.Descriptor{}, err
	}

	defer snapshotter.Remove(ctx, committedViewKey)
	opts := []diff.Opt{
		diff.WithMediaType(ocispec.MediaTypeImageLayerGzip),
		//diff.WithReference(TODO REFERENCE?),
		//diff.WithLabels(labels),
	}
	return client.DiffService().Compare(ctx, parentMounts, committedMounts, opts...)
}
