package esi

import (
	"encoding/json"
	"eve-corp-manager/global"
	"fmt"
	"io"
	"net/url"
	"regexp"
	"strconv"
	"strings"
)

// GetKillmail 获取击毁邮件详情
func GetKillmail(killmailID int, killmailHash string) (map[string]interface{}, error) {
	var result map[string]interface{}

	query := url.Values{}
	query.Set("datasource", "tranquility")

	path := fmt.Sprintf("/killmails/%d/%s/", killmailID, killmailHash)
	err := EsiClient.GetJSON(path, query, &result)
	if err != nil {
		return nil, err
	}

	return result, nil
}

// PostIdsToNames 批量获取ID对应的名称
func PostIdsToNames(ids []int) (map[string]string, error) {
	result := make(map[string]string)
	if len(ids) == 0 {
		return result, nil
	}

	// 过滤掉0值
	filteredIDs := make([]int, 0)
	for _, id := range ids {
		if id != 0 {
			filteredIDs = append(filteredIDs, id)
		}
	}

	if len(filteredIDs) == 0 {
		return result, nil
	}

	// 将整数ID转换为字符串
	idsJSON, err := json.Marshal(filteredIDs)
	if err != nil {
		return nil, err
	}

	// 发送POST请求
	resp, err := EsiClient.Post("/universe/names/", "application/json", idsJSON)
	if err != nil {
		return nil, err
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			global.Logger.Errorf("Failed to close response body: %v", err)
		}
	}(resp.Body)

	if resp.StatusCode >= 400 {
		return nil, fmt.Errorf("ESI API错误 (状态码: %d)", resp.StatusCode)
	}

	// 解析响应
	var namesData []struct {
		ID   int    `json:"id"`
		Name string `json:"name"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&namesData); err != nil {
		return nil, err
	}

	// 构建ID到名称的映射
	for _, item := range namesData {
		result[strconv.Itoa(item.ID)] = item.Name
	}

	return result, nil
}

// ExtractKillID 从zkillboard URL中提取killmail ID
func ExtractKillID(url string) (int, error) {
	pattern := regexp.MustCompile(`https://zkillboard\.com/kill/(\d+)/`)
	match := pattern.FindStringSubmatch(url)
	if match != nil && len(match) > 1 {
		return strconv.Atoi(match[1])
	}
	return 0, fmt.Errorf("无法从URL提取killmail ID")
}

// ExtractKillmailIDAndHash 从ESI URL中提取killmail ID和hash
func ExtractKillmailIDAndHash(url string) (int, string, error) {
	pattern := regexp.MustCompile(`https://esi\.evetech\.net/\w+/killmails/(\d+)/([a-f0-9]+).*`)
	match := pattern.FindStringSubmatch(url)
	if match != nil && len(match) > 2 {
		killmailID, err := strconv.Atoi(match[1])
		if err != nil {
			return 0, "", err
		}
		killmailHash := match[2]
		return killmailID, killmailHash, nil
	}
	return 0, "", fmt.Errorf("无法从URL提取killmail ID和hash")
}

// GetKillmailHash 从各种输入中获取killmail ID和hash
func GetKillmailHash(killmailURL string) (int, string, error) {
	killmailID := 0
	killmailHash := ""
	var err error

	// 尝试将输入解析为纯数字killmail ID
	killmailID, err = strconv.Atoi(strings.TrimSpace(killmailURL))
	if err == nil && killmailID != 0 {
		return killmailID, killmailHash, nil
	}

	// 处理URL
	killmailURL = strings.TrimSpace(killmailURL)
	if strings.HasPrefix(killmailURL, "https://zkillboard.com/kill/") {
		// 从zkillboard URL提取killmail ID
		killmailID, err = ExtractKillID(killmailURL)
		if err != nil || killmailID == 0 {
			return 0, "", fmt.Errorf("killmail不存在")
		}

		// 从zkillboard API获取hash
		_url := fmt.Sprintf("https://zkillboard.com/api/killID/%d/", killmailID)
		resp, err := EsiClient.client.Get(_url)
		if err != nil {
			return 0, "", err
		}
		defer resp.Body.Close()

		var data []struct {
			Zkb struct {
				Hash string `json:"hash"`
			} `json:"zkb"`
		}

		if err := json.NewDecoder(resp.Body).Decode(&data); err != nil || len(data) == 0 {
			return 0, "", fmt.Errorf("无法获取killmail hash")
		}

		killmailHash = data[0].Zkb.Hash
	} else if strings.HasPrefix(killmailURL, "https://esi.evetech.net") {
		// 从ESI URL提取killmail ID和hash
		killmailID, killmailHash, err = ExtractKillmailIDAndHash(killmailURL)
		if err != nil {
			return 0, "", err
		}
	}

	return killmailID, killmailHash, nil
}
