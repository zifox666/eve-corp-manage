package cache

import (
	"context"
	"encoding/json"
	"time"

	"github.com/redis/go-redis/v9"
)

type RedisCacheStruct[T any] struct {
	Redis             *redis.Client
	Ctx               context.Context
	HashKey           string
	Result            T
	DefaultExpiration time.Duration
	CleanupInterval   time.Duration
}

type RedisValue[T any] struct {
	ExpirationTimeStamp int64
	IsExpiration        bool // 是否有过期时间 false // 不过期
	Value               T
}

// NewRedisCache 创建一个新的Redis缓存实例
func NewRedisCache[T any](redisClient *redis.Client, hashKey string, defaultExpiration, cleanupInterval time.Duration) *RedisCacheStruct[T] {
	return &RedisCacheStruct[T]{
		Redis:             redisClient,
		Ctx:               context.Background(),
		HashKey:           hashKey,
		DefaultExpiration: defaultExpiration,
		CleanupInterval:   cleanupInterval,
	}
}

// Set 设置缓存，带过期时间
func (r *RedisCacheStruct[T]) Set(k string, v T, d time.Duration) {
	value := RedisValue[T]{
		Value: v,
	}

	if d > 0 {
		value.IsExpiration = true
		value.ExpirationTimeStamp = time.Now().Add(d).UnixNano()
	}

	jsonData, _ := json.Marshal(value)
	r.Redis.HSet(r.Ctx, r.HashKey, k, jsonData)
}

// Get 获取缓存
func (r *RedisCacheStruct[T]) Get(k string) (T, bool) {
	var value RedisValue[T]
	var empty T

	jsonData, err := r.Redis.HGet(r.Ctx, r.HashKey, k).Bytes()
	if err != nil {
		return empty, false
	}

	if err := json.Unmarshal(jsonData, &value); err != nil {
		return empty, false
	}

	// 检查是否过期
	if value.IsExpiration && time.Now().UnixNano() > value.ExpirationTimeStamp {
		r.Delete(k)
		return empty, false
	}

	return value.Value, true
}

// SetDefault 使用默认过期时间设置缓存
func (r *RedisCacheStruct[T]) SetDefault(k string, v T) {
	r.Set(k, v, r.DefaultExpiration)
}

// Delete 删除缓存项
func (r *RedisCacheStruct[T]) Delete(k string) {
	r.Redis.HDel(r.Ctx, r.HashKey, k)
}

// SetKeepExpiration 设置值但不重置过期时间
func (r *RedisCacheStruct[T]) SetKeepExpiration(k string, v T) {
	var value RedisValue[T]

	jsonData, err := r.Redis.HGet(r.Ctx, r.HashKey, k).Bytes()
	if err != nil {
		// 键不存在，使用默认值设置
		r.SetDefault(k, v)
		return
	}

	if err := json.Unmarshal(jsonData, &value); err != nil {
		// 反序列化失败，使用默认值设置
		r.SetDefault(k, v)
		return
	}

	// 保留过期时间，仅更新值
	value.Value = v
	updatedData, _ := json.Marshal(value)
	r.Redis.HSet(r.Ctx, r.HashKey, k, updatedData)
}

// ItemCount 获取缓存项目数量
func (r *RedisCacheStruct[T]) ItemCount() (int64, error) {
	return r.Redis.HLen(r.Ctx, r.HashKey).Result()
}

// Flush 清空缓存
func (r *RedisCacheStruct[T]) Flush() {
	r.Redis.Del(r.Ctx, r.HashKey)
}
