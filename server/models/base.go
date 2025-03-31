package models

import (
	_ "gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// int类型代表是否的常量
const (
	INT_FALSE = iota
	INT_TURE
)

// 分页的结构体
type PageLimitStruct struct {
	PageSize  int `gorm:"-"` //
	LimitSize int `gorm:"-"` //
}

// 计算分页
func calcPage(pageSize, limitSize int) (offset, limit int) {
	offset = limitSize * (pageSize - 1)
	limit = limitSize
	return
}

var Db *gorm.DB
