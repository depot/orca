package filter

import (
	"io"

	"github.com/containerd/continuity"
	"github.com/opencontainers/go-digest"
)

var _ continuity.Digester = (*Digester)(nil)

type Digester struct{}

func (d *Digester) Digest(r io.Reader) (digest.Digest, error) {
	return "", nil
}
