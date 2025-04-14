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

type StreamResponse struct {
	ID      string            `json:"id"`
	Choices []StreamingChoice `json:"choices"`
	Created int64             `json:"created"`
	Model   string            `json:"model"`
	Object  string            `json:"object"`
	Usage   *Usage            `json:"usage"`
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
	FinishReason string `json:"finish_reason"`
	Delta        *Delta `json:"delta"`
}

type Delta struct {
	Content      string        `json:"content"`
	Role         string        `json:"role"`
	ToolCalls    []ToolCall    `json:"tool_calls"`
	FunctionCall *FunctionCall `json:"function_call"`
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
