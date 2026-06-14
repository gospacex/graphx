package neo4jx

import (
	"context"
	"sync"
	"testing"

	"github.com/gospacex/graphx"
)

func TestConcurrent_Neo4jGet(t *testing.T) {
	defer Reset()
	ctx := context.Background()

	cfg := graphx.Config{Address: "127.0.0.1:17687"}
	var wg sync.WaitGroup
	for range 20 {
		wg.Add(1)
		go func() {
			defer wg.Done()
			_, err := Get(ctx, cfg)
			if err == nil {
				// all but first should see cached error — any result is fine
			}
		}()
	}
	wg.Wait()
}

func TestConcurrent_Neo4jSingletonRace(t *testing.T) {
	defer Reset()

	var wg sync.WaitGroup
	wg.Add(2)

	go func() {
		defer wg.Done()
		_, _ = Get(context.Background(), graphx.Config{Name: "race", Address: "127.0.0.1:17687"})
	}()
	go func() {
		defer wg.Done()
		Reset()
	}()
	wg.Wait()
}
