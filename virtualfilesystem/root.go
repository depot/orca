package virtualfilesystem

import (
	"context"
	"fmt"
	"os"
	"path"
	"path/filepath"
	"strings"
	"syscall"

	"github.com/depot/orca/layers/entry"
	layerv1 "github.com/depot/orca/proto/depot/orca/layer/v1"
	"github.com/hanwen/go-fuse/v2/fs"
	"github.com/hanwen/go-fuse/v2/fuse"
)

type virtualFilesystemRoot struct {
	fs.Inode

	layer layerv1.LayerEntries
}

func (r *virtualFilesystemRoot) OnAdd(ctx context.Context) {
	for _, e := range r.layer.Entries {
		inodeNumber := uint64(0)

		attr := entry.EntryToFuseAttr(e)
		inodeNumber += 1
		stableAttr := fs.StableAttr{Mode: attr.Mode, Ino: inodeNumber}
		mode := os.FileMode(e.Mode)

		var node fs.InodeEmbedder
		switch {
		case mode.IsRegular():
			chunks := make([]chunk, len(e.Blocks))
			for idx, b := range e.Blocks {
				chunks[idx] = chunk{
					size:     b.SizeBytes,
					location: digestToLocation(b.Digest),
				}
			}
			node = &virtualFile{size: e.SizeBytes, chunks: chunks, Attr: *attr}

		case mode&os.ModeSymlink != 0:
			node = &fs.MemRegularFile{Data: []byte(e.Target), Attr: *attr}

		default:
			node = &fs.MemRegularFile{Attr: *attr}
		}

		for _, path := range e.Path {
			dir, base := filepath.Split(path)
			parent := r.EmbeddedInode()

			for _, component := range strings.Split(dir, "/") {
				if component == "" {
					continue
				}
				child := parent.GetChild(component)
				if child == nil {
					child = r.NewPersistentInode(ctx, &fs.Inode{}, fs.StableAttr{Mode: fuse.S_IFDIR})
					parent.AddChild(component, child, false)
				}
				parent = child
			}

			inode := parent.NewPersistentInode(ctx, node, stableAttr)
			parent.AddChild(base, inode, false)
		}
	}

	ch := r.NewPersistentInode(
		ctx, &fs.MemRegularFile{
			Data: []byte("file.txt"),
			Attr: fuse.Attr{
				Mode: 0644,
			},
		}, fs.StableAttr{Ino: 2})
	r.AddChild("file.txt", ch, false)

	r.AddChild("hello.txt", r.NewPersistentInode(
		ctx, &virtualFile{
			size: 12,
			chunks: []chunk{
				{size: 6, location: "/mnt/hello"},
				{size: 6, location: "/mnt/world"},
			},
			Attr: fuse.Attr{
				Mode: 0644,
			},
		}, fs.StableAttr{Ino: 3}), false)

	dir := &fs.MemRegularFile{
		Attr: fuse.Attr{
			Mode: 0755,
		},
	}
	r.AddChild("hello", r.NewPersistentInode(
		ctx, dir, fs.StableAttr{Ino: 4, Mode: fuse.S_IFDIR}), false)

	dir.AddChild("hello.txt", r.NewPersistentInode(
		ctx, &virtualFile{
			size: 12,
			chunks: []chunk{
				{size: 6, location: "/mnt/hello"},
				{size: 6, location: "/mnt/world"},
			},
			Attr: fuse.Attr{
				Mode: 0644,
			},
		}, fs.StableAttr{Ino: 4}), false)
}

func (r *virtualFilesystemRoot) Getattr(ctx context.Context, fh fs.FileHandle, out *fuse.AttrOut) syscall.Errno {
	out.Mode = 0755
	return 0
}

func digestToLocation(digest *layerv1.Digest) string {
	// TODO: add some kind of block store configuration
	return path.Join("/tmp/blocks", digest.Algorithm.String(), fmt.Sprintf("%x", digest.Sum))
}
