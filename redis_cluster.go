package prefix

import (
	"context"
	"fmt"
	"net"
	"strings"

	"github.com/redis/go-redis/v9"
	"github.com/samber/lo"
	"github.com/spf13/cast"
)

type contextKey string

// global flag in the hook function
const skipPrefixKey contextKey = "skip"

// WithSkipPrefix define a context helper function, when a key does not need a prefix, use this function, example: Cli.Set(WithSkipPrefix(ctx), "key", "value")
func WithSkipPrefix(ctx context.Context) context.Context {
	return context.WithValue(ctx, skipPrefixKey, true)
}

func shouldSkipPrefix(ctx context.Context) bool {
	value := ctx.Value(skipPrefixKey)
	skip, ok := value.(bool)
	return ok && skip
}

// allow prefix single `key` command
var commandsWithPrefix = []string{
	"GET", "SET", "APPEND", "GETRANGE", "SETRANGE", "STRLEN", "GETSET", "SETNX", "SETEX", "PSETEX", "GETBIT", "SETBIT", "BITCOUNT", "BITPOS", "BITFIELD",
	"RPUSH", "LPOP", "RPOP", "LLEN", "LRANGE", "LPUSH", "LINDEX", "LSET", "LINSERT", "LREM", "LTRIM",
	"SADD", "SREM", "SISMEMBER", "SMEMBERS", "SCARD", "SPOP", "SRANDMEMBER",
	"HSET", "HMSET", "HGET", "HGETALL", "HVALS", "HLEN", "HEXISTS", "HDEL", "HKEYS", "HINCRBY", "HINCRBYFLOAT", "HSCAN", "HSTRLEN",
	"ZADD", "ZRANGE", "ZRANGEBYSCORE", "ZREVRANGEBYSCORE", "ZREM", "ZREVRANGE", "ZCARD", "ZSCORE", "ZRANK", "ZREVRANK", "ZINCRBY", "ZRANGEBYLEX", "ZREVRANGEBYLEX",
	"ZREMRANGEBYRANK", "ZREMRANGEBYSCORE", "ZREMRANGEBYLEX", "ZPOPMIN", "ZPOPMAX",
	"PFADD",
	"GEOADD", "GEOPOS", "GEODIST", "GEOSEARCH",
	"XADD", "XLEN", "XRANGE", "XREVRANGE", "XTRIM", "XDEL",
	"INCR", "INCRBY", "INCRBYFLOAT", "DECR", "DECRBY",
	"WATCH", "MULTI", "EXEC", "EXPIRE", "TTL", "TYPE", "DUMP", "RESTORE",
}

type AppPrefixHook struct {
	Prefix string
}

func (h AppPrefixHook) DialHook(next redis.DialHook) redis.DialHook {
	return func(ctx context.Context, network, addr string) (net.Conn, error) {
		return next(ctx, network, addr)
	}
}

func (h AppPrefixHook) ProcessHook(next redis.ProcessHook) redis.ProcessHook {
	return func(ctx context.Context, cmd redis.Cmder) error {
		if !shouldSkipPrefix(ctx) {
			h.addPrefixToArgs(ctx, cmd)
		}
		return next(ctx, cmd)
	}
}

func (h AppPrefixHook) ProcessPipelineHook(next redis.ProcessPipelineHook) redis.ProcessPipelineHook {
	return func(ctx context.Context, cmds []redis.Cmder) error {
		if !shouldSkipPrefix(ctx) {
			for _, cmd := range cmds {
				h.addPrefixToArgs(ctx, cmd)
			}
		}
		return next(ctx, cmds)
	}
}

