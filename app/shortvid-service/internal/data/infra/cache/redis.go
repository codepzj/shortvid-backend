package cache

import (
	"context"
	"log"

	"shortvid-backend/app/shortvid-service/internal/conf"

	"github.com/redis/go-redis/v9"
)

// prefixHook 为所有Redis键添加前缀的Hook
type prefixHook struct {
	prefix string
}

func (h *prefixHook) DialHook(next redis.DialHook) redis.DialHook {
	return next
}

func (h *prefixHook) ProcessHook(next redis.ProcessHook) redis.ProcessHook {
	return func(ctx context.Context, cmd redis.Cmder) error {
		if h.prefix != "" {
			h.addPrefix(cmd)
		}
		return next(ctx, cmd)
	}
}

func (h *prefixHook) ProcessPipelineHook(next redis.ProcessPipelineHook) redis.ProcessPipelineHook {
	return func(ctx context.Context, cmds []redis.Cmder) error {
		if h.prefix != "" {
			for _, cmd := range cmds {
				h.addPrefix(cmd)
			}
		}
		return next(ctx, cmds)
	}
}

// addPrefix 为命令中的key添加前缀
func (h *prefixHook) addPrefix(cmd redis.Cmder) {
	args := cmd.Args()
	if len(args) == 0 {
		return
	}

	cmdName := ""
	if name, ok := args[0].(string); ok {
		cmdName = name
	}

	// 根据不同的命令类型，为相应的key参数添加前缀
	switch cmdName {
	case "get", "set", "del", "exists", "expire", "ttl", "pttl", "type", "dump", "restore",
		"incr", "decr", "incrby", "decrby", "incrbyfloat",
		"append", "getrange", "setrange", "strlen",
		"hget", "hset", "hdel", "hexists", "hgetall", "hkeys", "hvals", "hlen", "hmget", "hmset", "hincrby", "hincrbyfloat", "hsetnx",
		"lpush", "rpush", "lpop", "rpop", "llen", "lrange", "lindex", "lset", "lrem", "ltrim",
		"sadd", "srem", "sismember", "smembers", "scard", "spop", "srandmember",
		"zadd", "zrem", "zscore", "zincrby", "zcard", "zcount", "zrange", "zrevrange", "zrangebyscore", "zrevrangebyscore", "zrank", "zrevrank", "zremrangebyrank", "zremrangebyscore":
		// 这些命令的第一个参数是key
		if len(args) > 1 {
			if key, ok := args[1].(string); ok {
				args[1] = h.prefix + key
			}
		}
	case "mget", "mset", "msetnx":
		// MGET: key [key ...]
		// MSET: key value [key value ...]
		if cmdName == "mget" {
			for i := 1; i < len(args); i++ {
				if key, ok := args[i].(string); ok {
					args[i] = h.prefix + key
				}
			}
		} else {
			// mset, msetnx: key value pairs
			for i := 1; i < len(args); i += 2 {
				if key, ok := args[i].(string); ok {
					args[i] = h.prefix + key
				}
			}
		}
	case "rename", "renamenx":
		// RENAME key newkey
		if len(args) > 2 {
			if key, ok := args[1].(string); ok {
				args[1] = h.prefix + key
			}
			if newKey, ok := args[2].(string); ok {
				args[2] = h.prefix + newKey
			}
		}
	case "keys", "scan":
		// KEYS pattern
		// SCAN cursor [MATCH pattern]
		// 对于KEYS和SCAN命令，需要在pattern前添加前缀
		if cmdName == "keys" && len(args) > 1 {
			if pattern, ok := args[1].(string); ok {
				args[1] = h.prefix + pattern
			}
		} else if cmdName == "scan" {
			// SCAN cursor [MATCH pattern] [COUNT count]
			for i := 1; i < len(args); i++ {
				if s, ok := args[i].(string); ok && s == "MATCH" && i+1 < len(args) {
					if pattern, ok := args[i+1].(string); ok {
						args[i+1] = h.prefix + pattern
					}
					break
				}
			}
		}
	}
}

func NewRedis(c *conf.Data) *redis.Client {
	log.Println("Redis connect start...")
	opts := &redis.Options{
		Addr:         c.Redis.Addr,
		Password:     c.Redis.Password,
		DB:           int(c.Redis.GetDb()),
		ReadTimeout:  c.Redis.GetReadTimeout().AsDuration(),
		WriteTimeout: c.Redis.GetWriteTimeout().AsDuration(),
	}

	client := redis.NewClient(opts)
	// 如果配置了前缀，添加前缀Hook
	if c.Redis.KeyPrefix != "" {
		client.AddHook(&prefixHook{prefix: c.Redis.KeyPrefix})
	}

	if err := client.Ping(context.Background()).Err(); err != nil {
		log.Fatalf("Redis connect failed: %v", err)
	}
	log.Println("Redis connect success...")
	return client
}
