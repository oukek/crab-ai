package _type

type Response struct {
	ID      string               `json:"id" mapstructure:"id"`
	Choices []NonStreamingChoice `json:"choices" mapstructure:"choices"`
	Created int64                `json:"created" mapstructure:"created"`
	Model   string               `json:"model" mapstructure:"model"`
	// chat.completion.chunk chat.completion
	Object string `json:"object" mapstructure:"object"`
	Usage  *Usage `json:"usage" mapstructure:"usage"`

	// azure才有，表示prompt的敏感检测结果
	PromptFilterResults []PromptFilterResults `json:"prompt_filter_results" mapstructure:"prompt_filter_results"`
}

// StreamItem represents a single chunk in a streaming response.
type StreamItem struct {
	Id      string            `json:"id" mapstructure:"id"` // Note: Field name is Id, json tag is id
	Choices []StreamingChoice `json:"choices" mapstructure:"choices"`
	Created int64             `json:"created" mapstructure:"created"`
	Model   string            `json:"model" mapstructure:"model"`
	Object  string            `json:"object" mapstructure:"object"`         // e.g., "chat.completion.chunk"
	Usage   *Usage            `json:"usage,omitempty" mapstructure:"usage"` // Often nil until the last chunk
	Error   *Error            `json:"error,omitempty"`                      // For potential error messages within the stream
	// Azure specific prompt filter result
	PromptFilterResults []PromptFilterResults `json:"prompt_filter_results,omitempty" mapstructure:"prompt_filter_results"`
}

// IsValid checks if the stream item is valid and handles special cases like Azure prompt filters or embedded errors.
func (si *StreamItem) IsValid() error {
	// Azure returns prompt filter results in the first chunk without an ID.
	if len(si.PromptFilterResults) > 0 && si.Id == "" {
		// This is a valid Azure prompt filter result chunk, treat as valid but maybe skip processing content.
		return nil
	}
	// Check if it's an error message embedded in the stream
	if si.Error != nil {
		// This is an error chunk
		return si.Error
	}
	// For a regular data chunk, ID and Choices are expected.
	// Allow chunks with nil/empty choices (like finish chunks) but require Id.
	if si.Id == "" {
		// Use the helper function to create a standard error
		return NewError("", "Stream item missing ID", "invalid_stream_data", nil)
	}
	// No error found
	return nil
}

// Embedding represents a single embedding vector.
type Embedding struct {
	Index     int       `json:"index"`
	Object    string    `json:"object"` // e.g., "embedding"
	Embedding []float64 `json:"embedding"`
}

// EmbeddingsResponse represents the response structure for embedding requests.
type EmbeddingsResponse struct {
	Object string      `json:"object"` // e.g., "list"
	Data   []Embedding `json:"data"`
	Model  string      `json:"model"`
	Usage  Usage       `json:"usage"`
}

type Usage struct {
	CompletionTokens int `json:"completion_tokens" mapstructure:"completion_tokens"`
	PromptTokens     int `json:"prompt_tokens" mapstructure:"prompt_tokens"`
	TotalTokens      int `json:"total_tokens" mapstructure:"total_tokens"`
	TotalCost        int `json:"total_cost" mapstructure:"total_cost"`
}

// 过滤
type ContentFilterResults struct {
	Hate struct {
		Filtered bool   `json:"filtered" mapstructure:"filtered"`
		Severity string `json:"severity" mapstructure:"severity"`
	} `json:"hate" mapstructure:"hate"`
	SelfHarm struct {
		Filtered bool   `json:"filtered" mapstructure:"filtered"`
		Severity string `json:"severity" mapstructure:"severity"`
	} `json:"self_harm" mapstructure:"self_harm"`
	Sexual struct {
		Filtered bool   `json:"filtered" mapstructure:"filtered"`
		Severity string `json:"severity" mapstructure:"severity"`
	} `json:"sexual" mapstructure:"sexual"`
	Violence struct {
		Filtered bool   `json:"filtered" mapstructure:"filtered"`
		Severity string `json:"severity" mapstructure:"severity"`
	} `json:"violence" mapstructure:"violence"`
}

// PromptFilterResults prompt敏感检测
type PromptFilterResults struct {
	//prompt_index
	PromptIndex int `json:"prompt_index" mapstructure:"prompt_index"`
	// content_filter_results
	ContentFilterResults *ContentFilterResults `json:"content_filter_results" mapstructure:"content_filter_results"`
}

type NonStreamingChoice struct {
	Index        int         `json:"index" mapstructure:"index"`
	FinishReason string      `json:"finish_reason" mapstructure:"finish_reason"`
	Message      *ResMessage `json:"message" mapstructure:"message"`
	// azure才有
	ContentFilterResults *ContentFilterResults `json:"content_filter_results" mapstructure:"content_filter_results"`
}

type ResMessage struct {
	Content      string        `json:"content" mapstructure:"content"`
	Role         string        `json:"role" mapstructure:"role"`
	ToolCalls    []ToolCall    `json:"tool_calls" mapstructure:"tool_calls"`
	FunctionCall *FunctionCall `json:"function_call" mapstructure:"function_call"`
}

type StreamingChoice struct {
	Index        int    `json:"index" mapstructure:"index"` // Add index to streaming choice
	FinishReason string `json:"finish_reason,omitempty" mapstructure:"finish_reason"`
	Delta        *Delta `json:"delta" mapstructure:"delta"`
	// azure才有
	ContentFilterResults *ContentFilterResults `json:"content_filter_results,omitempty" mapstructure:"content_filter_results"`
}

type Delta struct {
	Content      string        `json:"content,omitempty" mapstructure:"content"`
	Role         string        `json:"role,omitempty" mapstructure:"role"`
	ToolCalls    []ToolCall    `json:"tool_calls,omitempty" mapstructure:"tool_calls"`
	FunctionCall *FunctionCall `json:"function_call,omitempty" mapstructure:"function_call"`
}

type ToolCall struct {
	ID       string       `json:"id" mapstructure:"id"`
	Type     string       `json:"type" mapstructure:"type"`
	Function FunctionCall `json:"function" mapstructure:"function"`
}

type FunctionCall struct {
	Name      string `json:"name" mapstructure:"name"`
	Arguments string `json:"arguments" mapstructure:"arguments"`
}
