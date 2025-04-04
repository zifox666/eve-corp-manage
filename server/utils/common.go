package utils

import (
	"strconv"
	"strings"
	"time"
)

// StringToIntList 将逗号分隔的字符串转换为无符号整数切片
func StringToIntList(s string) ([]uint, error) {
	if s == "" {
		return []uint{}, nil
	}

	parts := strings.Split(s, ",")
	result := make([]uint, 0, len(parts))

	for _, part := range parts {
		// 去除空白
		trimmed := strings.TrimSpace(part)
		if trimmed == "" {
			continue
		}

		// 转换为无符号整数
		num, err := strconv.ParseUint(trimmed, 10, 64)
		if err != nil {
			return nil, err
		}
		result = append(result, uint(num))
	}

	return result, nil
}

// StringToStringList 将逗号分隔的字符串转换为字符串切片
func StringToStringList(s string) []string {
	if s == "" {
		return []string{}
	}

	parts := strings.Split(s, ",")
	result := make([]string, 0, len(parts))

	for _, part := range parts {
		trimmed := strings.TrimSpace(part)
		if trimmed != "" {
			result = append(result, trimmed)
		}
	}

	return result
}

// StringToDate 将字符串转换为时间
func StringToDate(dateStr string) time.Time {
	t, err := time.Parse(time.RFC3339, dateStr)
	if err != nil {
		return time.Now()
	}
	return t
}
