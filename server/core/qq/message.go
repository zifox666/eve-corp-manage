package qq

import (
	"encoding/json"
)

// MessageSegment 表示一个消息段
type MessageSegment struct {
	Type string                 `json:"type"`
	Data map[string]interface{} `json:"data"`
}

// Message 表示一个完整的消息，由多个消息段组成
type Message []MessageSegment

// NewMessage 创建一个新的消息
func NewMessage() Message {
	return Message{}
}

// Text 添加文本消息段
func (m Message) Text(text string) Message {
	segment := MessageSegment{
		Type: "text",
		Data: map[string]interface{}{
			"text": text,
		},
	}
	return append(m, segment)
}

// At 添加@用户消息段
func (m Message) At(qq string) Message {
	segment := MessageSegment{
		Type: "at",
		Data: map[string]interface{}{
			"qq": qq,
		},
	}
	return append(m, segment)
}

// Image 添加图片消息段
func (m Message) Image(file string) Message {
	segment := MessageSegment{
		Type: "image",
		Data: map[string]interface{}{
			"file": file,
		},
	}
	return append(m, segment)
}

// Face 添加表情消息段
func (m Message) Face(id string) Message {
	segment := MessageSegment{
		Type: "face",
		Data: map[string]interface{}{
			"id": id,
		},
	}
	return append(m, segment)
}

// Reply 添加回复消息段
func (m Message) Reply(id string) Message {
	segment := MessageSegment{
		Type: "reply",
		Data: map[string]interface{}{
			"id": id,
		},
	}
	return append(m, segment)
}

// ToJSON 将消息转换为JSON字符串
func (m Message) ToJSON() (string, error) {
	bytes, err := json.Marshal(m)
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}

// ToString 将消息转换为纯文本（提取所有文本消息段）
func (m Message) ToString() string {
	var text string
	for _, segment := range m {
		if segment.Type == "text" {
			if textValue, ok := segment.Data["text"].(string); ok {
				text += textValue
			}
		}
	}
	return text
}
