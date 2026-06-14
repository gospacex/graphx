package graphx

import (
	"context"
	"fmt"
	"sync"
	"sync/atomic"
)

// Closer closes the underlying driver connection.
type Closer interface {
	Close(ctx context.Context) error
}

// Pinger checks driver health.
type Pinger interface {
	Ping(ctx context.Context) error
}

// Singleton manages a lazy singleton instance with atomic double-checked locking.
// A failed factory is not cached so the next call retries.
type Singleton[T any] struct {
	mu       sync.Mutex
	instance atomic.Pointer[T]
	name     string
}

// NewSingleton creates a Singleton. name is used only for error wrapping.
func NewSingleton[T any](name string) *Singleton[T] {
	return &Singleton[T]{name: name}
}

// Get returns the cached instance, calling factory on the first invocation.
func (s *Singleton[T]) Get(factory func() (*T, error)) (*T, error) {
	if inst := s.instance.Load(); inst != nil {
		return inst, nil
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	if inst := s.instance.Load(); inst != nil {
		return inst, nil
	}
	inst, err := factory()
	if err != nil {
		return nil, fmt.Errorf("graphx/%s: %w", s.name, err)
	}
	s.instance.Store(inst)
	return inst, nil
}

// Reset clears the cached instance. Intended for testing only.
func (s *Singleton[T]) Reset() {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.instance.Store(nil)
}
