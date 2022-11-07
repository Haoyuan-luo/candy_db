package common

type Element struct {
	Data   *node
	Levels []*Element
}

func NewElement(oriNode *node, level int) *Element {
	return &Element{Data: oriNode, Levels: make([]*Element, level)}
}

func (e *Element) AddElement(elem *Element, level int) {
	for i := 0; i < level; i++ {
		elem.Levels[i] = e.Levels[i]
		e.Levels[i] = elem
	}
}
