package llm

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"oukek/crab-ai/app/common"
	"oukek/crab-ai/app/model"
	_type "oukek/crab-ai/app/service/llm/type"
)

// GoogleAIStudioProvider 实现了 LLMProvider 接口，用于 Google AI Studio 服务
type GoogleAIStudioProvider struct{}

// NewGoogleAIStudioProvider 创建一个新的 GoogleAIStudioProvider 实例
func NewGoogleAIStudioProvider() *GoogleAIStudioProvider {
	return &GoogleAIStudioProvider{}
}

// makeRequest 创建通用的请求参数
func (g *GoogleAIStudioProvider) makeRequest(ctx context.Context, provider *model.Provider, commonRequest *_type.ChatRequest, stream bool) (*RequestParameters, error) {
	geminiReq, err := _type.NewGeminiReqFromChatRequest(commonRequest)
	if err != nil {
		return nil, fmt.Errorf("创建 Gemini 请求失败: %w", err)
	}
	data, err := json.Marshal(geminiReq)
	if err != nil {
		return nil, fmt.Errorf("序列化 Gemini 请求失败: %w", err)
	}

	url := provider.BaseUrl + "/v1beta/models/" + string(commonRequest.Model)
	if stream {
		url += ":streamGenerateContent?alt=sse&key=" + provider.ApiKey
	} else {
		url += ":generateContent?key=" + provider.ApiKey
	}

	headers := map[string]string{
		"Content-Type": "application/json",
		"User-Agent":   "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/122.0.0.0 Safari/537.36",
	}
	if stream {
		headers["Accept"] = "text/event-stream"
	}

	return &RequestParameters{
		URL:     url,
		Method:  "POST",
		Headers: headers,
		Body:    data,
	}, nil
}

// MakeChatRequest 实现 LLMProvider 接口
func (g *GoogleAIStudioProvider) MakeChatRequest(ctx context.Context, provider *model.Provider, commonRequest *_type.ChatRequest) (*RequestParameters, error) {
	return g.makeRequest(ctx, provider, commonRequest, false)
}

// ParseChatResponse 实现 LLMProvider 接口
func (g *GoogleAIStudioProvider) ParseChatResponse(ctx context.Context, responseBody []byte) (*_type.Response, error) {
	var geminiRes _type.GeminiRes
	var rawMap map[string]any

	// 尝试解析为成功的 GeminiRes 结构
	if err := json.Unmarshal(responseBody, &geminiRes); err == nil && len(geminiRes.Candidates) > 0 {
		// Retrieve model name from context or another source if needed by ToResponse

		return geminiRes.ToResponse()
	}

	// 如果直接解析 GeminiRes 失败，尝试解析为 map[string]any 以便检查错误结构
	if err := json.Unmarshal(responseBody, &rawMap); err != nil {
		common.BaseLogger.WithError(err).WithField("body", string(responseBody)).Error("GoogleAIStudio: 解析响应体 JSON 失败")
		return nil, fmt.Errorf("无法解析响应体: %w", err)
	}

	// Handle error via parseGoogleAPIError helper
	apiErr := g.parseGoogleAPIError(rawMap)
	if apiErr != nil {
		return nil, apiErr
	}

	common.BaseLogger.WithField("body", string(responseBody)).Warn("GoogleAIStudio: 未知的响应格式")
	return nil, fmt.Errorf("未知的 Google AI Studio 响应格式: %s", string(responseBody))
}

// MakeStreamChatRequest 实现 LLMProvider 接口
func (g *GoogleAIStudioProvider) MakeStreamChatRequest(ctx context.Context, provider *model.Provider, commonRequest *_type.ChatRequest) (*RequestParameters, error) {
	// Pass model name into context for potential use in ParseStreamEvent
	ctx = context.WithValue(ctx, "modelName", commonRequest.Model)
	return g.makeRequest(ctx, provider, commonRequest, true)
}

// ParseStreamEvent 实现 LLMProvider 接口
func (g *GoogleAIStudioProvider) ParseStreamEvent(ctx context.Context, eventData []byte) (item *_type.StreamItem, isDone bool, err error) {
	// Google AI Studio's stream does not use [DONE] marker, it just closes.
	// It sends JSON objects per event.

	if len(eventData) == 0 {
		// Empty event, might happen, ignore.
		return nil, false, nil
	}

	var geminiStreamItem _type.GeminiStreamItem
	var rawMap map[string]any

	if err := json.Unmarshal(eventData, &geminiStreamItem); err == nil {
		// Retrieve model name from context
		modelName := "" // Default
		if modelVal := ctx.Value("modelName"); modelVal != nil {
			modelName = modelVal.(string)
		}

		streamItem := geminiStreamItem.ToStreamItem(modelName)
		// Check if the converted item contains an error
		if streamItem.Error != nil {
			return nil, false, streamItem.Error // Return the error parsed by ToStreamItem
		}
		return streamItem, false, nil // Not done, return the item
	}

	// If direct parsing fails, check for error structure
	if errMap := json.Unmarshal(eventData, &rawMap); errMap != nil {
		// Not a valid JSON object at all
		common.BaseLogger.WithError(errMap).WithField("data", string(eventData)).Error("GoogleAIStudio: 解析流式事件 JSON 失败")
		return nil, false, fmt.Errorf("解析流式事件失败: %w", errMap)
	}

	// Check if it's a Google API error
	apiErr := g.parseGoogleAPIError(rawMap)
	if apiErr != nil {
		// It's an error event
		return nil, false, apiErr // Return the specific API error
	}

	// It's some other JSON structure we don't recognize in the stream
	common.BaseLogger.WithField("data", string(eventData)).Warn("GoogleAIStudio: 未知的流式事件格式")
	// Return a structured error instead of just fmt.Errorf
	return nil, false, _type.NewError("", fmt.Sprintf("未知的 Google AI Studio 流式事件格式: %s", string(eventData)), "invalid_stream_data", nil)
}

