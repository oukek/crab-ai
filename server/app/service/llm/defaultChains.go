package llm

import (
	"oukek/crab-ai/app/model"
	_type "oukek/crab-ai/app/service/llm/type"

	"github.com/gin-gonic/gin"
)

func SimpleChatWithReq(c *gin.Context, req *_type.ChatRequest, provider model.Provider) (*_type.Response, error) {

	chains := NewChatChains()
	chains.Provider = provider
	var res *_type.Response
	var err error
	if req.Stream != nil && *req.Stream {
		res, err = chains.ChatWithStream(c, req)
	} else {
		res, err = chains.Chat(req)
	}

	if err != nil {
		return nil, err
	}
	return res, nil
}
