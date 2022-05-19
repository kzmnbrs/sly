package sly

import "golang.org/x/exp/constraints"

// Compare is a comparator func. See bounded methods for more.
type Compare[T any] func(T, T) int

// Less is `less than`.
func (c Compare[T]) Less(a T, b T) bool {
	return c(a, b) < 0
}

// Greater is `greater than`.
func (c Compare[T]) Greater(a, b T) bool {
	return c(a, b) > 0
}

// Equal is `equal to`.
func (c Compare[T]) Equal(a, b T) bool {
	return c(a, b) == 0
}

// NotEqual is `not equal to`.
func (c Compare[T]) NotEqual(a, b T) bool {
	return c(a, b) != 0
}

// LessOrEqual is `less than or equal to`.
func (c Compare[T]) LessOrEqual(a, b T) bool {
	return c(a, b) <= 0
}

// GreaterOrEqual is `greater than or equal to`.
func (c Compare[T]) GreaterOrEqual(a, b T) bool {
	return c(a, b) >= 0
}

// CompareOrdered is a built-in comparator for ordered types.
func CompareOrdered[T constraints.Ordered](a, b T) int {
	switch {
	case a < b:
		return -1
	case a > b:
		return 1
	}
	return 0
}

// CompareReverse returns the negated comparator.
func CompareReverse[T any](compare Compare[T]) Compare[T] {
	return func(a T, b T) int {
		return compare(a, b) * -1
	}
}
