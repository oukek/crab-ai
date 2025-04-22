package provider

import (
	"encoding/json"
	"net/http"
	"oukek/crab-ai/app/common"
	"oukek/crab-ai/app/service/llm/type"
	"strings"
)

func makeRequestGoogleAIStudio(chains any) error {
	c := chains.(*type.ChatChains) // 断言类型
	isStream := c.Req.Stream != nil && *c.Req.Stream
	geminiReq, err := type.NewGeminiReqFromRequest(c.Req)
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

// MakeRequest 分发入口
func MakeRequest(chains any) error {
	c := chains.(*type.ChatChains)
	switch c.Provider.Type {
	case "GOOGLE_AI_STUDIO":
		return makeRequestGoogleAIStudio(chains)
	default:
		c.Set("reqBody", c.Req)
		return nil
	}
} 