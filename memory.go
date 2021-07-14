package brainit

// Memory is main struct of linked list.
type Memory struct {
	Len     uint
	Front   *Element
	Back    *Element
	Current *Element
}

// Init help easily to get memory struct.
func (m *Memory) Init(v interface{}) *Memory {
	element := &Element{
		list: m,
	}
	m.Front = element
	m.Back = element
	m.Current = element
	m.Current.Value = v

	m.Len = 1

	return m
}

// RemoveUntil remove code from front of memory an element.
// if element not inside of memory, nothing change.
func (m *Memory) RemoveUntil(e *Element) {
	current := m.Front
	count := 0
	for current != nil && current != e {
		count++
		current = current.nextElement
	}

	if current != nil && count > 0 {
		current.prevElement = nil
		m.Len = m.Len - uint(count)
		m.Front = current
		m.Current = current
	}
}

// NewMemory is a helper function to get a functional memory.
func NewMemory(v interface{}) *Memory {
	return new(Memory).Init(v)
}

// Element is an struct of duble-linked list.
type Element struct {
	nextElement *Element
	prevElement *Element

	// Know about belong memlist
	list  *Memory
	Value interface{}
}

// NewElement generate new element and set size of memory.
func (e *Element) NewElement() *Element {
	e.list.Len++
	return &Element{
		list: e.list,
	}
}

// Next generate new element or exist next element.
func (e *Element) Next(v interface{}) *Element {
	if e.nextElement == nil {
		e.nextElement = e.NewElement()
		e.nextElement.prevElement = e
		e.nextElement.Value = v

		e.list.Back = e.nextElement
	}

	return e.nextElement
}

// Prev generate new element with argument or return exist previous element.
func (e *Element) Prev(v interface{}) *Element {
	if e.prevElement == nil {
		e.prevElement = e.NewElement()
		e.prevElement.nextElement = e
		e.prevElement.Value = v

		e.list.Front = e.prevElement
	}

	return e.prevElement
}
