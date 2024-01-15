package snapshots

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"time"

	"github.com/containerd/containerd/errdefs"
	"github.com/containerd/containerd/mount"
	"github.com/containerd/containerd/snapshots"
)

// Prepare will mount all layers of the parentSnapshotKey ont mountPoint.
// Once mounted, the caller can read and write files in the mountPoint.
//
// Use the returned "Committer" to unmount and commit your file changes.
//
// HINT: Mount points could be temporary directories if you want.
//
// NOTE: containerd uses the term "target" to refer to the mount point.
func Prepare(ctx context.Context, snapshotter snapshots.Snapshotter, parentSnapshotKey, mountPoint string) (*Committer, error) {
	// This key is ephemeral and exists only until the prepared snapshot is committed.
	activeSnapshotKey := uniqueKey()
	mounts, err := snapshotter.Prepare(ctx, activeSnapshotKey, parentSnapshotKey)
	if err != nil {
		if errdefs.IsAlreadyExists(err) {
			mounts, err = snapshotter.Mounts(ctx, activeSnapshotKey)
		}
		if err != nil {
			return nil, err
		}
	}

	if err := mount.All(mounts, mountPoint); err != nil {
		if err := snapshotter.Remove(ctx, activeSnapshotKey); err != nil && !errdefs.IsNotFound(err) {
			return nil, err
		}
		return nil, err
	}

	return &Committer{
		ActiveSnapshotKey: activeSnapshotKey,
		MountPoint:        mountPoint,
	}, nil
}

type Committer struct {
	ActiveSnapshotKey string
	MountPoint        string
}

func (c *Committer) Unmount() error {
	// TODO: buildkit uses syscall.MNT_DETACH.  Why? Seems dangerous, yes?
	flags := 0
	return mount.UnmountAll(c.MountPoint, flags)
}

func (c *Committer) Commit(ctx context.Context, snapshotter snapshots.Snapshotter, committedKey string) error {
	return snapshotter.Commit(ctx, committedKey, c.ActiveSnapshotKey)
}

func uniqueKey() string {
	t := time.Now()
	var b [3]byte
	// Ignore read failures, just decreases uniqueness
	_, _ = rand.Read(b[:])
	return fmt.Sprintf("%d-%s", t.Nanosecond(), base64.URLEncoding.EncodeToString(b[:]))
}
