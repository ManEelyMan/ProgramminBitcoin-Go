package collections_test

import (
	"bitcoin-go/collections"
	"bitcoin-go/utility"
	"bytes"
	"testing"
)

func TestStackEmptyPop(t *testing.T) {
	s := collections.NewStack()

	i, ok := s.Pop()
	if ok || i != nil {
		t.Error()
	}
}
func TestStackRoundTrip(t *testing.T) {
	s := collections.NewStack()

	rand := utility.RandomData(32)

	s.Push(rand)

	res, ok := s.Pop()
	if !ok {
		t.Error()
	}

	if !bytes.Equal(rand, res) {
		t.Error()
	}
}

func TestStackMultipleElements(t *testing.T) {
	s := collections.NewStack()

	one := utility.RandomData(32)
	two := utility.RandomData(32)
	three := utility.RandomData(32)

	s.Push(one)
	s.Push(two)
	s.Push(three)

	i, ok := s.Pop()
	if !ok || i == nil {
		t.Error()
	}

	if !bytes.Equal(three, i) {
		t.Error()
	}

	i, ok = s.Pop()
	if !ok || i == nil {
		t.Error()
	}

	if !bytes.Equal(two, i) {
		t.Error()
	}

	i, ok = s.Pop()
	if !ok || i == nil {
		t.Error()
	}

	if !bytes.Equal(one, i) {
		t.Error()
	}
}

func TestStackPeek(t *testing.T) {
	s := collections.NewStack()

	i, ok := s.Peek()
	if ok || i != nil {
		t.Error()
	}

	one := utility.RandomData(4)
	two := utility.RandomData(4)

	s.Push(one)
	s.Push(two)

	i, ok = s.Peek()
	if !ok {
		t.Error()
	}

	if !bytes.Equal(two, i) {
		t.Error()
	}

	i, ok = s.PeekAt(1)
	if !ok {
		t.Error()
	}

	if !bytes.Equal(one, i) {
		t.Error()
	}

	i, ok = s.PeekAt(2)
	if ok || i != nil {
		t.Error()
	}
}
