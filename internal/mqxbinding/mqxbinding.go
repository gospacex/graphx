package mqxbinding

import (
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
	"github.com/redis/go-redis/v9"
)

var (
	cacheMu sync.RWMutex
	cache   = map[string]any{}
)

func cacheKey(driver, mode string, addrs []string) string {
	return driver + ":" + mode + ":" + strings.Join(addrs, ",")
}

// Reset clears all cached Kafka and Redis clients. Intended for graceful
// shutdown: call before process exit to release file descriptors and
// network connections. Safe to call concurrently with Redis()/Kafka().
// Safe to call when the cache is empty.
func Reset() {
	cacheMu.Lock()
	cache = map[string]any{}
	cacheMu.Unlock()
}

func Redis(cfg MQConfig) (*redis.Client, error) {
	key := cacheKey(cfg.Driver, cfg.Mode, cfg.Addrs)
	cacheMu.RLock()
	if existing, ok := cache[key]; ok {
		cacheMu.RUnlock()
		return existing.(*redis.Client), nil
	}
	cacheMu.RUnlock()

	cacheMu.Lock()
	defer cacheMu.Unlock()
	if existing, ok := cache[key]; ok {
		return existing.(*redis.Client), nil
	}

	opts := &redis.Options{
		Addr:         cfg.Addrs[0],
		Password:     cfg.Auth.Password,
		DB:           cfg.Redis.DB,
		PoolSize:     cfg.Redis.PoolSize,
		DialTimeout:  5 * time.Second,
		ReadTimeout:  3 * time.Second,
		WriteTimeout: 3 * time.Second,
	}
	if cfg.Redis.PoolSize <= 0 {
		opts.PoolSize = 10
	}
	if cfg.Mode == "cluster" {
		// For cluster mode, return client pointed at first addr; real cluster uses ClusterClient
		opts.Addr = cfg.Addrs[0]
	}

	client := redis.NewClient(opts)
	cache[key] = client
	return client, nil
}

func Kafka(cfg MQConfig) (*kafka.Producer, error) {
	key := cacheKey(cfg.Driver, cfg.Mode, cfg.Addrs)
	cacheMu.RLock()
	if existing, ok := cache[key]; ok {
		cacheMu.RUnlock()
		return existing.(*kafka.Producer), nil
	}
	cacheMu.RUnlock()

	cacheMu.Lock()
	defer cacheMu.Unlock()
	if existing, ok := cache[key]; ok {
		return existing.(*kafka.Producer), nil
	}

	conf := &kafka.ConfigMap{
		"bootstrap.servers": strings.Join(cfg.Addrs, ","),
	}
	if cfg.Kafka.SecurityProtocol != "" {
		conf.SetKey("security.protocol", cfg.Kafka.SecurityProtocol)
	}
	if cfg.Kafka.SASLMechanism != "" {
		conf.SetKey("sasl.mechanism", cfg.Kafka.SASLMechanism)
	}
	// Default: no auth needed for local dev
	if cfg.Auth.Username != "" {
		conf.SetKey("sasl.username", cfg.Auth.Username)
		conf.SetKey("sasl.password", cfg.Auth.Password)
	}

	producer, err := kafka.NewProducer(conf)
	if err != nil {
		return nil, fmt.Errorf("mqxbinding: kafka producer: %w", err)
	}

	cache[key] = producer
	return producer, nil
}
