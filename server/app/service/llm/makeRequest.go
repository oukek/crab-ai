package llm

import (
	"encoding/json"
	"net/http"
	"oukek/crab-ai/app/common"
	_type "oukek/crab-ai/app/service/llm/type"
	"strings"
)

// 根据不同的模型，将其转为json数据
func makeRequest(chains *ChatChains) error {
	var req *http.Request
	isStream := chains.Req.Stream != nil && *chains.Req.Stream
	switch chains.Provider {
	case _type.ProviderAIStudio:
		geminiReq, err := _type.NewGeminiReqFromRequest(chains.Req)
		if err != nil {
			return err
		}
		data, err := json.Marshal(geminiReq)
		if err != nil {
			return err
		}
		if isStream {
			req, err = http.NewRequest("POST", chains.Host+"/v1beta/models/"+string(chains.Req.Model)+":streamGenerateContent?alt=sse&key="+chains.ApiKey, strings.NewReader(string(data)))
		} else {
			req, err = http.NewRequest("POST", chains.Host+"/v1beta/models/"+string(chains.Req.Model)+":generateContent?key="+chains.ApiKey, strings.NewReader(string(data)))
		}
		if err != nil {
			common.BaseLogger.WithError(err).Error("GPT：请求接口失败")
			return err
		}
	default:
		chains.Set("reqBody", chains.Req)
	}

	// 设置请求头的Content-Type为application/json
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/122.0.0.0 Safari/537.36")
	if isStream {
		req.Header.Add("Accept", "text/event-stream")
	}

	chains.Set("req", req)
	return nil
}
