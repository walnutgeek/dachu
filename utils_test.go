package dachu_test

import (
	"testing"
	u "github.com/walnutgeek/dachu"
)

func testIt(t *testing.T, a *u.BytesCmp, b []byte, expectedEqual, expectedLess bool) {
	if a.EqualArray(b) != expectedEqual {
		t.Errorf("Equal failed expected:%v a:%v b:%v", expectedEqual, a, b)
	}
	if a.LessArray(b) != expectedLess {
		t.Errorf("Less failed expected:%v a:%v b:%v", expectedLess, a, b)
	}
	bb := u.NewBytesCmp(b)
	if a.Equal(&bb) != expectedEqual {
		t.Errorf("Equal failed expected:%v a:%v b:%v", expectedEqual, a, b)
	}
	if a.Less(&bb) != expectedLess {
		t.Errorf("Less failed expected:%v a:%v b:%v", expectedLess, a, b)
	}

}

func TestBytesCmpNil(t *testing.T) {
	var bb u.BytesCmp
	if bb.IsNil() != true {
		t.Errorf("Error:%#v expected to be nil", bb)
	}
	bb = u.NewBytesCmp([]byte{1, 2, 3})
	if bb.IsNil() != false {
		t.Errorf("Error:%#v expected to be not nil", bb)
	}
}

func TestBytesCmp(t *testing.T) {
	a := u.NewBytesCmp([]byte{1, 2, 3})
	testIt(t, &a, []byte{1, 2, 3}, true, false)
	testIt(t, &a, []byte{1, 2, 3, 4}, false, true)
	testIt(t, &a, []byte{1, 2, 4}, false, true)
	testIt(t, &a, []byte{1, 3}, false, true)
	testIt(t, &a, []byte{1, 2}, false, false)
	testIt(t, &a, []byte{1, 2, 1}, false, false)
	testIt(t, &a, []byte{1, 1}, false, false)
}

func TestAppendSize(t *testing.T) {
	testSize(0, []byte{0x80}, t)
	testSize(255, []byte{0x7f, 0x81}, t)
	testSize(u.MaxSize, []byte{0x7f, 0x7f, 0x7f, 0x7f, 0x8f}, t)
	testParseSizeError([]byte{0x7f, 0x7f, 0x7f, 0x7f, 0x9f}, "Overflow: 10000 i:4 0x9f", t)
	testParseSizeError([]byte{0x7f, 0x7f, 0x7f, 0x7f}, "Unexpected end of input", t)
	testParseSizeError([]byte{0x7f, 0x7f, 0x7f, 0x7f, 0x0f, 0x0f}, "Size is too big", t)
}

func testParseSizeError(toParse []byte, expectedError string, t *testing.T) {
	n, err := u.ParseSize(&u.SliceIterator{Src: toParse})
	if err == nil {
		t.Errorf("No error: n=%d", n)
	} else if err.Error() != expectedError {
		t.Errorf("Error:%#v is not what expected: %#v", err.Error(), expectedError)
	}
}

func testSize(n u.Size, expected []byte, t *testing.T) {
	bc_expected := u.NewBytesCmp(expected)
	a := u.AppendSize([]byte{}, n)
	if !(&bc_expected).EqualArray(a) {
		t.Errorf("AppendSize failed  %v %#v", expected, a)
	}
	n2, err := u.ParseSize(&u.SliceIterator{Src: a})
	if err != nil {
		t.Error(err)
	}
	if n2 != n {
		t.Errorf("ParseSize failed %b %b", n, n2)
	}
}
