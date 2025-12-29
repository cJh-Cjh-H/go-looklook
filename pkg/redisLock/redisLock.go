package redisLock

import (
	"github.com/zeromicro/go-zero/core/stores/redis"

	"github.com/google/uuid"
)

// Lock 获取锁
func Lock(redis *redis.Redis, key string, v string, ex int) error {
	// 使用 SET NX EX 命令
	err := redis.Setex(key, v, ex)
	if err != nil {
		return err
	}
	return nil
}

// Unlock 释放锁
func Unlock(redis redis.Redis, key string) error {
	// 使用Lua脚本保证原子性：只有锁的持有者才能删除
	script := `
if redis.call("GET", KEYS[1]) == ARGV[1] then
    return redis.call("DEL", KEYS[1])
else
    return 0
end
`
	result, err := l.client.Eval(ctx, script, []string{l.key}, l.value).Result()
	if err != nil {
		return err
	}

	if result.(int64) == 0 {
		return ErrUnlockFailed
	}
	return nil
}

// 生成随机值
func generateRandomValue() string {
	return "lock:" + uuid.New().String()
}
