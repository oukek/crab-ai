package _type

import (
	"fmt"
	"strings"
	"time"
)

type GeminiReq struct {
	Contents         []GeminiReqContent       `json:"contents"`
	SafetySettings   []GeminiReqSafetySetting `json:"safetySettings,omitempty"`
	GenerationConfig *GenerationConfig        `json:"generationConfig,omitempty"`

	SystemInstruction *GeminiReqContent `json:"systemInstruction,omitempty"`
}

type GenerationConfig struct {
	StopSequences    []string `json:"stopSequences,omitempty"`
	ResponseMimeType *string  `json:"responseMimeType,omitempty"`
	CandidateCount   *int     `json:"candidateCount,omitempty"`
	MaxOutputTokens  *int     `json:"maxOutputTokens,omitempty"`
	Temperature      *float64 `json:"temperature,omitempty"`
	TopP             *float64 `json:"topP,omitempty"`
	TopK             *int     `json:"topK,omitempty"`
}

type GeminiReqContent struct {
	Role  string                 `json:"role,omitempty"`
	Parts []GeminiReqContentPart `json:"parts"`
}

type GeminiReqContentPart struct {
	Text string `json:"text"`
}

type GeminiReqSafetySetting struct {
	Category  string `json:"category"`
	Threshold string `json:"threshold"`
}

func WithResponseMimeTypeText() *string {
	return nil
}

// 返回json有两种方式，一种是在prompt中提示，一种是在response_format中提示
// 我们全部都使用prompt中提示，这样比较容易，功能也一样
func WithResponseMimeTypeJson() *string {
	a := "application/json"
	return &a
}
func WithSystemPrompt(prompt string) *GeminiReqContent {
	return &GeminiReqContent{
		Parts: []GeminiReqContentPart{
			{
				Text: prompt,
			},
		},
	}

}

func WithSafetySettingNoBlock() []GeminiReqSafetySetting {
	return []GeminiReqSafetySetting{
		{
			Category:  "HARM_CATEGORY_HATE_SPEECH",
			Threshold: "BLOCK_NONE",
		},
		{
			Category:  "HARM_CATEGORY_SEXUALLY_EXPLICIT",
			Threshold: "BLOCK_NONE",
		},
		{
			Category:  "HARM_CATEGORY_DANGEROUS_CONTENT",
			Threshold: "BLOCK_NONE",
		},
		{
			Category:  "HARM_CATEGORY_HARASSMENT",
			Threshold: "BLOCK_NONE",
		},
		{
			Category:  "HARM_CATEGORY_CIVIC_INTEGRITY",
			Threshold: "BLOCK_NONE",
		},
	}
}

// NewGeminiReqFromChatRequest converts a generic ChatRequest to a Gemini specific request.
func NewGeminiReqFromChatRequest(req *ChatRequest) (*GeminiReq, error) {
	geminiReq := &GeminiReq{
		Contents:       []GeminiReqContent{},
		SafetySettings: WithSafetySettingNoBlock(),
		GenerationConfig: &GenerationConfig{
			MaxOutputTokens: req.MaxTokens,
			Temperature:     req.Temperature,
			TopP:            req.TopP,
		},
	}
	// 如果是json的话，则设置返回格式为json
	if req.ResponseFormat != nil && req.ResponseFormat.Type == "json_object" {
		geminiReq.GenerationConfig.ResponseMimeType = WithResponseMimeTypeJson()
	}

	for _, message := range req.Messages {
		role := ""
		switch message.Role {
		case UserRole:
			role = "user"
		case AssistantRole:
			role = "model"
		case SystemRole:
			geminiReq.SystemInstruction = WithSystemPrompt(message.Content)

		default:
			return nil, fmt.Errorf("role %s not supported", message.Role)
		}
		msg := GeminiReqContent{
			Role: role,
			Parts: []GeminiReqContentPart{
				{
					Text: message.Content,
				},
			},
		}
		geminiReq.Contents = append(geminiReq.Contents, msg)
	}
	return geminiReq, nil
}

type GeminiRes struct {
	Candidates     []Candidate    `json:"candidates"`
	PromptFeedback PromptFeedback `json:"promptFeedback"`
	UsageMetadata  UsageMetadata  `json:"usageMetadata"`
	Error          struct {
		Code    int    `json:"code"`
		Message string `json:"message"`
		Status  int    `json:"status"`
	}
}

