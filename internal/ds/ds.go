package ds

import (
	"fmt"
	"sort"
)

type Set[T comparable] map[T]struct{}

func NewSet[T comparable](values []T) Set[T] {
	s := Set[T]{}
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

func (s Set[T]) SortedValues() []T {
	keys := make([]T, 0, len(s))

	for value := range s {
		keys = append(keys, value)
	}

	sort.Slice(keys, func(i, j int) bool {
		return fmt.Sprint(keys[i]) < fmt.Sprint(keys[j])
	})

	return keys
}
