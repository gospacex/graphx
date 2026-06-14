package tigergraphx

import (
	"context"
	"sync"
	"testing"

	"github.com/gospacex/graphx"
)

func TestConcurrent_TigerGet(t *testing.T) {
	defer Reset()
	cfg := graphx.Config{Address: "127.0.0.1:14240"}
	var wg sync.WaitGroup
	for range 20 {
		wg.Add(1)
		go func() {
			defer wg.Done()
			_, _ = Get(context.Background(), cfg)
		}()
	}
	wg.Wait()
}
