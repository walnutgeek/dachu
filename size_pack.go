package dachu

import (
	"fmt"
)

type Size uint32

const MaxSize Size = ^Size(0)
const shift int = 7
const end_mask byte = 1 << shift
const meat_mask byte = ^end_mask
const max_size_bytes = 5
const should_be_empty_on_last_mask byte = 0b01110000

func AppendSize(dst []byte, size Size) []byte {
	for {
		current := byte(size) & meat_mask
		size >>= shift
		if size == 0 {
			return append(dst, current|end_mask)
		} else {
			dst = append(dst, current)
		}
	}
}

// reads source byte by byte. Returns next byte and `true` if byte is available, 
// otherwise `(byte)0` and `false`
type ByteIterator interface {
	Next() (byte, bool)
}

type SliceIterator struct {
	Src  []byte
	Position int
}

func (it *SliceIterator) Next() (byte, bool) {
	if it.Position >= len(it.Src) {
		return 0, false
	}
	it.Position++
	return it.Src[it.Position-1], true
}

func ParseSize(src ByteIterator) (Size, error) {
	var size Size
	for i := 0; i < max_size_bytes; i++ {
		current, ok := src.Next()
		if !ok {
			return 0, fmt.Errorf("Unexpected end of input")
		}
		if i == max_size_bytes-1 {
			if x := current & should_be_empty_on_last_mask; x != 0 {
				return 0, fmt.Errorf("Overflow: %b i:%d 0x%02x", x, i, current)
			}
		}
		size |= Size(current&meat_mask) << (shift * i)
		if current&end_mask != 0 {
			return size, nil
		}
	}
	return 0, fmt.Errorf("Size is too big")
}
