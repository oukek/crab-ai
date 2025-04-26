package llm

import (
	"oukek/crab-ai/app/model"
	"oukek/crab-ai/app/service/chains"
	_type "oukek/crab-ai/app/service/llm/type"

	"github.com/gin-gonic/gin"
)

type ChatChains struct {
	Context *chains.Context[*ChatChains]

	Req *_type.ChatRequest
	Res *_type.Response

	Provider model.Provider `json:"provider"`
}

func (c *ChatChains) Chat(req *_type.ChatRequest) (*_type.Response, error) {
	c.Req = req

	c.Context.Reset(-1)

	c.Context.Use(Chat)

	err := c.Context.Next(c)

	return c.Res, err
}

func (c *ChatChains) ChatWithStream(ctx *gin.Context, req *_type.ChatRequest) (*_type.Response, error) {
	c.Req = req

	c.Context.Reset(-1)

	c.Context.Use(MakeChatWithStream(ctx))

	err := c.Context.Next(c)

	return c.Res, err

}

func (c *ChatChains) Next() error {
	return c.Context.Next(c)
}

// Get retrieves a value from the context.
func (c *ChatChains) Get(key string) (value any, exists bool) {
	return c.Context.Get(key)
}

// GetWithDefault retrieves a value from the context with a default fallback.
func (c *ChatChains) GetWithDefault(key string, defaultValue any) any {
	return c.Context.GetWithDefault(key, defaultValue)
}

// Set stores a value in the context.
func (c *ChatChains) Set(key string, value any) {
	c.Context.Set(key, value)
}

func NewChatChains(handlers ...chains.HandlerFunc[*ChatChains]) *ChatChains {
	ctx := chains.NewContext[*ChatChains]()

	ctx.Use(handlers...)

	return &ChatChains{
		Context: ctx,
	}
}
