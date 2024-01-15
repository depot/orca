package chunker

import (
	"bufio"
	"encoding/binary"
	"io"
	"math/bits"
)

const (
	// DefaultMinChunkSize is the default minimum chunk size.
	DefaultMinChunkSize = 256 << 10 // 256 KiB

	// DefaultMaxChunkSize is the default maximum chunk size.
	DefaultMaxChunkSize = 8 << 20 // 8 MiB

	// DefaultNormalChunkSize is the default normal chunk size.
	DefaultNormalChunkSize = 256<<10 + (8 << 10) // 256 KiB + 8 KiB
)

// ChunkerOptions contains options for a chunker.
type ChunkerOptions struct {
	MinChunkSize    int
	MaxChunkSize    int
	NormalChunkSize int
}

// Chunker is an implementation of the UltraCDC algorithm.
type Chunker struct {
	options    *ChunkerOptions
	buf        *bufio.Reader
	splitpoint int
}

var _ io.WriterTo = (*Chunker)(nil)

// NewChunker creates a new chunker.
func NewChunker(reader io.Reader) *Chunker {
	return NewChunkerWithOptions(reader, nil)
}

// NewChunker creates a new chunker with options.
func NewChunkerWithOptions(reader io.Reader, options *ChunkerOptions) *Chunker {
	chunker := &Chunker{
		options: options,
	}

	if chunker.options == nil {
		chunker.options = &ChunkerOptions{}
	}

	if chunker.options.MinChunkSize == 0 {
		chunker.options.MinChunkSize = DefaultMinChunkSize
	}

	if chunker.options.MaxChunkSize == 0 {
		chunker.options.MaxChunkSize = DefaultMaxChunkSize
	}

	if chunker.options.NormalChunkSize == 0 {
		chunker.options.NormalChunkSize = DefaultNormalChunkSize
	}

	chunker.buf = bufio.NewReaderSize(reader, int(2*chunker.options.MaxChunkSize))

	return chunker
}

// Next returns the next chunk.
func (chunker *Chunker) Next() ([]byte, error) {
	if chunker.splitpoint != 0 {
		_, _ = chunker.buf.Discard(chunker.splitpoint)
		chunker.splitpoint = 0
	}

	data, err := chunker.buf.Peek(chunker.options.MaxChunkSize)
	if err != nil && err != io.EOF {
		return nil, err
	}

	length := len(data)
	if length == 0 {
		return nil, io.EOF
	}

	splitpoint := chunker.ultraCDC(data, length)
	chunker.splitpoint = splitpoint

	if splitpoint < chunker.options.MinChunkSize {
		return data[:splitpoint], io.EOF
	}

	return data[:splitpoint], nil
}

// WriteTo implements io.WriterTo, writing all chunks to the given writer.
func (chunker *Chunker) WriteTo(writer io.Writer) (int64, error) {
	var written int64

	for {
		data, err := chunker.Next()
		if err != nil && err != io.EOF {
			return written, err
		}

		length := len(data)

		if length == 0 {
			break
		}

		if _, err := writer.Write(data); err != nil {
			return written, err
		}

		written += int64(length)
	}

	return written, io.EOF
}

const (
	pattern                    uint64 = 0xAAAAAAAAAAAAAAAA
	maskSmaller                uint64 = 0x2F
	maskLarger                 uint64 = 0x2C
	lowEntropyStringsThreshold uint32 = 64
)

// An implementation of the UltraCDC algorithm.
// See https://ieeexplore.ieee.org/document/9894295.
func (chunker *Chunker) ultraCDC(data []byte, size int) int {
	normalChunkSize := chunker.options.NormalChunkSize

	if size <= chunker.options.MinChunkSize {
		return size
	}

	if size >= chunker.options.MaxChunkSize {
		size = chunker.options.MaxChunkSize
	} else if size <= normalChunkSize {
		normalChunkSize = size
	}

	i := chunker.options.MinChunkSize
	outWindow := binary.LittleEndian.Uint64(data[i:])
	distance := uint64(bits.OnesCount64(outWindow ^ pattern))
	i += 8

	mask := maskSmaller
	count := uint32(0)

	for i < size {
		if i == normalChunkSize {
			mask = maskLarger
		}

		inWindow := binary.LittleEndian.Uint64(data[i:])
		if outWindow^inWindow == 0 {
			count++
			if count == lowEntropyStringsThreshold {
				return i + 8
			}
		} else {
			count = 0
			for j := 0; j < 8; j++ {
				if (distance & mask) == 0 {
					return i + 8
				}
				inByte := data[i+j]
				outByte := data[i+j-8]
				distance = distance + uint64(hammingDistance[outByte][inByte])
			}
			outWindow = inWindow
		}

		i += 8
	}

	return size
}
