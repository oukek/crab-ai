package llm

import (
	"context"
	"net/http"
	"oukek/crab-ai/app/model"
	_type "oukek/crab-ai/app/service/llm/type"
)

// RequestParameters 包含构建 HTTP 请求所需的信息
type RequestParameters struct {
	URL     string
	Method  string
	Headers map[string]string
	Body    []byte
}

// LLMProvider 定义了 LLM 服务商需要实现的接口
type LLMProvider interface {
	// MakeChatRequest 将通用的聊天请求转换为服务商特定的请求参数
	MakeChatRequest(ctx context.Context, provider *model.Provider, commonRequest *_type.ChatRequest) (*RequestParameters, error)

	// ParseChatResponse 将服务商的响应解析为通用格式
	// responseBody 是从 HTTP 响应中读取的原始 []byte
	ParseChatResponse(ctx context.Context, responseBody []byte) (*_type.Response, error)

	// MakeStreamChatRequest 将通用的流式聊天请求转换为服务商特定的请求参数
	// 返回值同 MakeChatRequest
	MakeStreamChatRequest(ctx context.Context, provider *model.Provider, commonRequest *_type.ChatRequest) (*RequestParameters, error)

	// ParseStreamEvent 解析流式响应中的单个事件
	// eventData 是从 SSE 事件中提取的数据部分 (去除 "data: ")
	// Returns:
	// 1. The parsed stream item (_type.StreamItem) or nil if it's not a data event.
	// 2. A boolean indicating if the event signals the end of the stream ("[DONE]").
	// 3. An error if parsing fails.
	ParseStreamEvent(ctx context.Context, eventData []byte) (item *_type.StreamItem, isDone bool, err error)

	// MakeEmbeddingsRequest 转换嵌入请求
	MakeEmbeddingsRequest(ctx context.Context, provider *model.Provider, commonRequest *_type.EmbeddingsRequest) (*RequestParameters, error)

	// ParseEmbeddingsResponse 解析嵌入响应
	ParseEmbeddingsResponse(ctx context.Context, responseBody []byte) (*_type.EmbeddingsResponse, error)

	// HandleError (可选) 处理或转换服务商特定的错误
	// inputErr 可能是 HTTP 错误或解析错误
	// responseBody 可能是 nil 如果错误发生在读取 body 之前
	// 返回转换后的错误，如果无法处理则返回原始错误
	HandleError(ctx context.Context, inputErr error, response *http.Response, responseBody []byte) error
}
