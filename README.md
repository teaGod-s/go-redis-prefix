# go-redis-prefix

`go-redis-prefix` is a Go library for adding prefixes to Redis keys. It leverages Redis hooks to automatically add a specified prefix to keys when executing Redis commands.

## Installation

Install using `go get`:

```sh
go get -u github.com/teaGod-s/go-redis-prefix
```

## Usage

### 1. Initialize Redis Client and Add Prefix Hook

When initializing the Redis client, add the `AppPrefixHook` hook and specify the prefix:

```go
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

    // Example: Key with prefix
    // GET prefix4k:hello
    Cli.Get(ctx, "hello")

    // Example: Key without prefix
    // GET hello
    Cli.Get(prefix.WithSkipPrefix(ctx), "hello")
}
```

### 2. Skip Prefix

If certain keys do not need a prefix, use the `WithSkipPrefix` function:

```go
ctx := prefix.WithSkipPrefix(context.Background())
Cli.Get(ctx, "hello") // This will not add a prefix
```

## Testing

Run tests using `go test`:

```sh
go test
```

Test files are located in `redis_cluster_test.go`.

## Contributing

Contributions are welcome! Please submit a Pull Request or report an Issue.

## License

This project is licensed under the MIT License. See the LICENSE file for details.

## Contact

For any questions, please contact the project maintainers.

---
Thank you for using go-redis-prefix! We hope it helps you manage Redis keys more conveniently.