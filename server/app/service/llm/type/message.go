package _type

import "errors"

type ContentType string

const (
	TextContentType  ContentType = "text"
	ImageContentType ContentType = "image_url"
)

type ContentPart struct {
	Type     ContentType `json:"type"`
	Text     string      `json:"text,omitempty"`
	ImageURL *struct {
		URL    string `json:"url"`
		Detail string `json:"detail,omitempty"`
	} `json:"image_url,omitempty"`
}

type Content interface {
	string | []ContentPart
}

type StreamItemChoiceDelta struct {
	Content string `json:"content"`
}

type StreamItemChoice struct {
	Delta        StreamItemChoiceDelta `json:"delta"`
	Index        int                   `json:"index"`
	FinishReason string                `json:"finish_reason"`
}

type StreamItem struct {
	Id      string             `json:"id"`
	Object  string             `json:"object"`
	Created int                `json:"created"`
	Model   string             `json:"model"`
	Choices []StreamItemChoice `json:"choices"`
}

// 是否是合法的响应
func (r *StreamItem) IsValid() error {
	if len(r.Choices) == 0 {
		return errors.New("gpt stream解析流失败")
	}
	if r.Choices[0].FinishReason == "content_filter" {
		return errors.New("内容存在敏感内容，已被过滤，请不要生成违法内容。如有疑问请咨询客服")
	}
	return nil
}
