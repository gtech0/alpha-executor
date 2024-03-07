package entity

type Queue[T any] struct {
	slice []T
}

func NewQueue[T any]() *Queue[T] {
	return &Queue[T]{}
}

func (stack *Queue[T]) Enqueue(key T) {
	stack.slice = append(stack.slice, key)
}

func (stack *Queue[T]) Dequeue() (T, bool) {
	var x T
	if len(stack.slice) > 0 {
		x, stack.slice = stack.slice[0], stack.slice[1:]
		return x, true
	}
	return x, false
}

func (stack *Queue[T]) Length() int {
	return len(stack.slice)
}

func (stack *Queue[T]) Slice() []T {
	return stack.slice
}
