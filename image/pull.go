package image

import (
	"context"

	"github.com/containerd/containerd"
	"github.com/containerd/containerd/platforms"
	"github.com/opencontainers/image-spec/identity"
)

// TODO: add ref string helper (e.g. alpine:latest -> "docker.io/library/alpine:latest")

type UnpackedImage struct {
	Image       containerd.Image
	SnapshotKey string
}

// Pull pulls a remote image and unpacks it into the default snapshotter (typically, overlayfs).
// It returns an UnpackedImage which contains the image and the snapshot key.
//
// Containerd's snapshot key of an image is the ChainID of the image's diffIDs.
func Pull(ctx context.Context, client *containerd.Client, ref string) (*UnpackedImage, error) {
	opts := []containerd.RemoteOpt{
		containerd.WithPlatform(platforms.DefaultString()),
		containerd.WithPullUnpack,
		containerd.WithPullSnapshotter(containerd.DefaultSnapshotter),
	}
	image, err := client.Pull(ctx, ref, opts...)
	if err != nil {
		return nil, err
	}

	diffIDs, err := image.RootFS(ctx)
	if err != nil {
		return nil, err
	}

	snapshotKey := identity.ChainID(diffIDs).String()

	return &UnpackedImage{
		Image:       image,
		SnapshotKey: snapshotKey,
	}, nil
}
