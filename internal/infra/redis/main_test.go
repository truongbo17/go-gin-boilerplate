package redis

import (
	"testing"

	"github.com/truongbo17/go-gin-boilerplate/config"
)

func TestFailedConnectionDoesNotPublishClient(t *testing.T) {
	previousConfig, previousClient := config.EnvConfig, ClientRedis
	t.Cleanup(func() {
		config.EnvConfig, ClientRedis = previousConfig, previousClient
	})
	ClientRedis = nil
	config.EnvConfig = &config.Config{Cache: config.Cache{CacheStore: config.CacheStoreRedis, RedisHost: "127.0.0.1", RedisPort: "0"}}
	if err := ConnectRedis(); err == nil {
		t.Fatal("expected Redis connection failure")
	}
	if ClientRedis != nil {
		t.Fatal("published a failed Redis client")
	}
}
