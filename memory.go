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

	return m
}

// RemoveUntil remove code from front of memory an element.
func (m *Memory) RemoveUntil(e *Element) {
	current := m.Front
	for current != e {
		current.prevElement = nil
		current = current.nextElement
		// remove prev element values
		current.prevElement.nextElement = nil
		current.prevElement.list = nil
		current.prevElement = nil
	}

	m.Front = current
	m.Back = current
	m.Current = current
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
		e.list.Front = e.nextElement
		e.nextElement.Value = v

		e.list.Back = e.nextElement
	}

	return e.nextElement
}

// Prev generate new element or exist previous element.
func (e *Element) Prev(v interface{}) *Element {
	if e.prevElement == nil {
		e.prevElement = e.NewElement()
		e.prevElement.nextElement = e
		e.list.Back = e.prevElement
		e.prevElement.Value = v

		e.list.Front = e.nextElement
	}

	return e.prevElement
}
