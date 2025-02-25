package ds

import (
	"sort"
)

// Set
// TODO: Implement other funcs to support interface for Set struct
// Contains(value T) bool - to check if an element exists in the set
// Remove(value T) - to remove an element from the set
// Union(other Set[T]) Set[T] - to create a union of two sets
// Intersection(other Set[T]) Set[T] - to create an intersection of two sets
type Set[T comparable] map[T]struct{}

func NewSet[T comparable](values []T) Set[T] {
	s := Set[T]{}

	if values == nil {
		return s
	}

	for _, value := range values {
		s.Add(value)
	}

	return s
}

func (s Set[T]) Add(value T) {
	s[value] = struct{}{}
}

func (s Set[T]) Size() int {
	return len(s)
}

func (s Set[T]) SortedValues(less func(a, b T) bool) []T {
	keys := make([]T, 0, len(s))

	for value := range s {
		keys = append(keys, value)
	}

	sort.Slice(keys, func(i, j int) bool {
		return less(keys[i], keys[j])
	})

	return keys
}
