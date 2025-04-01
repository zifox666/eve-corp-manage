package esi

import (
	"net/url"
)

// getServerStatus 获取服务器状态
func getServerStatus() (map[string]interface{}, error) {
	var result map[string]interface{}

	query := url.Values{}
	query.Set("datasource", "tranquility")

	err := Client.GetJSON("/status/", query, &result)
	if err != nil {
		return nil, err
	}

	return result, nil
}
