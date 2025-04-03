package service

import (
	"eve-corp-manager/global"
	"eve-corp-manager/models/service"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// GetUserPapList 获取用户PAP记录列表
func GetUserPapList(c *gin.Context) {
	var req struct {
		UserID uint `json:"userId" form:"userId"`
		Page   int  `json:"page" form:"page"`
		Limit  int  `json:"limit" form:"limit"`
	}

	if err := c.ShouldBind(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "参数错误"})
		return
	}

	if req.Page <= 0 {
		req.Page = 1
	}
	if req.Limit <= 0 {
		req.Limit = 10
	}

	var papRecords []service.CorpPap
	var total int64

	db := global.Db.Model(&service.CorpPap{})
	if req.UserID > 0 {
		db = db.Where("user_id = ?", req.UserID)
	}

	db.Count(&total)
	offset := (req.Page - 1) * req.Limit
	result := db.Order("created_at DESC").Offset(offset).Limit(req.Limit).Find(&papRecords)
	if result.Error != nil {
		global.Logger.Error("获取PAP记录失败:", result.Error)
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "获取PAP记录失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "获取PAP记录成功",
		"data": gin.H{
			"total": total,
			"items": papRecords,
		},
	})
}

// GetUserPapBalance 获取用户PAP余额
func GetUserPapBalance(c *gin.Context) {
	var req struct {
		UserID uint `json:"userId" form:"userId"`
	}

	if err := c.ShouldBind(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "参数错误"})
		return
	}

	if req.UserID <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "用户ID不能为空"})
		return
	}

	// 获取最新的PAP记录以获取余额
	var latestPap service.CorpPap
	result := global.Db.Where("user_id = ?", req.UserID).Order("created_at DESC").First(&latestPap)

	balance := 0
	if result.Error == nil {
		balance = latestPap.Balance
	} else if result.Error != gorm.ErrRecordNotFound {
		global.Logger.Error("获取PAP余额失败:", result.Error)
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "获取PAP余额失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "获取PAP余额成功",
		"data":    balance,
	})
}

// AddUserPap 增加用户PAP
func AddUserPap(c *gin.Context) {
	var req struct {
		UserID   uint   `json:"userId" binding:"required"`
		Amount   int    `json:"amount" binding:"required"`
		Source   string `json:"source"`
		SourceID uint   `json:"sourceId"`
		Remark   string `json:"remark"`
		Operator uint   `json:"operator"` // 操作人ID
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "参数错误"})
		return
	}

	if req.Amount <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "PAP数量必须大于0"})
		return
	}

	// 获取当前余额
	var currentBalance int
	var latestPap service.CorpPap
	result := global.Db.Where("user_id = ?", req.UserID).Order("created_at DESC").First(&latestPap)
	if result.Error == nil {
		currentBalance = latestPap.Balance
	} else if result.Error != gorm.ErrRecordNotFound {
		global.Logger.Error("获取PAP余额失败:", result.Error)
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "获取PAP余额失败"})
		return
	}

	// 计算新余额
	newBalance := currentBalance + req.Amount

	// 开始事务
	tx := global.Db.Begin()

	// 创建PAP记录
	papRecord := service.CorpPap{
		UserID:     req.UserID,
		Amount:     req.Amount,
		Balance:    newBalance,
		Source:     req.Source,
		SourceID:   req.SourceID,
		Type:       1, // 1-获取
		CreateTime: time.Now(),
		Remark:     req.Remark,
	}

	if err := tx.Create(&papRecord).Error; err != nil {
		tx.Rollback()
		global.Logger.Error("创建PAP记录失败:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "创建PAP记录失败"})
		return
	}

	// 创建操作日志
	papLog := service.CorpPapLog{
		UserID:     req.UserID,
		Operation:  "增加PAP",
		Amount:     req.Amount,
		BeforeVal:  currentBalance,
		AfterVal:   newBalance,
		Operator:   req.Operator,
		CreateTime: time.Now(),
		Remark:     req.Remark,
	}

	if err := tx.Create(&papLog).Error; err != nil {
		tx.Rollback()
		global.Logger.Error("创建PAP操作日志失败:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "创建PAP操作日志失败"})
		return
	}

	// 提交事务
	if err := tx.Commit().Error; err != nil {
		global.Logger.Error("提交事务失败:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "操作失败"})
		return
	}

	global.Logger.Info("用户PAP增加成功, 用户ID:", req.UserID, "数量:", req.Amount, "来源:", req.Source)

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "PAP增加成功",
		"data":    newBalance,
	})
}

