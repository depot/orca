package chunk

import (
	"encoding"
	"encoding/binary"
	"encoding/hex"
	"fmt"
)

// Represents the size of the Chunk struct in bytes.
const chunkSize = 12

var ErrInvalidChunkSize = fmt.Errorf("invalid chunk size")

type Chunk struct {
	Size uint32
	Hash uint64
}

var _ fmt.Stringer = (*Chunk)(nil)
var _ encoding.BinaryMarshaler = (*Chunk)(nil)
var _ encoding.BinaryUnmarshaler = (*Chunk)(nil)

func (c *Chunk) MarshalBinary() ([]byte, error) {
	buf := make([]byte, chunkSize)
	binary.BigEndian.PutUint32(buf, c.Size)
	binary.BigEndian.PutUint64(buf[4:], c.Hash)
	return buf, nil
}

func (c *Chunk) UnmarshalBinary(buf []byte) error {
	if len(buf) != chunkSize {
		return ErrInvalidChunkSize
	}

	c.Size = binary.BigEndian.Uint32(buf)
	c.Hash = binary.BigEndian.Uint64(buf[4:])

	return nil
}

func (c *Chunk) String() string {
	bytes, _ := c.MarshalBinary()
	return hex.EncodeToString(bytes)
}

func NewChunkFromString(s string) (*Chunk, error) {
	c := &Chunk{}
	bytes, err := hex.DecodeString(s)
	if err != nil {
		return nil, err
	}
	err = c.UnmarshalBinary(bytes)
	if err != nil {
		return nil, err
	}
	return c, nil
}
