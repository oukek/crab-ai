package _type

type Request struct {
	Messages []Message `json:"messages,omitempty"`
	// 模型名
	Model          Model           `json:"model"`
	ResponseFormat *ResponseFormat `json:"response_format,omitempty"`
	// 是否使用流式
	Stream      *bool    `json:"stream,omitempty"`
	MaxTokens   *int     `json:"max_tokens,omitempty"`
	Temperature *float64 `json:"temperature,omitempty"`
	Tools       []Tool   `json:"tools,omitempty"`
	ToolChoice  any      `json:"tool_choice,omitempty"`

	// 下面这些基本不用
	Stop              []string        `json:"stop,omitempty"`
	TopP              *float64        `json:"top_p,omitempty"`
	TopK              *int            `json:"top_k,omitempty"`
	FrequencyPenalty  *float64        `json:"frequency_penalty,omitempty"`
	PresencePenalty   *float64        `json:"presence_penalty,omitempty"`
	RepetitionPenalty *float64        `json:"repetition_penalty,omitempty"`
	Seed              *int            `json:"seed,omitempty"`
	LogitBias         map[int]float64 `json:"logit_bias,omitempty"`
}

type Message struct {
	Role    Role   `json:"role"`
	Content string `json:"content"`
	Name    string `json:"name,omitempty"`
}

type Role string

const (
	UserRole      Role = "user"
	AssistantRole Role = "assistant"
	SystemRole    Role = "system"
	ToolRole      Role = "tool"
)

type ResponseFormat struct {
	Type string `json:"type"`
}

func WithJSONResponseFormat() *ResponseFormat {
	return &ResponseFormat{Type: "json_object"}
}

func WithStream() *bool {
	var isStream = true
	return &isStream
}

func WithMaxToken(maxToken int) *int {
	if maxToken == 0 {
		return nil
	}
	return &maxToken
}

func WithTemperature(temperature float64) *float64 {
	return &temperature
}
