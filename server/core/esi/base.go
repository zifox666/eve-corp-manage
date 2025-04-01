package esi

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"sync"
	"time"
)

const (
	baseESIURL = "https://esi.evetech.net/latest"
)

// ESIClient 表示与EVE ESI API通信的HTTP客户端
type ESIClient struct {
	client    *http.Client
	userAgent string
	baseURL   string
}

// Client 全局ESI客户端实例
var Client *ESIClient

// NewESIClient 创建一个新的ESI客户端
func NewESIClient(proxyHost, proxyPort, userAgent string) *ESIClient {
	transport := &http.Transport{
		MaxIdleConns:    2000,
		IdleConnTimeout: 90 * time.Second,
	}

	// 设置代理
	if proxyHost != "" && proxyPort != "" {
		proxyURL, err := url.Parse(fmt.Sprintf("http://%s:%s", proxyHost, proxyPort))
		if err == nil {
			transport.Proxy = http.ProxyURL(proxyURL)
		}
	}

	client := &http.Client{
		Transport: transport,
		Timeout:   time.Second * 30,
	}

	return &ESIClient{
		client:    client,
		userAgent: userAgent,
		baseURL:   baseESIURL,
	}
}

// Get 发送GET请求到ESI API
func (c *ESIClient) Get(path string, query url.Values) (*http.Response, error) {
	reqURL := c.baseURL + path
	if query != nil {
		reqURL += "?" + query.Encode()
	}

	req, err := http.NewRequest("GET", reqURL, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("User-Agent", c.userAgent)
	req.Header.Set("Accept", "application/json")

	return c.client.Do(req)
}

// Post 发送POST请求到ESI API
func (c *ESIClient) Post(path string, contentType string, body []byte) (*http.Response, error) {
	reqURL := c.baseURL + path
	req, err := http.NewRequest("POST", reqURL, bytes.NewBuffer(body))
	if err != nil {
		return nil, err
	}

	req.Header.Set("User-Agent", c.userAgent)
	req.Header.Set("Content-Type", contentType)
	req.Header.Set("Accept", "application/json")

	return c.client.Do(req)
}

// GetJSON 发送GET请求并将结果解析为JSON
func (c *ESIClient) GetJSON(path string, query url.Values, result interface{}) error {
	resp, err := c.Get(path, query)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("ESI API错误 (状态码: %d): %s", resp.StatusCode, string(bodyBytes))
	}

	return json.NewDecoder(resp.Body).Decode(result)
}

// AuthorizedGet 发送带授权的GET请求
func (c *ESIClient) AuthorizedGet(path string, query url.Values, token string) (*http.Response, error) {
	reqURL := c.baseURL + path
	if query != nil {
		reqURL += "?" + query.Encode()
	}

	req, err := http.NewRequest("GET", reqURL, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("User-Agent", c.userAgent)
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)

	return c.client.Do(req)
}

// AuthorizedGetJSON 发送带授权的GET请求并解析JSON
func (c *ESIClient) AuthorizedGetJSON(path string, query url.Values, token string, result interface{}) error {
	resp, err := c.AuthorizedGet(path, query, token)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("ESI API错误 (状态码: %d): %s", resp.StatusCode, string(bodyBytes))
	}

	return json.NewDecoder(resp.Body).Decode(result)
}

