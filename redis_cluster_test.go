package prefix

import (
	"context"
	"testing"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/spf13/cast"
	"github.com/stretchr/testify/assert"
)

func TestAppPrefixHook(t *testing.T) {
	Cli := redis.NewUniversalClient(&redis.UniversalOptions{
		Addrs:    []string{"127.0.0.1:7001", "127.0.0.1:7002", "127.0.0.1:7003", "127.0.0.1:7004", "127.0.0.1:7005", "127.0.0.1:7006"},
		Password: "",
	})
	prefix := "prefix4key:"
	Cli.AddHook(AppPrefixHook{Prefix: prefix})
	ctx := context.Background()
	tests := []struct {
		name     string
		cmd      redis.Cmder
		expected []interface{}
	}{
		{
			name:     "SET command",
			cmd:      Cli.Set(ctx, "key", "value", time.Minute),
			expected: []interface{}{"set", prefix + "key", "value", "ex", 60},
		},
		{
			name:     "GET command",
			cmd:      Cli.Get(ctx, "key"),
			expected: []interface{}{"get", prefix + "key"},
		},
		{
			name:     "APPEND command",
			cmd:      Cli.Append(ctx, "key", "value"),
			expected: []interface{}{"append", prefix + "key", "value"},
		},
		{
			name:     "GETRANGE command",
			cmd:      Cli.GetRange(ctx, "key", 0, -1),
			expected: []interface{}{"getrange", prefix + "key", "0", "-1"},
		},
		{
			name:     "SETRANGE command",
			cmd:      Cli.SetRange(ctx, "key", 10, "value"),
			expected: []interface{}{"setrange", prefix + "key", "10", "value"},
		},
		{
			name:     "STRLEN command",
			cmd:      Cli.StrLen(ctx, "key"),
			expected: []interface{}{"strlen", prefix + "key"},
		},
		{
			name:     "GETSET command",
			cmd:      Cli.GetSet(ctx, "key", "value"),
			expected: []interface{}{"getset", prefix + "key", "value"},
		},
		{
			name:     "SETNX command",
			cmd:      Cli.SetNX(ctx, "key", "value", time.Minute),
			expected: []interface{}{"set", prefix + "key", "value", "ex", 60, "nx"},
		},
		{
			name:     "SETEX command",
			cmd:      Cli.SetEx(ctx, "key", "value", time.Minute),
			expected: []interface{}{"setex", prefix + "key", 60, "value"},
		},
		//{
		//	name:     "PSETEX command",
		//	cmd: ,
		//	expected: []interface{}{}
		//},
		{
			name:     "GETBIT command",
			cmd:      Cli.GetBit(ctx, "key", 0),
			expected: []interface{}{"getbit", prefix + "key", 0},
		},
		{
			name:     "SETBIT command",
			cmd:      Cli.SetBit(ctx, "key", 0, 1),
			expected: []interface{}{"setbit", prefix + "key", 0, "1"},
		},
		{
			name:     "BITCOUNT command",
			cmd:      Cli.BitCount(ctx, "key", nil),
			expected: []interface{}{"bitcount", prefix + "key"},
		},
		{
			name:     "BITPOS command",
			cmd:      Cli.BitPos(ctx, "key", 0),
			expected: []interface{}{"bitpos", prefix + "key", 0},
		},
		{
			name:     "BITFIELD command",
			cmd:      Cli.BitField(ctx, "key", "value"),
			expected: []interface{}{"bitfield", prefix + "key", "value"},
		},
		{
			name:     "RPUSH command",
			cmd:      Cli.RPush(ctx, "key", "value"),
			expected: []interface{}{"rpush", prefix + "key", "value"},
		},
		{
			name:     "LPUSH command",
			cmd:      Cli.LPush(ctx, "key", "value"),
			expected: []interface{}{"lpush", prefix + "key", "value"},
		},
		{
			name:     "LPOP command",
			cmd:      Cli.LPop(ctx, "key"),
			expected: []interface{}{"lpop", prefix + "key"},
		},
		{
			name:     "RPOP command",
			cmd:      Cli.RPop(ctx, "key"),
			expected: []interface{}{"rpop", prefix + "key"},
		},
		{
			name:     "LLEN command",
			cmd:      Cli.LLen(ctx, "key"),
			expected: []interface{}{"llen", prefix + "key"},
		},
		{
			name:     "LRANGE command",
			cmd:      Cli.LRange(ctx, "key", 0, -1),
			expected: []interface{}{"lrange", prefix + "key", "0", "-1"},
		},
		{
			name:     "LINDEX command",
			cmd:      Cli.LIndex(ctx, "key", 0),
			expected: []interface{}{"lindex", prefix + "key", "0"},
		},
		{
			name:     "LSET command",
			cmd:      Cli.LSet(ctx, "key", 0, "value"),
			expected: []interface{}{"lset", prefix + "key", "0", "value"},
		},
		{
			name:     "LINSERT command",
			cmd:      Cli.LInsert(ctx, "key", "BEFORE", "作为参考点的现有元素", "value"),
			expected: []interface{}{"linsert", prefix + "key", "BEFORE", "作为参考点的现有元素", "value"},
		},
		{
			name:     "LREM command",
			cmd:      Cli.LRem(ctx, "key", 0, "value"),
			expected: []interface{}{"lrem", prefix + "key", "0", "value"},
		},
		{
			name:     "LTRIM command",
			cmd:      Cli.LTrim(ctx, "key", 0, -1),
			expected: []interface{}{"ltrim", prefix + "key", "0", "-1"},
		},
		{
			name:     "RPOPLPUSH command",
			cmd:      Cli.RPopLPush(ctx, "key1", "key2"),
			expected: []interface{}{"rpoplpush", prefix + "key1", prefix + "key2"},
		},
		{
			name:     "SADD command",
			cmd:      Cli.SAdd(ctx, "key", "value"),
			expected: []interface{}{"sadd", prefix + "key", "value"},
		},
		{
			name:     "SREM command",
			cmd:      Cli.SRem(ctx, "key", "value"),
			expected: []interface{}{"srem", prefix + "key", "value"},
		},
		{
			name:     "SISMEMBER command",
			cmd:      Cli.SIsMember(ctx, "key", "value"),
			expected: []interface{}{"sismember", prefix + "key", "value"},
		},
		{
			name:     "SMEMBERS command",
			cmd:      Cli.SMembers(ctx, "key"),
			expected: []interface{}{"smembers", prefix + "key"},
		},
		{
			name:     "SCARD command",
			cmd:      Cli.SCard(ctx, "key"),
			expected: []interface{}{"scard", prefix + "key"},
		},
		{
			name:     "SPOP command",
			cmd:      Cli.SPop(ctx, "key"),
			expected: []interface{}{"spop", prefix + "key"},
		},
		{
			name:     "SRANDMEMBER command",
			cmd:      Cli.SRandMember(ctx, "key"),
			expected: []interface{}{"srandmember", prefix + "key"},
		},
		{
			name:     "HSET command",
			cmd:      Cli.HSet(ctx, "key", "value", "value"),
			expected: []interface{}{"hset", prefix + "key", "value", "value"},
		},
		{
			name:     "HMSET command",
			cmd:      Cli.HMSet(ctx, "key", "value", "value"),
			expected: []interface{}{"hmset", prefix + "key", "value", "value"},
		},
		{
			name:     "HGET command",
			cmd:      Cli.HGet(ctx, "key", "value"),
			expected: []interface{}{"hget", prefix + "key", "value"},
		},
		{
			name:     "HGETALL command",
			cmd:      Cli.HGetAll(ctx, "key"),
			expected: []interface{}{"hgetall", prefix + "key"},
		},
		{
			name:     "HVALS command",
			cmd:      Cli.HVals(ctx, "key"),
			expected: []interface{}{"hvals", prefix + "key"},
		},
		{
			name:     "HLEN command",
			cmd:      Cli.HLen(ctx, "key"),
			expected: []interface{}{"hlen", prefix + "key"},
		},
		{
			name:     "HEXISTS command",
			cmd:      Cli.HExists(ctx, "key", "value"),
			expected: []interface{}{"hexists", prefix + "key", "value"},
		},
		{
			name:     "HDEL command",
			cmd:      Cli.HDel(ctx, "key", "value"),
			expected: []interface{}{"hdel", prefix + "key", "value"},
		},
		{
			name:     "HKEYS command",
			cmd:      Cli.HKeys(ctx, "key"),
			expected: []interface{}{"hkeys", prefix + "key"},
		},
		{
			name:     "HINCRBY command",
			cmd:      Cli.HIncrBy(ctx, "key", "value", 1),
			expected: []interface{}{"hincrby", prefix + "key", "value", "1"},
		},
		{
			name:     "HINCRBYFLOAT command",
			cmd:      Cli.HIncrByFloat(ctx, "key", "value", 1),
			expected: []interface{}{"hincrbyfloat", prefix + "key", "value", "1"},
		},
		{
			name:     "HSCAN command",
			cmd:      Cli.HScan(ctx, "key", 0, "field前缀", 10),
			expected: []interface{}{"hscan", prefix + "key", 0, "match", "field前缀", "count", 10},
		},
		//{
		//	name:     "HSTRLEN command",
		//	cmd:      ,
		//	expected: []interface{}{"hstrlen", prefix+"key"},
		//},
		{
			name: "ZADD command",
			cmd: Cli.ZAdd(ctx, "key", redis.Z{
				Score:  0,
				Member: "",
			}),
			expected: []interface{}{"zadd", prefix + "key", "0", redis.Z{
				Score:  0,
				Member: "",
			}},
		},
		{
			name:     "ZRANGE command",
			cmd:      Cli.ZRange(ctx, "key", 0, -1),
			expected: []interface{}{"zrange", prefix + "key", 0, -1},
		},
		{
			name:     "ZRANGEBYSCORE command",
			cmd:      Cli.ZRangeByScore(ctx, "key", &redis.ZRangeBy{}),
			expected: []interface{}{"zrangebyscore", prefix + "key", "", ""},
		},
		{
			name:     "ZREVRANGEBYSCORE command",
			cmd:      Cli.ZRevRangeByScore(ctx, "key", &redis.ZRangeBy{}),
			expected: []interface{}{"zrevrangebyscore", prefix + "key", "", ""},
		},
		{
			name:     "ZREM command",
			cmd:      Cli.ZRem(ctx, "key", "value"),
			expected: []interface{}{"zrem", prefix + "key", "value"},
		},
		{
			name:     "ZREVRANGE command",
			cmd:      Cli.ZRevRange(ctx, "key", 0, -1),
			expected: []interface{}{"zrevrange", prefix + "key", 0, -1},
		},
		{
			name:     "ZCARD command",
			cmd:      Cli.ZCard(ctx, "key"),
			expected: []interface{}{"zcard", prefix + "key"},
		},
		{
			name:     "ZSCORE command",
			cmd:      Cli.ZScore(ctx, "key", "value"),
			expected: []interface{}{"zscore", prefix + "key", "value"},
		},
		{
			name:     "ZRANK command",
			cmd:      Cli.ZRank(ctx, "key", "value"),
			expected: []interface{}{"zrank", prefix + "key", "value"},
		},
		{
			name:     "ZREVRANK command",
			cmd:      Cli.ZRevRank(ctx, "key", "value"),
			expected: []interface{}{"zrevrank", prefix + "key", "value"},
		},
		{
			name:     "ZINCRBY command",
			cmd:      Cli.ZIncrBy(ctx, "key", 1.0, "value"),
			expected: []interface{}{"zincrby", prefix + "key", 1.0, "value"},
		},
		{
			name:     "ZRANGEBYLEX command",
			cmd:      Cli.ZRangeByLex(ctx, "key", &redis.ZRangeBy{}),
			expected: []interface{}{"zrangebylex", prefix + "key", "", ""},
		},
		{
			name:     "ZREVRANGEBYLEX command",
			cmd:      Cli.ZRevRangeByLex(ctx, "key", &redis.ZRangeBy{}),
			expected: []interface{}{"zrevrangebylex", prefix + "key", "", ""},
		},
		{
			name:     "ZREMRANGEBYRANK command",
			cmd:      Cli.ZRemRangeByRank(ctx, "key", 0, -1),
			expected: []interface{}{"zremrangebyrank", prefix + "key", 0, -1},
		},
		{
			name:     "ZREMRANGEBYSCORE command",
			cmd:      Cli.ZRemRangeByScore(ctx, "key", "", ""),
			expected: []interface{}{"zremrangebyscore", prefix + "key", "", ""},
		},
		{
			name:     "ZREMRANGEBYLEX command",
			cmd:      Cli.ZRemRangeByLex(ctx, "key", "", ""),
			expected: []interface{}{"zremrangebylex", prefix + "key", "", ""},
		},
		{
			name:     "ZPOPMIN command",
			cmd:      Cli.ZPopMin(ctx, "key"),
			expected: []interface{}{"zpopmin", prefix + "key"},
		},
		{
			name:     "ZPOPMAX command",
			cmd:      Cli.ZPopMax(ctx, "key"),
			expected: []interface{}{"zpopmax", prefix + "key"},
		},
		{
			name:     "PFADD command",
			cmd:      Cli.PFAdd(ctx, "key", "value"),
			expected: []interface{}{"pfadd", prefix + "key", "value"},
		},
		{
			name: "GEOADD command",
			cmd: Cli.GeoAdd(ctx, "key", &redis.GeoLocation{
				Name:      "",
				Longitude: 0,
				Latitude:  0,
				Dist:      0,
				GeoHash:   0,
			}),
			expected: []interface{}{"geoadd", prefix + "key", "0", "0", ""},
		},
		{
			name:     "GEOPOS command",
			cmd:      Cli.GeoPos(ctx, "key", "value"),
			expected: []interface{}{"geopos", prefix + "key", "value"},
		},
		{
			name:     "GEODIST command",
			cmd:      Cli.GeoDist(ctx, "key", "member1", "member2", "km"),
			expected: []interface{}{"geodist", prefix + "key", "member1", "member2", "km"},
		},
		{
			name:     "GEOSEARCH",
			cmd:      Cli.GeoSearch(ctx, "key", &redis.GeoSearchQuery{}),
			expected: []interface{}{"geosearch", prefix + "key", "fromlonlat", "0", "0", "bybox", "0", "0", "km"},
		},
		{
			name: "XADD command",
			cmd: Cli.XAdd(ctx, &redis.XAddArgs{
				Stream:     "key",
				NoMkStream: false,
				MaxLen:     0,
				MinID:      "111",
				Approx:     false,
				Limit:      0,
				ID:         "我的ID",
				Values:     []interface{}{"field1", "value1", "field2", "value2"},
			}),
			expected: []interface{}{"xadd", prefix + "key", "minid", "111", "我的ID", "field1", "value1", "field2", "value2"},
		},
		{
			name:     "XLEN command",
			cmd:      Cli.XLen(ctx, "key"),
			expected: []interface{}{"xlen", prefix + "key"},
		},
		{
			name:     "XRANGE command",
			cmd:      Cli.XRange(ctx, "key", "start", "stop"),
			expected: []interface{}{"xrange", prefix + "key", "start", "stop"},
		},
		{
			name:     "XREVRANGE command",
			cmd:      Cli.XRevRange(ctx, "key", "start", "stop"),
			expected: []interface{}{"xrevrange", prefix + "key", "start", "stop"},
		},
		{
			name:     "XTRIM command",
			cmd:      Cli.XTrimMinID(ctx, "key", "1"),
			expected: []interface{}{"xtrim", prefix + "key", "minid", "1"},
		},
		{
			name:     "XDEL command",
			cmd:      Cli.XDel(ctx, "key", "1"),
			expected: []interface{}{"xdel", prefix + "key", "1"},
		},
		{
			name:     "INCR command",
			cmd:      Cli.Incr(ctx, "key"),
			expected: []interface{}{"incr", prefix + "key"},
		},
		{
			name:     "INCRBY command",
			cmd:      Cli.IncrBy(ctx, "key", 1.0),
			expected: []interface{}{"incrby", prefix + "key", 1.0},
		},
		{
			name:     "INCRBYFLOAT command",
			cmd:      Cli.IncrByFloat(ctx, "key", 1.0),
			expected: []interface{}{"incrbyfloat", prefix + "key", 1.0},
		},
		{
			name:     "DECR command",
			cmd:      Cli.Decr(ctx, "key"),
			expected: []interface{}{"decr", prefix + "key"},
		},
		{
			name:     "DECRBY command",
			cmd:      Cli.DecrBy(ctx, "key", 1.0),
			expected: []interface{}{"decrby", prefix + "key", 1.0},
		},
		//{
		//	name: "WATCH command",
		//	cmd:  Cli.Watch(ctx, nil, "key1", "key2"),
		//},
		//{
		//	name: "MULTI command",
		//	cmd: Cli.,
		//},
		//{
		//	name:     "EXEC command",
		//	cmd: Cli.exe,
		//},
		{
			name:     "EXPIRE command",
			cmd:      Cli.Expire(ctx, "key", 0),
			expected: []interface{}{"expire", prefix + "key", 0},
		},
		{
			name:     "TTL command",
			cmd:      Cli.TTL(ctx, "key"),
			expected: []interface{}{"ttl", prefix + "key"},
		},
		{
			name:     "TYPE command",
			cmd:      Cli.Type(ctx, "key"),
			expected: []interface{}{"type", prefix + "key"},
		},
		{
			name:     "DUMP command",
			cmd:      Cli.Dump(ctx, "key"),
			expected: []interface{}{"dump", prefix + "key"},
		},
		{
			name:     "RESTORE command",
			cmd:      Cli.Restore(ctx, "key", time.Minute, "value"),
			expected: []interface{}{"restore", prefix + "key", "60000", "value"},
		},
		{
			name:     "MGET command",
			cmd:      Cli.MGet(ctx, "key1", "key2", "key3"),
			expected: []interface{}{"mget", prefix + "key1", prefix + "key2", prefix + "key3"},
		},
		{
			name:     "DEL command",
			cmd:      Cli.Del(ctx, "key1", "key2", "key3"),
			expected: []interface{}{"del", prefix + "key1", prefix + "key2", prefix + "key3"},
		},
		{
			name:     "EXISTS command",
			cmd:      Cli.Exists(ctx, "key1", "key2", "key3"),
			expected: []interface{}{"exists", prefix + "key1", prefix + "key2", prefix + "key3"},
		},
		{
			name:     "TOUCH command",
			cmd:      Cli.Touch(ctx, "key1", "key2", "key3"),
			expected: []interface{}{"touch", prefix + "key1", prefix + "key2", prefix + "key3"},
		},
		{
			name:     "UNLINK command",
			cmd:      Cli.Unlink(ctx, "key1", "key2", "key3"),
			expected: []interface{}{"unlink", prefix + "key1", prefix + "key2", prefix + "key3"},
		},
		{
			name:     "RENAME command",
			cmd:      Cli.Rename(ctx, "key1", "key2"),
			expected: []interface{}{"rename", prefix + "key1", prefix + "key2"},
		},
		{
			name:     "RENAMENX command",
			cmd:      Cli.RenameNX(ctx, "key1", "key2"),
			expected: []interface{}{"renamenx", prefix + "key1", prefix + "key2"},
		},
		{
			name:     "PFMERGE command",
			cmd:      Cli.PFMerge(ctx, "key1", "key2", "key3"),
			expected: []interface{}{"pfmerge", prefix + "key1", prefix + "key2", prefix + "key3"},
		},
		{
			name:     "SINTERSTORE command",
			cmd:      Cli.SInterStore(ctx, "key1", "key2", "key3"),
			expected: []interface{}{"sinterstore", prefix + "key1", prefix + "key2", prefix + "key3"},
		},
		{
			name:     "SUNIONSTORE command",
			cmd:      Cli.SUnionStore(ctx, "key1", "key2", "key3"),
			expected: []interface{}{"sunionstore", prefix + "key1", prefix + "key2", prefix + "key3"},
		},
		{
			name:     "SDIFFSTORE command",
			cmd:      Cli.SDiffStore(ctx, "key1", "key2", "key3"),
			expected: []interface{}{"sdiffstore", prefix + "key1", prefix + "key2", prefix + "key3"},
		},
		{
			name:     "SDIFF command",
			cmd:      Cli.SDiff(ctx, "key1", "key2", "key3"),
			expected: []interface{}{"sdiff", prefix + "key1", prefix + "key2", prefix + "key3"},
		},
		{
			name:     "SINTER command",
			cmd:      Cli.SInter(ctx, "key1", "key2", "key3"),
			expected: []interface{}{"sinter", prefix + "key1", prefix + "key2", prefix + "key3"},
		},
		{
			name:     "SUNION command",
			cmd:      Cli.SUnion(ctx, "key1", "key2", "key3"),
			expected: []interface{}{"sunion", prefix + "key1", prefix + "key2", prefix + "key3"},
		},
		{
			name:     "PFCOUNT command",
			cmd:      Cli.PFCount(ctx, "key1", "key2", "key3"),
			expected: []interface{}{"pfcount", prefix + "key1", prefix + "key2", prefix + "key3"},
		},
		{
			name:     "MSET command",
			cmd:      Cli.MSet(ctx, "key1", "value1", "key2", "value2"),
			expected: []interface{}{"mset", prefix + "key1", "value1", prefix + "key2", "value2"},
		},
		{
			name:     "BITOP command",
			cmd:      Cli.BitOpOr(ctx, "key1", "key2", "key3"),
			expected: []interface{}{"bitop", "or", prefix + "key1", prefix + "key2", prefix + "key3"},
		},
		{
			name:     "BRPOP command",
			cmd:      Cli.BRPop(ctx, time.Minute, "key1", "key2", "key3"),
			expected: []interface{}{"brpop", prefix + "key1", prefix + "key2", prefix + "key3", 60},
		},
		{
			name:     "BLPOP command",
			cmd:      Cli.BLPop(ctx, time.Minute, "key1", "key2", "key3"),
			expected: []interface{}{"blpop", prefix + "key1", prefix + "key2", prefix + "key3", 60},
		},
		{
			name:     "BRPOPLPUSH command",
			cmd:      Cli.BRPopLPush(ctx, "key1", "key2", time.Minute),
			expected: []interface{}{"brpoplpush", prefix + "key1", prefix + "key2", 60},
		},
		{
			name:     "BZPOPMIN command",
			cmd:      Cli.BZPopMin(ctx, time.Minute, "key1", "key2", "key3"),
			expected: []interface{}{"bzpopmin", prefix + "key1", prefix + "key2", prefix + "key3", 60},
		},
		{
			name:     "BZPOPMAX command",
			cmd:      Cli.BZPopMax(ctx, time.Minute, "key1", "key2", "key3"),
			expected: []interface{}{"bzpopmax", prefix + "key1", prefix + "key2", prefix + "key3", 60},
		},
		{
			name:     "XINFO command",
			cmd:      Cli.XInfoStream(ctx, "key"),
			expected: []interface{}{"xinfo", "stream", prefix + "key"},
		},
		{
			name:     "XGROUP command",
			cmd:      Cli.XGroupCreateMkStream(ctx, "key", "", ""),
			expected: []interface{}{"xgroup", "create", prefix + "key", "", "", "mkstream"},
		},
		{
			name:     "XGROUP command",
			cmd:      Cli.XGroupCreateConsumer(ctx, "key", "", ""),
			expected: []interface{}{"xgroup", "createconsumer", prefix + "key", "", ""},
		},
		{
			name:     "RPOPLPUSH command",
			cmd:      Cli.RPopLPush(ctx, "key1", "key2"),
			expected: []interface{}{"rpoplpush", prefix + "key1", prefix + "key2"},
		},
		{
			name:     "LMOVE command",
			cmd:      Cli.LMove(ctx, "key1", "key2", "RIGHT", "LEFT"),
			expected: []interface{}{"lmove", prefix + "key1", prefix + "key2", "RIGHT", "LEFT"},
		},
		{
			name:     "BLMOVE command",
			cmd:      Cli.BLMove(ctx, "key1", "key2", "RIGHT", "LEFT", time.Minute),
			expected: []interface{}{"blmove", prefix + "key1", prefix + "key2", "RIGHT", "LEFT", "60"},
		},
		{
			name:     "SMOVE command",
			cmd:      Cli.SMove(ctx, "key1", "key2", "value"),
			expected: []interface{}{"smove", prefix + "key1", prefix + "key2", "value"},
		},
		{
			name: "GEOSEARCHSTORE command",
			cmd: Cli.GeoSearchStore(ctx, "key1", "key2", &redis.GeoSearchStoreQuery{
				GeoSearchQuery: redis.GeoSearchQuery{},
				StoreDist:      false,
			}),
			expected: []interface{}{"geosearchstore", prefix + "key2", prefix + "key1", "fromlonlat", "0", "0", "bybox", "0", "0", "km"},
		},
		{
			name:     "SCAN command",
			cmd:      Cli.Scan(ctx, 0, "no:prefix:key", 100),
			expected: []interface{}{"scan", "0", "match", prefix + "no:prefix:key", "count", "100"},
		},
		{
			name:     "SSCAN command",
			cmd:      Cli.SScan(ctx, "key1", 0, "no:prefix:key", 100),
			expected: []interface{}{"sscan", prefix + "key1", 0, "match", prefix + "no:prefix:key", "count", "100"},
		},
		{
			name:     "ZSCAN command",
			cmd:      Cli.ZScan(ctx, "key1", 0, "no:prefix:key", 100),
			expected: []interface{}{"zscan", prefix + "key1", "0", "match", prefix + "no:prefix:key", "count", "100"},
		},
		{
			name: "SORT command",
			cmd: Cli.Sort(ctx, "key", &redis.Sort{
				By:     "id",
				Offset: 0,
				Count:  0,
				Get:    []string{"name", "age"},
				Order:  "desc",
				Alpha:  false,
			}),
			expected: []interface{}{"sort", prefix + "key", "by", prefix + "id", "get", prefix + "name", "get", prefix + "age", "desc"},
		},
		{
			name:     "ZDIFF command",
			cmd:      Cli.ZDiff(ctx, "key1", "key2", "key3"),
			expected: []interface{}{"zdiff", 3, prefix + "key1", prefix + "key2", prefix + "key3"},
		},
		{
			name: "ZINTER command",
			cmd: Cli.ZInter(ctx, &redis.ZStore{
				Keys:      []string{"key1", "key2", "key3"},
				Weights:   []float64{1, 2},
				Aggregate: "SUM",
			}),
			expected: []interface{}{"zinter", "3", prefix + "key1", prefix + "key2", prefix + "key3", "weights", "1", "2", "aggregate", "SUM"},
		},
		{
			name: "ZUNION command",
			cmd: Cli.ZUnion(ctx, redis.ZStore{
				Keys:      []string{"key1", "key2", "key3"},
				Weights:   []float64{1, 2},
				Aggregate: "SUM",
			}),
			expected: []interface{}{"zunion", "3", prefix + "key1", prefix + "key2", prefix + "key3", "weights", "1", "2", "aggregate", "SUM"},
		},
		{
			name: "ZUNIONSTORE command",
			cmd: Cli.ZUnionStore(ctx, "keyall", &redis.ZStore{
				Keys:      []string{"key1", "key2", "key3"},
				Weights:   []float64{1, 2},
				Aggregate: "SUM",
			}),
			expected: []interface{}{"zunionstore", prefix + "keyall", "3", prefix + "key1", prefix + "key2", prefix + "key3", "weights", "1", "2", "aggregate", "SUM"},
		},
		{
			name: "ZINTERSTORE command",
			cmd: Cli.ZInterStore(ctx, "keyall", &redis.ZStore{
				Keys:      []string{"key1", "key2", "key3"},
				Weights:   []float64{1, 2},
				Aggregate: "SUM",
			}),
			expected: []interface{}{"zinterstore", prefix + "keyall", "3", prefix + "key1", prefix + "key2", prefix + "key3", "weights", "1", "2", "aggregate", "SUM"},
		},
		{
			name:     "EVAL command",
			cmd:      Cli.Eval(ctx, "", []string{"key1", "key2", "key3"}, 1, 2),
			expected: []interface{}{"eval", "", 3, prefix + "key1", prefix + "key2", prefix + "key3", 1, 2},
		},
		{
			name:     "EVALSHA command",
			cmd:      Cli.EvalSha(ctx, "hash", []string{"key1", "key2", "key3"}, 1, 2),
			expected: []interface{}{"evalsha", "hash", 3, prefix + "key1", prefix + "key2", prefix + "key3", 1, 2},
		},
		//{
		//	name:     "GEORADIUS command",
		//	cmd:      Cli.GeoRadius(ctx, "key", 100, 100, &redis.GeoRadiusQuery{}),
		//	expected: []interface{}{"georadius", prefix+"key"},
		//},
		//{
		//	name:     "GEORADIUSBYMEMBER command",
		//	cmd:      Cli.GeoRadiusByMember(ctx, "key", "hehe", &redis.GeoRadiusQuery{}),
		//	expected: []interface{}{"georadiusbymember", prefix+"key", "hehe", 1, 2},
		//},
		{
			name:     "MIGRATE command",
			cmd:      Cli.Migrate(ctx, "127.0.0.1", "6379", "key", 0, time.Minute),
			expected: []interface{}{"migrate", "127.0.0.1", "6379", prefix + "key", "0", "60000"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Log(cast.ToStringSlice(tt.expected))
			assert.Equal(t, cast.ToStringSlice(tt.expected), cast.ToStringSlice(tt.cmd.Args()))
		})
	}
}
