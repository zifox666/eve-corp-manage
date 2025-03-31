package cache

import (
	"time"
)

// Cache 缓存接口-支持Redis和内存使用
type Cache[T any] interface {
	// Set 设置
	Set(k string, v T, d time.Duration)

	// Get 取值
	Get(k string) (T, bool)

	// SetDefault 设置-过期时间采用默认值
	SetDefault(k string, v T)

	// Delete 删除
	Delete(k string)

	// 只有在给定Key项尚未存在，或者现有项已过期时，才能将项添加到缓存中。否则返回错误。
	// Add(k string, v T, d time.Duration)
	// IncrementInt(k string, n int) (num int, err error)

	// SetKeepExpiration 设置值，但不重置过期时间
	SetKeepExpiration(k string, v T)

	// ItemCount 项目总数
	ItemCount() (int64, error)

	// Flush 清空
	Flush()
}
