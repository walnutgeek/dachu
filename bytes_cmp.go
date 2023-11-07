package dachu

import "golang.org/x/exp/constraints"

func Min[T constraints.Ordered](a, b T) T {
	if a < b {
		return a
	}
	return b
}

func Max[T constraints.Ordered](a, b T) T {
	if a < b {
		return b
	}
	return a
}

type BytesCmp struct {
	src []byte
}

func NewBytesCmp(src []byte) BytesCmp {
	return BytesCmp{src: src}
}

func (b BytesCmp) IsNil() bool {
	return b.src == nil
}

func (b *BytesCmp) Equal(dst *BytesCmp) bool {
	return b.EqualArray(dst.src)
}

func (b *BytesCmp) EqualArray(dst []byte) bool {
	if len(b.src) != len(dst) {
		return false
	}
	for i, v := range b.src {
		if v != dst[i] {
			return false
		}
	}
	return true
}

func (b *BytesCmp) Less(dst *BytesCmp) bool {
	return b.LessArray(dst.src)
}

func (b *BytesCmp) LessArray(dst []byte) bool {
	minLength := Min(len(b.src), len(dst))
	for i := 0; i < minLength; i++ {
		if b.src[i] < dst[i] {
			return true
		} else if b.src[i] > dst[i] {
			return false
		}
	}
	return len(b.src) < len(dst) // if all bytes are equal, the shorter one is less
}
