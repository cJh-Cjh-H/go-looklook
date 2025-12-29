package redisLock

import (
	"errors"
	"fmt"
	"github.com/zeromicro/go-zero/core/stores/redis"
)

// Lock 获取锁
func Lock(redis *redis.Redis, key string, v string, ex int) error {
	// go-zero使用 SetnxEx 方法
	// SetnxEx(key, value string, seconds int) 返回 bool, error
	success, err := redis.SetnxEx(key, v, ex)
	if err != nil {
		return fmt.Errorf("redis setnxex error: %w", err)
	}

	if !success {
		return errors.New("lock already exists")
	}

	return nil
}

// Unlock 释放锁
func Unlock(redisClient *redis.Redis, key string, value string) {
	// 使用Eval执行Lua脚本
	script := `
if redis.call("GET", KEYS[1]) == ARGV[1] then
    return redis.call("DEL", KEYS[1])
else
    return 0
end
`

	// go-zero的Eval方法参数格式不同
	_, err := redisClient.Eval(script, []string{key}, value)
	if err != nil {
		return
	}
	return
}
