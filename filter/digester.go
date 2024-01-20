package filter

import (
	"io"

	"github.com/cespare/xxhash"
	"github.com/containerd/continuity"
	"github.com/depot/orca/util/chunk"
	"github.com/depot/orca/util/chunker"
	"github.com/opencontainers/go-digest"
)

var _ continuity.Digester = (*Digester)(nil)

type Digester struct{}

func (d *Digester) Digest(r io.Reader) (digest.Digest, error) {
	var chunks chunk.Chunks
	chunker := chunker.NewChunker(r)

	for {
		data, err := chunker.Next()
		if err != nil && err != io.EOF {
			return "", err
		}

		hash := xxhash.Sum64(data)
		chunk := chunk.NewChunk(len(data), hash)
		chunks = append(chunks, chunk)
		if err == io.EOF {
			break
		}
	}

	return digest.Digest(chunks.String()), nil
}
