package virtualfilesystem

import "github.com/containerd/continuity"

type VirtualFilesystem struct {
	manifest continuity.Manifest
}

func NewVirtualFilesystem(manifest continuity.Manifest) *VirtualFilesystem {
	return &VirtualFilesystem{
		manifest: manifest,
	}
}
