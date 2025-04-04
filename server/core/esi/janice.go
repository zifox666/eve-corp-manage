package esi

import (
	"encoding/json"
	"eve-corp-manager/global"
	"fmt"
	"io"
)

// GetAppraisal 从Janice获取物品估价
func GetAppraisal(queryStr string) (float64, error) {
	if queryStr == "" {
		return 0, nil
	}

	// 准备请求体
	reqBody := map[string]interface{}{
		"market_name":    "jita",
		"pricelist_name": "default",
		"pricedata": map[string]interface{}{
			"pricing_type":     "immediate",
			"pricing_modifier": 0,
			"evepraisal_url":   "",
			"raw_textarea":     queryStr,
			"live_update":      false,
		},
	}

	reqJSON, err := json.Marshal(reqBody)
	if err != nil {
		return 0, err
	}

	// 发送POST请求
	resp, err := JaniceClient.Post("/appraisal", "application/json", reqJSON)
	if err != nil {
		return 0, err
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			global.Logger.Errorf("Failed to close response body: %v", err)
		}
	}(resp.Body)

	if resp.StatusCode >= 400 {
		return 0, fmt.Errorf("janice API错误 (状态码: %d)", resp.StatusCode)
	}

	// 解析响应
	var result struct {
		Appraisal struct {
			Prices struct {
				Sell struct {
					Min float64 `json:"min"`
				} `json:"sell"`
			} `json:"prices"`
		} `json:"appraisal"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return 0, err
	}

	return result.Appraisal.Prices.Sell.Min, nil
}
