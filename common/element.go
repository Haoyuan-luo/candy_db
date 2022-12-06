package common

type Element struct {
	Data  *node
	Level []*Element
}

func NewElement(oriNode *node, level int) *Element {
	return &Element{Data: oriNode, Level: make([]*Element, level)}
}
