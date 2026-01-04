package domain

type Id[T any] struct {
	Value int64
}

func Equal[T any](id1, id2 Id[T]) bool {
	return id1.Value == id2.Value
}

func NewId[T any](value int64) Id[T] {
	return Id[T]{Value: value}
}
