package chunk

import (
	"encoding"
	"encoding/hex"
	"fmt"
)

type Chunks []Chunk

var _ fmt.Stringer = (*Chunk)(nil)
var _ encoding.BinaryMarshaler = (*Chunk)(nil)
var _ encoding.BinaryUnmarshaler = (*Chunk)(nil)

func (c Chunks) MarshalBinary() ([]byte, error) {
	buf := make([]byte, chunkSize*len(c))
	for i, chunk := range c {
		bytes, err := chunk.MarshalBinary()
		if err != nil {
			return nil, err
		}
		copy(buf[i*chunkSize:], bytes)
	}
	return buf, nil
}

func (c *Chunks) UnmarshalBinary(buf []byte) error {
	if len(buf)%chunkSize != 0 {
		return ErrInvalidChunkSize
	}

	chunks := make(Chunks, 0, len(buf)/chunkSize)

	for i := 0; i < len(buf); i += chunkSize {
		chunk := Chunk{}
		err := chunk.UnmarshalBinary(buf[i : i+chunkSize])
		if err != nil {
			return err
		}
		chunks = append(chunks, chunk)
	}

	*c = chunks
	return nil
}

func (c Chunks) String() string {
	bytes, _ := c.MarshalBinary()
	return hex.EncodeToString(bytes)
}

func NewChunksFromString(s string) (Chunks, error) {
	c := Chunks{}
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
