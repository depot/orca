package entry

import (
	"os"

	layerv1 "github.com/depot/orca/proto/depot/orca/layer/v1"
	"github.com/hanwen/go-fuse/v2/fuse"
)

func EntryToFuseAttr(e *layerv1.LayerEntry) *fuse.Attr {
	attr := &fuse.Attr{
		Mode: e.Mode,
		Owner: fuse.Owner{
			Uid: uint32(e.Uid),
			Gid: uint32(e.Gid),
		},
		Mtime: uint64(e.Mtime.AsTime().Unix()),
	}

	mode := os.FileMode(e.Mode)
	switch {
	case mode.IsRegular():
		attr.Size = e.SizeBytes

	case mode&os.ModeDevice != 0:
		attr.Rdev = uint32(e.Major)<<8 | uint32(e.Minor)
	}

	return attr
}
