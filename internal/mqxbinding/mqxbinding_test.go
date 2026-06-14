package mqxbinding

import (
	"sync"
	"testing"
)

// TestReset_ClearsRedisCache verifies that Reset() empties the Redis client
// cache, so a subsequent Redis() call with the same key creates a fresh
// *redis.Client instance (regression test for F-2: cache growth).
func TestReset_ClearsRedisCache(t *testing.T) {
	Reset()
	t.Cleanup(Reset)

	cfg := MQConfig{
		Driver: "redis",
		Mode:    "single",
		Addrs:   []string{"127.0.0.1:6379"},
	}

	c1, err := Redis(cfg)
	if err != nil {
		t.Fatal(err)
	}

	c2, err := Redis(cfg)
	if err != nil {
		t.Fatal(err)
	}
	if c1 != c2 {
		t.Fatal("expected cached client on second call")
	}

	Reset()

	c3, err := Redis(cfg)
	if err != nil {
		t.Fatal(err)
	}
	if c1 == c3 {
		t.Fatal("Reset() should have cleared cache; expected new client")
	}
}

// TestReset_ClearsKafkaCache verifies the same behavior for Kafka producers.
func TestReset_ClearsKafkaCache(t *testing.T) {
	Reset()
	t.Cleanup(Reset)

	cfg := MQConfig{
		Driver: "kafka",
		Mode:    "single",
		Addrs:   []string{"127.0.0.1:9092"},
	}

	p1, err := Kafka(cfg)
	if err != nil {
		t.Fatal(err)
	}
	p2, err := Kafka(cfg)
	if err != nil {
		t.Fatal(err)
	}
	if p1 != p2 {
		t.Fatal("expected cached producer on second call")
	}

	Reset()

	p3, err := Kafka(cfg)
	if err != nil {
		t.Fatal(err)
	}
	if p1 == p3 {
		t.Fatal("Reset() should have cleared cache; expected new producer")
	}
}

// TestReset_ConcurrentSafe verifies that Reset() and Redis()/Kafka() can run
// concurrently without data race or panic. Run under `go test -race`.
func TestReset_ConcurrentSafe(t *testing.T) {
	Reset()
	t.Cleanup(Reset)

	cfg := MQConfig{
		Driver: "redis",
		Mode:    "single",
		Addrs:   []string{"127.0.0.1:6379"},
	}

	const goroutines = 8
	const iterations = 50

	var wg sync.WaitGroup
	wg.Add(goroutines)
	for i := 0; i < goroutines; i++ {
		go func() {
			defer wg.Done()
			for j := 0; j < iterations; j++ {
				_, _ = Redis(cfg)
			}
		}()
	}
	// Hammer Reset() concurrently
	for j := 0; j < iterations; j++ {
		Reset()
	}
	wg.Wait()
}

// TestReset_Idempotent verifies Reset() can be called multiple times safely.
func TestReset_Idempotent(t *testing.T) {
	Reset()
	Reset()
	Reset()
}
