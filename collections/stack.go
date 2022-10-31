package collections

type StackElement []byte

type Stack struct {
	data []StackElement
}

func NewStack() Stack {
	s := Stack{}
	s.data = make([]StackElement, 0)
	return s
}

func (s *Stack) Push(item []byte) {
	s.data = append(s.data, item)
}

func (s *Stack) Pop() ([]byte, bool) {

	if len(s.data) == 0 {
		return nil, false
	}

	item := s.data[len(s.data)-1]
	s.data = s.data[0 : len(s.data)-1]
	return item, true
}

func (s *Stack) Peek() ([]byte, bool) {
	var item []byte = nil

	if s.Length() > 0 {
		item = s.data[len(s.data)-1]
	} else {
		return nil, false
	}

	return item, true
}

func (s *Stack) Length() uint32 {
	return (uint32)(len(s.data))
}

func (s *Stack) PeekAt(index uint32) ([]byte, bool) {
	if index >= s.Length() {
		return nil, false
	}

	index = (uint32)(len(s.data)) - index - 1
	return s.data[index], true
}