// ConsumeUserPap 消费用户PAP
func ConsumeUserPap(c *gin.Context) {
	var req struct {
		UserID   uint   `json:"userId" binding:"required"`
		Amount   int    `json:"amount" binding:"required"`
		Source   string `json:"source"`
		SourceID uint   `json:"sourceId"`
		Remark   string `json:"remark"`
		Operator uint   `json:"operator"` // 操作人ID
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "参数错误"})
		return
	}

	if req.Amount <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "PAP数量必须大于0"})
		return
	}

	// 获取当前余额
	var currentBalance int
	var latestPap service.CorpPap
	result := global.Db.Where("user_id = ?", req.UserID).Order("created_at DESC").First(&latestPap)
	if result.Error == nil {
		currentBalance = latestPap.Balance
	} else if result.Error != gorm.ErrRecordNotFound {
		global.Logger.Error("获取PAP余额失败:", result.Error)
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "获取PAP余额失败"})
		return
	} else {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "用户PAP余额不足"})
		return
	}

	// 检查余额是否足够
	if currentBalance < req.Amount {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "用户PAP余额不足"})
		return
	}

	// 计算新余额
	newBalance := currentBalance - req.Amount

	// 开始事务
	tx := global.Db.Begin()

	// 创建PAP记录
	papRecord := service.CorpPap{
		UserID:     req.UserID,
		Amount:     -req.Amount, // 负数表示消费
		Balance:    newBalance,
		Source:     req.Source,
		SourceID:   req.SourceID,
		Type:       2, // 2-消费
		CreateTime: time.Now(),
		Remark:     req.Remark,
	}

	if err := tx.Create(&papRecord).Error; err != nil {
		tx.Rollback()
		global.Logger.Error("创建PAP记录失败:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "创建PAP记录失败"})
		return
	}

	// 创建操作日志
	papLog := service.CorpPapLog{
		UserID:     req.UserID,
		Operation:  "消费PAP",
		Amount:     req.Amount,
		BeforeVal:  currentBalance,
		AfterVal:   newBalance,
		Operator:   req.Operator,
		CreateTime: time.Now(),
		Remark:     req.Remark,
	}

	if err := tx.Create(&papLog).Error; err != nil {
		tx.Rollback()
		global.Logger.Error("创建PAP操作日志失败:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "创建PAP操作日志失败"})
		return
	}

	// 提交事务
	if err := tx.Commit().Error; err != nil {
		global.Logger.Error("提交事务失败:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "操作失败"})
		return
	}

	global.Logger.Info("用户PAP消费成功, 用户ID:", req.UserID, "数量:", req.Amount, "来源:", req.Source)

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "PAP消费成功",
		"data":    newBalance,
	})
}

// GetPapLogs 获取PAP操作日志
func GetPapLogs(c *gin.Context) {
	var req struct {
		UserID uint `json:"userId" form:"userId"`
		Page   int  `json:"page" form:"page"`
		Limit  int  `json:"limit" form:"limit"`
	}

	if err := c.ShouldBind(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "参数错误"})
		return
	}

	if req.Page <= 0 {
		req.Page = 1
	}
	if req.Limit <= 0 {
		req.Limit = 10
	}

	var papLogs []service.CorpPapLog
	var total int64

	db := global.Db.Model(&service.CorpPapLog{})
	if req.UserID > 0 {
		db = db.Where("user_id = ?", req.UserID)
	}

	db.Count(&total)
	offset := (req.Page - 1) * req.Limit
	result := db.Order("created_at DESC").Offset(offset).Limit(req.Limit).Find(&papLogs)
	if result.Error != nil {
		global.Logger.Error("获取PAP操作日志失败:", result.Error)
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "获取PAP操作日志失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "获取PAP操作日志成功",
		"data": gin.H{
			"total": total,
			"items": papLogs,
		},
	})
}
