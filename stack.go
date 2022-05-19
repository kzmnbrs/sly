package sly

// StackPop pops the top value from the stack.
func StackPop[T any](stackPtr *[]T) T {
	stack := *stackPtr

	val := stack[len(stack)-1]
	stack = stack[:len(stack)-1]

	*stackPtr = stack
	return val
}
