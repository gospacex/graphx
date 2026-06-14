package tugraphx

import (
	"context"
	"sync"
	"testing"

	"github.com/gospacex/graphx"
)

func TestConcurrent_TuGraphGet(t *testing.T) {
	defer Reset()
	cfg := graphx.Config{Address: "127.0.0.1:17688"}
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
