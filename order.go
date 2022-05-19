package sly

// Cmp is a classical comparator func. See bounded methods for more.
type Cmp[T any] func(T, T) int

// Lt is `less than`.
func (c Cmp[T]) Lt(a T, b T) bool {
	return c(a, b) < 0
}

// Gt is `greater than`.
func (c Cmp[T]) Gt(a, b T) bool {
	return c(a, b) > 0
}

// Eq is `equal to`.
func (c Cmp[T]) Eq(a, b T) bool {
	return c(a, b) == 0
}

// Neq is `not equal to`.
func (c Cmp[T]) Neq(a, b T) bool {
	return c(a, b) != 0
}

// Leq is `less than or equal to`.
func (c Cmp[T]) Leq(a, b T) bool {
	return c(a, b) <= 0
}

// Geq is `greater than or equal to`.
func (c Cmp[T]) Geq(a, b T) bool {
	return c(a, b) >= 0
}

// Swap swaps two slice items by their indices.
func Swap[T any](s []T, i, j int) {
	s[i], s[j] = s[j], s[i]
}