type Candidate struct {
	Content struct {
		Parts []struct {
			Text string `json:"text"`
		} `json:"parts"`
		Role string `json:"role"`
	} `json:"content"`
	FinishReason  string `json:"finishReason"`
	Index         int    `json:"index"`
	SafetyRatings []struct {
		Category    string `json:"category"`
		Probability string `json:"probability"`
	} `json:"safetyRatings"`
}

type PromptFeedback struct {
	BlockReason   string `json:"blockReason"`
	SafetyRatings []struct {
		Category    string `json:"category"`
		Probability string `json:"probability"`
	} `json:"safetyRatings"`
}

type UsageMetadata struct {
	PromptTokenCount     int `json:"promptTokenCount"`
	CandidatesTokenCount int `json:"candidatesTokenCount"`
	TotalTokenCount      int `json:"totalTokenCount"`
}

// 将geminiRes转为response
func (g *GeminiRes) ToResponse() (*Response, error) {
	// 创建一个新的Response实例
	res := &Response{}

	if g.PromptFeedback.BlockReason == "OTHER" {
		return nil, ErrPromptBlocked
	}

	// 将GeminiRes中的Candidates字段转换为Response中的Choices字段
	for _, candidate := range g.Candidates {
		choice := NonStreamingChoice{
			Index:        candidate.Index,
			FinishReason: candidate.FinishReason,
			Message: &ResMessage{
				Content: candidate.Content.Parts[0].Text,
				Role:    candidate.Content.Role,
			},
		}
		res.Choices = append(res.Choices, choice)
	}

	// 将GeminiRes中的UsageMetadata字段转换为Response中的Usage字段
	res.Usage = &Usage{
		CompletionTokens: g.UsageMetadata.CandidatesTokenCount,
		PromptTokens:     g.UsageMetadata.PromptTokenCount,
		TotalTokens:      g.UsageMetadata.TotalTokenCount,
	}

	// 返回新的Response实例
	return res, nil
}

type GeminiStreamItem struct {
	Candidates []struct {
		Content struct {
			Parts []struct {
				Text string `json:"text"`
			} `json:"parts"`
			Role string `json:"role"`
		} `json:"content"`
		FinishReason  string `json:"finishReason"`
		Index         int    `json:"index"`
		SafetyRatings []struct {
			Category    string `json:"category"`
			Probability string `json:"probability"`
		} `json:"safetyRatings"`
	} `json:"candidates"`
	Error          *Error          `json:"error,omitempty"`
	PromptFeedback *PromptFeedback `json:"promptFeedback,omitempty"`
	UsageMetadata  *UsageMetadata  `json:"usageMetadata,omitempty"`
}

func (g *GeminiStreamItem) ToStreamItem(modelName string) *StreamItem {
	// 转为通用 stream item
	gptItem := &StreamItem{
		Id:      fmt.Sprintf("gemini-chunk-%d", time.Now().UnixNano()), // Generate pseudo chunk ID
		Object:  "chat.completion.chunk",
		Created: time.Now().Unix(), // Use int64 directly
		Model:   modelName,
	}

	// Handle potential stream-level error or feedback
	if g.Error != nil {
		gptItem.Error = g.Error
		return gptItem // Return immediately as it's an error chunk
	}
	if g.PromptFeedback != nil && g.PromptFeedback.BlockReason != "" {
		gptItem.Error = NewError("", fmt.Sprintf("Prompt blocked: %s", g.PromptFeedback.BlockReason), "prompt_blocked", nil)
		// Map safety ratings if needed
		return gptItem
	}
	// Update usage if present in the chunk
	if g.UsageMetadata != nil {
		gptItem.Usage = &Usage{
			PromptTokens:     g.UsageMetadata.PromptTokenCount,
			CompletionTokens: g.UsageMetadata.CandidatesTokenCount, // Gemini uses candidatesTokenCount
			TotalTokens:      g.UsageMetadata.TotalTokenCount,
		}
	}

	for _, candidate := range g.Candidates {
		var contentParts []string
		for _, part := range candidate.Content.Parts {
			contentParts = append(contentParts, part.Text)
		}

		// Use the types from response.go
		choice := StreamingChoice{ // Use StreamingChoice
			Delta: &Delta{ // Use Delta
				Content: strings.Join(contentParts, ""),
				Role:    "assistant", // Map role if available and needed (g.Content.Role?)
			},
			Index:        candidate.Index,
			FinishReason: candidate.FinishReason,
			// TODO: Map SafetyRatings to ContentFilterResults if needed
		}

		gptItem.Choices = append(gptItem.Choices, choice)
	}

	return gptItem

}
