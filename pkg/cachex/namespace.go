package cachex

import (
	"os"
	"strings"
	"sync"
)

// Namespace isolates keys between different cache use-cases. (e.g. "channel_affinity:v1").
type Namespace string

var (
	globalRedisKeyPrefix     string
	globalRedisKeyPrefixOnce sync.Once
)

// getGlobalRedisKeyPrefix 获取全局 Redis key 前缀（从环境变量）
func getGlobalRedisKeyPrefix() string {
	globalRedisKeyPrefixOnce.Do(func() {
		globalRedisKeyPrefix = os.Getenv("REDIS_KEY_PREFIX")
	})
	return globalRedisKeyPrefix
}

func (n Namespace) prefix() string {
	ns := strings.TrimSpace(string(n))
	ns = strings.TrimRight(ns, ":")
	if ns == "" {
		return ""
	}
	// 添加全局 Redis key 前缀支持
	globalPrefix := getGlobalRedisKeyPrefix()
	if globalPrefix != "" {
		return globalPrefix + ":" + ns + ":"
	}
	return ns + ":"
}

func (n Namespace) FullKey(key string) string {
	key = strings.TrimSpace(key)
	if key == "" {
		return ""
	}
	p := n.prefix()
	if p == "" {
		return strings.TrimLeft(key, ":")
	}
	if strings.HasPrefix(key, p) {
		return key
	}
	return p + strings.TrimLeft(key, ":")
}

func (n Namespace) MatchPattern() string {
	p := n.prefix()
	if p == "" {
		return "*"
	}
	return p + "*"
}
