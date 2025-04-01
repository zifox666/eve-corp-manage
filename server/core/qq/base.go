package qq

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"
)

// Client QQ机器人API客户端
type Client struct {
	BaseURL    string
	HTTPClient *http.Client
}

// 消息类型常量
const (
	MessageTypePrivate = "private" // 私聊消息
	MessageTypeGroup   = "group"   // 群聊消息
)

// QQClient 全局QQ客户端实例
var QQClient *Client

// NewClient 创建新的QQ客户端
func NewClient(baseURL string) *Client {
	return &Client{
		BaseURL: baseURL,
		HTTPClient: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

// SendMessageRequest 发送消息请求
type SendMessageRequest struct {
	MessageType string      `json:"message_type"`
	UserID      string      `json:"user_id,omitempty"`
	GroupID     string      `json:"group_id,omitempty"`
	Message     interface{} `json:"message"` // 可以是字符串或消息段列表
	AutoEscape  bool        `json:"auto_escape"`
}

// APIResponse API响应
type APIResponse struct {
	Status  string      `json:"status"`
	RetCode int         `json:"retcode"`
	Msg     string      `json:"msg"`
	Data    interface{} `json:"data"`
}

// SendPrivateMsg 发送私聊消息
func (c *Client) SendPrivateMsg(userID string, message interface{}, autoEscape bool) (*APIResponse, error) {
	req := SendMessageRequest{
		MessageType: MessageTypePrivate,
		UserID:      userID,
		Message:     message,
		AutoEscape:  autoEscape,
	}
	return c.sendMsg(req)
}

// SendGroupMsg 发送群聊消息
func (c *Client) SendGroupMsg(groupID string, message interface{}, autoEscape bool) (*APIResponse, error) {
	req := SendMessageRequest{
		MessageType: MessageTypeGroup,
		GroupID:     groupID,
		Message:     message,
		AutoEscape:  autoEscape,
	}
	return c.sendMsg(req)
}

// 发送消息的内部方法
func (c *Client) sendMsg(req SendMessageRequest) (*APIResponse, error) {
	if c == nil {
		return nil, errors.New("QQ client not initialized")
	}

	url := fmt.Sprintf("%s/send_msg", c.BaseURL)

	// 序列化请求体
	reqBody, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %v", err)
	}

	// 创建HTTP请求
	httpReq, err := http.NewRequest("POST", url, bytes.NewBuffer(reqBody))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %v", err)
	}
	httpReq.Header.Set("Content-Type", "application/json")

	// 发送请求
	resp, err := c.HTTPClient.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %v", err)
	}
	defer resp.Body.Close()

	// 解析响应
	var apiResp APIResponse
	if err := json.NewDecoder(resp.Body).Decode(&apiResp); err != nil {
		return nil, fmt.Errorf("failed to decode response: %v", err)
	}

	// 检查响应状态
	if apiResp.Status != "ok" {
		return &apiResp, fmt.Errorf("API error: %s (code: %d)", apiResp.Msg, apiResp.RetCode)
	}

	return &apiResp, nil
}
