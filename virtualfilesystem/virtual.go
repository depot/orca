package virtualfilesystem

import (
	"context"
	"fmt"
	"os"
	"sync"
	"syscall"

	"github.com/hanwen/go-fuse/v2/fs"
	"github.com/hanwen/go-fuse/v2/fuse"
)

type virtualFile struct {
	fs.Inode

	Attr fuse.Attr

	size   uint64
	chunks []chunk
	mu     sync.Mutex
}

type chunk struct {
	size     uint64
	location string
}

var _ = (fs.NodeOpener)((*virtualFile)(nil))
var _ = (fs.NodeReader)((*virtualFile)(nil))
var _ = (fs.NodeGetattrer)((*virtualFile)(nil))

func (f *virtualFile) Open(ctx context.Context, flags uint32) (fh fs.FileHandle, fuseFlags uint32, errno syscall.Errno) {
	return nil, fuse.FOPEN_KEEP_CACHE, fs.OK
}

func (f *virtualFile) Read(ctx context.Context, fh fs.FileHandle, dest []byte, off int64) (fuse.ReadResult, syscall.Errno) {
	f.mu.Lock()
	defer f.mu.Unlock()

	offset := uint64(off)

	end := offset + uint64(len(dest))
	if end > f.size {
		end = f.size
	}
	targetFilledBytes := end - offset

	fmt.Printf("Reading from %d to %d\n", offset, end)

	// Determine starting source
	var sourceIdx int
	var sourceOffset uint64
	for idx, s := range f.chunks {
		if offset < s.size {
			sourceIdx = idx
			sourceOffset = offset
			break
		}
		offset -= s.size
	}

	filledBytes := uint64(0)
	for i := sourceIdx; i < len(f.chunks); i++ {
		fmt.Printf("Reading chunk %d, filledBytes %d, len %d\n", i, filledBytes, len(dest))

		if filledBytes >= targetFilledBytes {
			break
		}

		chunk := f.chunks[i]
		file, err := os.Open(chunk.location)
		if err != nil {
			panic(err)
		}

		// If this is the first chunk, we need to seek to the offset
		if i == sourceIdx && sourceOffset > 0 {
			_, err = file.Seek(int64(sourceOffset), 0)
			if err != nil {
				panic(err)
			}
		}

		// Read the data
		bytesToRead := uint64(len(dest)) - filledBytes
		if bytesToRead > chunk.size {
			bytesToRead = chunk.size
		}

		bytesRead, err := file.Read(dest[filledBytes : filledBytes+bytesToRead])
		if err != nil {
			panic(err)
		}

		filledBytes += uint64(bytesRead)

		// Close the file
		err = file.Close()
		if err != nil {
			panic(err)
		}
	}

	return fuse.ReadResultData(dest), fs.OK
}

func (f *virtualFile) Getattr(ctx context.Context, fh fs.FileHandle, out *fuse.AttrOut) syscall.Errno {
	f.mu.Lock()
	defer f.mu.Unlock()
	out.Attr = f.Attr
	out.Attr.Size = uint64(f.size)
	return fs.OK
}
