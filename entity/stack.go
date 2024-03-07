package entity

type Stack[T any] struct {
	slice []T
}

func NewStack[T any]() *Stack[T] {
	return &Stack[T]{}
}

func (stack *Stack[T]) Push(key T) {
	stack.slice = append(stack.slice, key)
}

func (stack *Stack[T]) Peek() (T, bool) {
	var x T
	if len(stack.slice) > 0 {
		x = stack.slice[len(stack.slice)-1]
		return x, true
	}
	return x, false
}

func (stack *Stack[T]) Pop() (T, bool) {
	var x T
	if len(stack.slice) > 0 {
		x, stack.slice = stack.slice[len(stack.slice)-1], stack.slice[:len(stack.slice)-1]
		return x, true
	}
	return x, false
}

func (stack *Stack[T]) Length() int {
	return len(stack.slice)
}