// MakeEmbeddingsRequest 实现 LLMProvider 接口 (Not Implemented)
func (g *GoogleAIStudioProvider) MakeEmbeddingsRequest(ctx context.Context, provider *model.Provider, commonRequest *_type.EmbeddingsRequest) (*RequestParameters, error) {
	return nil, _type.ErrMethodNotImplemented
}

// ParseEmbeddingsResponse 实现 LLMProvider 接口 (Not Implemented)
func (g *GoogleAIStudioProvider) ParseEmbeddingsResponse(ctx context.Context, responseBody []byte) (*_type.EmbeddingsResponse, error) {
	return nil, _type.ErrMethodNotImplemented
}

// HandleError 实现 LLMProvider 接口
func (g *GoogleAIStudioProvider) HandleError(ctx context.Context, inputErr error, response *http.Response, responseBody []byte) error {
	// 1. Handle HTTP level errors passed in inputErr (e.g., network issues)
	if inputErr != nil {
		// If it's a known HTTP error code, potentially wrap it
		if response != nil && response.StatusCode == http.StatusTooManyRequests {
			common.BaseLogger.Warn("GoogleAIStudio: 收到 HTTP 429 Too Many Requests")
			return _type.ErrResourceExhausted // Wrap as specific error type
		}
		// Add handling for other relevant status codes if needed

		// Otherwise, return the original http/network error
		return inputErr // Return the original error (could be network error, etc.)
	}

	// 2. Handle API level errors parsed from responseBody (if available)
	if response != nil && response.StatusCode >= 400 {
		// Attempt to parse API error from body if we have one
		if responseBody != nil {
			var rawMap map[string]any
			// Use tolerant JSON parsing if body might not be JSON
			if err := json.Unmarshal(responseBody, &rawMap); err == nil {
				apiErr := g.parseGoogleAPIError(rawMap)
				if apiErr != nil {
					return apiErr // Return the parsed API error
				}
			}
			// else: body is not JSON or not a Google error structure, fall through
		}
		// If body couldn't be parsed or didn't contain a known error, return a generic HTTP error
		return _type.NewError("", fmt.Sprintf("HTTP Error: %d %s", response.StatusCode, http.StatusText(response.StatusCode)), "http_error", response.StatusCode)
	}

	// If no error identified (inputErr is nil, status code < 400)
	return nil // Indicate success or that the error wasn't handled here
}

// parseGoogleAPIError Helper function to parse standard Google API error structure
func (g *GoogleAIStudioProvider) parseGoogleAPIError(rawMap map[string]any) error {
	if errObj, ok := rawMap["error"].(map[string]interface{}); ok {
		codeF, _ := errObj["code"].(float64) // JSON 数字默认为 float64
		code := int(codeF)
		message, _ := errObj["message"].(string)
		status, _ := errObj["status"].(string)
		// Use the standard Error struct
		apiError := _type.NewError("", message, status, code)
		common.BaseLogger.WithField("error", rawMap["error"]).Error(apiError.Error())

		// Check for specific known error statuses
		if status == "RESOURCE_EXHAUSTED" {
			return _type.ErrResourceExhausted // Return the specific error variable
		}
		// Return the structured API error
		return apiError
	}
	return nil // Not a Google API error structure
}

/*
// 原有的 makeRequestGoogleAIStudio 函数，现在已被整合到 GoogleAIStudioProvider 的方法中
func makeRequestGoogleAIStudio(chains any) error {
	c := chains.(*ChatChains) // 断言类型
	isStream := c.Req.Stream != nil && *c.Req.Stream
	geminiReq, err := _type.NewGeminiReqFromRequest(c.Req)
	if err != nil {
		return err
	}
	data, err := json.Marshal(geminiReq)
	if err != nil {
		return err
	}
	provider := c.Provider
	var req *http.Request
	if isStream {
		req, err = http.NewRequest("POST", provider.BaseUrl+"/v1beta/models/"+string(c.Req.Model)+":streamGenerateContent?alt=sse&key="+provider.ApiKey, strings.NewReader(string(data)))
	} else {
		req, err = http.NewRequest("POST", provider.BaseUrl+"/v1beta/models/"+string(c.Req.Model)+":generateContent?key="+provider.ApiKey, strings.NewReader(string(data)))
	}
	if err != nil {
		common.BaseLogger.WithError(err).Error("GPT：请求接口失败")
		return err
	}
	// 设置请求头
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/122.0.0.0 Safari/537.36")
	if isStream {
		req.Header.Add("Accept", "text/event-stream")
	}
	c.Set("req", req)
	return nil
}
*/
