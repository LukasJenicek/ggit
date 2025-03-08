package main

import (
	"strconv"
	"sync"
)

type Something struct {
	entries map[string]string
	mu      sync.RWMutex
}

func (s *Something) Add(key, value string) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.entries[key] = value
}

var something Something

func main() {
	i := 0
	something = Something{
		entries: make(map[string]string),
	}

	for ; i < 10; i++ {
		go func(i int) {
			something.Add(strconv.Itoa(i), strconv.Itoa(i))
		}(i)
	}
}
