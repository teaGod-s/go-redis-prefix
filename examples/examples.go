package main

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/teaGod-s/go-redis-prefix"
)

func main() {
	Cli := redis.NewUniversalClient(&redis.UniversalOptions{
		Addrs:    []string{"localhost:6379"},
		Password: "",
	})

	Cli.AddHook(prefix.AppPrefixHook{Prefix: "prefix4k:"})

	// Add context with timeout for ping
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := Cli.Ping(ctx).Err(); err != nil {
		panic(fmt.Sprintf("failed to connect to redis: %v", err))
	}

	// a example for the key with prefix
	// GET prefix4k:hello
	Cli.Get(ctx, "hello")

	// a example for the key with no prefix
	// GET hello
	Cli.Get(prefix.WithSkipPrefix(ctx), "hello")
}