// public prefix processing function
func (h AppPrefixHook) addPrefixToArgs(ctx context.Context, cmd redis.Cmder) {
	// directly change the args variable, because the memory address is the same
	args := cmd.Args()
	if len(args) <= 1 {
		return
	}

	name := strings.ToUpper(cmd.Name())
	switch name {
	case "MGET", "DEL", "EXISTS", "TOUCH", "UNLINK", "RENAME", "RENAMENX", "PFMERGE", "SINTERSTORE",
		"SUNIONSTORE", "SDIFFSTORE", "SDIFF", "SINTER", "SUNION", "PFCOUNT":
		// common multi `key` command
		for i := 1; i < len(args); i++ {
			args[i] = h.Prefix + cast.ToString(args[i])
		}
	case "MSET": // MSET key1 value1 key2 value2 ...
		for i := 1; i < len(args); i += 2 {
			args[i] = h.Prefix + cast.ToString(args[i])
		}
	case "BITOP": // BITOP operation destkey key1 key2 ...
		for i := 2; i < len(args); i++ {
			args[i] = h.Prefix + cast.ToString(args[i])
		}
	case "BRPOP", "BLPOP", "BRPOPLPUSH", "BZPOPMIN", "BZPOPMAX": // BRPOP key [key ...] timeout
		for i := 1; i < len(args)-1; i++ {
			args[i] = h.Prefix + cast.ToString(args[i])
		}
	case "XINFO", "XGROUP":
		if len(args) > 2 {
			args[2] = h.Prefix + cast.ToString(args[2])
		}
	case "RPOPLPUSH", "LMOVE", "BLMOVE", "SMOVE", "GEOSEARCHSTORE":
		if len(args) > 2 {
			args[1] = h.Prefix + cast.ToString(args[1])
			args[2] = h.Prefix + cast.ToString(args[2])
		}
	case "SCAN":
		if len(args) > 2 {
			for i := 2; i < len(args); i += 2 {
				if strings.ToUpper(cast.ToString(args[i])) == "MATCH" && i+1 < len(args) {
					args[i+1] = h.Prefix + cast.ToString(args[i+1])
					break
				}
			}
		}
	case "SSCAN", "ZSCAN":
		if len(args) > 3 {
			args[1] = h.Prefix + cast.ToString(args[1])
			for i := 3; i < len(args); i += 2 {
				if strings.ToUpper(cast.ToString(args[i])) == "MATCH" && i+1 < len(args) {
					args[i+1] = h.Prefix + cast.ToString(args[i+1])
					break
				}
			}
		}
	case "SORT":
		// SORT command may have `key` and `BY` clause
		if len(args) > 1 {
			args[1] = h.Prefix + cast.ToString(args[1])
			for i := 2; i < len(args); i++ {
				argsI := strings.ToUpper(cast.ToString(args[i]))
				if argsI == "BY" || argsI == "GET" {
					if i+1 < len(args) {
						args[i+1] = h.Prefix + cast.ToString(args[i+1])
					}
				}
			}
		}
	case "ZDIFF", "ZINTER", "ZUNION":
		// ZUNION `key` parameter starts from the second parameter
		if len(args) > 2 {
			numKeys := cast.ToInt64(args[1])
			if numKeys > 0 {
				for i := 2; i < 2+int(numKeys); i++ {
					args[i] = h.Prefix + cast.ToString(args[i])
				}
			}
		}
	case "ZUNIONSTORE", "ZINTERSTORE":
		if len(args) > 1 {
			args[1] = h.Prefix + cast.ToString(args[1])
		}
		if len(args) > 3 {
			numKeys := cast.ToInt64(args[2])
			if numKeys > 0 {
				for i := 3; i < 3+int(numKeys); i++ {
					args[i] = h.Prefix + cast.ToString(args[i])
				}
			}
		}
	case "EVAL", "EVALSHA":
		// EVAL and EVALSHA `key` parameter starts from the third parameter
		if len(args) > 3 {
			numKeys := cast.ToInt64(args[2])
			if numKeys > 0 {
				for i := 3; i < 3+int(numKeys); i++ {
					args[i] = h.Prefix + cast.ToString(args[i])
				}
			}
		}
	case "MIGRATE":
		if len(args) > 4 {
			if cast.ToString(args[3]) != "" {
				args[3] = h.Prefix + cast.ToString(args[3])
			}
			keysIndex := -1
			for i := 4; i < len(args); i++ {
				if strings.ToUpper(cast.ToString(args[i])) == "KEYS" {
					keysIndex = i
					break
				}
			}
			if keysIndex > 0 {
				for i := keysIndex; i < len(args); i++ {
					args[i] = h.Prefix + cast.ToString(args[i])
				}
			}
		}
	default:
		if lo.IndexOf[string](commandsWithPrefix, name) != -1 {
			args[1] = h.Prefix + cast.ToString(args[1])
		} else {
			fmt.Println("unsupport app prefix command: ", name)
		}
	}
}
