package graphx

import (
	"context"
	"sync"
	"testing"
)

func TestSingleton_ConcurrentGet(t *testing.T) {
	var calls int32
	s := NewSingleton[int]("test")

	var wg sync.WaitGroup
	for range 100 {
		wg.Add(1)
		go func() {
			defer wg.Done()
			val, err := s.Get(func() (*int, error) {
				v := 42
				return &v, nil
			})
			if err != nil {
				t.Error(err)
			}
			if *val != 42 {
				t.Errorf("got %d, want 42", *val)
			}
		}()
	}
	wg.Wait()

	_ = calls
}

func TestSingleton_ConcurrentGetFailed(t *testing.T) {
	// All concurrent callers should see the same error (not cached, retried per caller)
	s := NewSingleton[int]("fail")

	var wg sync.WaitGroup
	errs := make(chan error, 50)
	for range 50 {
		wg.Add(1)
		go func() {
			defer wg.Done()
			_, err := s.Get(func() (*int, error) {
				return nil, context.Canceled
			})
			errs <- err
		}()
	}
	wg.Wait()
	close(errs)

	for err := range errs {
		if err == nil {
			t.Fatal("expected error, got nil")
		}
	}
}

func TestSingleton_ResetUnderGet(t *testing.T) {
	// Reset should not deadlock with concurrent Get
	s := NewSingleton[int]("concurrent")

	done := make(chan struct{})
	go func() {
		for {
			select {
			case <-done:
				return
			default:
				s.Get(func() (*int, error) { v := 99; return &v, nil })
			}
		}
	}()

	for range 100 {
		s.Reset()
	}
	close(done)
}