// GetAllPages 并发获取所有分页数据
func (c *ESIClient) GetAllPages(path string, query url.Values, resultContainer interface{}) error {
	// 首先获取第一页来确定总页数
	resp, err := c.Get(path, query)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("ESI API错误 (状态码: %d): %s", resp.StatusCode, string(bodyBytes))
	}

	// 获取总页数
	totalPagesStr := resp.Header.Get("X-Pages")
	totalPages, err := strconv.Atoi(totalPagesStr)
	if err != nil {
		return fmt.Errorf("解析X-Pages失败: %w", err)
	}

	// 如果只有一页，直接解析并返回
	if totalPages <= 1 {
		return json.NewDecoder(resp.Body).Decode(resultContainer)
	}

	// 解析第一页数据
	var firstPageData []interface{}
	if err := json.NewDecoder(resp.Body).Decode(&firstPageData); err != nil {
		return err
	}

	// 创建一个等待组和结果通道
	var wg sync.WaitGroup
	resultChan := make(chan []interface{}, totalPages)
	errorChan := make(chan error, totalPages)

	// 已经获取了第一页，将其放入结果通道
	resultChan <- firstPageData

	// 限制最大并发数
	semaphore := make(chan struct{}, 100)

	// 并发请求剩余页面
	for page := 2; page <= totalPages; page++ {
		wg.Add(1)
		go func(pageNum int) {
			defer wg.Done()

			semaphore <- struct{}{}
			defer func() { <-semaphore }()

			// 复制查询参数并添加页码
			pageQuery := url.Values{}
			for k, v := range query {
				pageQuery[k] = v
			}
			pageQuery.Set("page", strconv.Itoa(pageNum))

			// 获取当前页数据
			var pageData []interface{}
			err := c.GetJSON(path, pageQuery, &pageData)
			if err != nil {
				errorChan <- fmt.Errorf("获取第%d页失败: %w", pageNum, err)
				return
			}

			resultChan <- pageData
		}(page)
	}

	// 等待所有请求完成
	wg.Wait()
	close(resultChan)
	close(errorChan)

	// 检查是否有错误
	if len(errorChan) > 0 {
		return <-errorChan
	}

	// 合并所有结果
	var allResults []interface{}
	for result := range resultChan {
		allResults = append(allResults, result...)
	}

	// 将合并的结果转换为预期类型
	resultBytes, err := json.Marshal(allResults)
	if err != nil {
		return err
	}
	return json.Unmarshal(resultBytes, resultContainer)
}

// AuthorizedGetAllPages 带授权的并发获取所有分页数据
func (c *ESIClient) AuthorizedGetAllPages(path string, query url.Values, token string, resultContainer interface{}) error {
	// 首先获取第一页来确定总页数
	resp, err := c.AuthorizedGet(path, query, token)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("ESI API错误 (状态码: %d): %s", resp.StatusCode, string(bodyBytes))
	}

	// 获取总页数
	totalPagesStr := resp.Header.Get("X-Pages")
	totalPages, err := strconv.Atoi(totalPagesStr)
	if err != nil {
		return fmt.Errorf("解析X-Pages失败: %w", err)
	}

	// 如果只有一页，直接解析并返回
	if totalPages <= 1 {
		return json.NewDecoder(resp.Body).Decode(resultContainer)
	}

	// 解析第一页数据
	var firstPageData []interface{}
	if err := json.NewDecoder(resp.Body).Decode(&firstPageData); err != nil {
		return err
	}

	// 创建一个等待组和结果通道
	var wg sync.WaitGroup
	resultChan := make(chan []interface{}, totalPages)
	errorChan := make(chan error, totalPages)

	// 已经获取了第一页，将其放入结果通道
	resultChan <- firstPageData

	// 限制最大并发数
	semaphore := make(chan struct{}, 100)

	// 并发请求剩余页面
	for page := 2; page <= totalPages; page++ {
		wg.Add(1)
		go func(pageNum int) {
			defer wg.Done()

			semaphore <- struct{}{}
			defer func() { <-semaphore }()

			// 复制查询参数并添加页码
			pageQuery := url.Values{}
			for k, v := range query {
				pageQuery[k] = v
			}
			pageQuery.Set("page", strconv.Itoa(pageNum))

			// 获取当前页数据
			var pageData []interface{}
			err := c.AuthorizedGetJSON(path, pageQuery, token, &pageData)
			if err != nil {
				errorChan <- fmt.Errorf("获取第%d页失败: %w", pageNum, err)
				return
			}

			resultChan <- pageData
		}(page)
	}

	// 等待所有请求完成
	wg.Wait()
	close(resultChan)
	close(errorChan)

	// 检查是否有错误
	if len(errorChan) > 0 {
		return <-errorChan
	}

	// 合并所有结果
	var allResults []interface{}
	for result := range resultChan {
		allResults = append(allResults, result...)
	}

	// 将合并的结果转换为预期类型
	resultBytes, err := json.Marshal(allResults)
	if err != nil {
		return err
	}
	return json.Unmarshal(resultBytes, resultContainer)
}
