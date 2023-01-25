package util

import "encoding/json"

type SerializeService[T any] interface {
	Serialize(input T) (ret []byte, err error)
	Deserialize(input []byte) (ret T, err error)
}

type serialize[T any] struct{}

func NewSerializeService[T any]() SerializeService[T] {
	return &serialize[T]{}
}

func (s serialize[T]) Serialize(input T) ([]byte, error) {
	return json.Marshal(input)
}

func (s serialize[T]) Deserialize(bytes []byte) (ret T, err error) {
	err = json.Unmarshal(bytes, &ret)
	return
}
