package llm

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"

	"oukek/crab-ai/app/common"
	_type "oukek/crab-ai/app/service/llm/type"
)

func Chat(chains *ChatChains) error {
	defaultError := errors.New("服务器异常，请联系客服")

	req, ok := chains.Get("req")
	if !ok {
		return errors.New("GPT：请求体不存在")
	}

	var body []byte
	// 发送请求并获取响应
	resp, err := http.DefaultClient.Do(req.(*http.Request))
	if err != nil {
		common.BaseLogger.WithError(err).Error("GPT：发送请求失败")
		return defaultError
	}
	// 延迟关闭响应体
	defer resp.Body.Close()

	// 读取响应体的内容
	body, err = io.ReadAll(resp.Body)
	if err != nil {
		common.BaseLogger.WithError(err).Error("GPT：读取body失败")
		return defaultError
	}
	// 判断状态码
	if resp.StatusCode != http.StatusOK {
		// 429 代表请求过于频繁
		if resp.StatusCode == http.StatusTooManyRequests {
			common.BaseLogger.WithField("body", string(body)).WithField("status", resp.StatusCode).Error("GPT：请求过于频繁")
			return _type.ErrResourceExhausted
		}
		common.BaseLogger.WithField("body", string(body)).WithField("status", resp.StatusCode).Error("GPT：请求失败")
		return defaultError
	}

	res := map[string]any{}
	err = json.Unmarshal(body, &res)
	if err != nil {
		common.BaseLogger.WithError(err).WithField("res", string(body)).Error("GPT：转换成结构体失败")
		return defaultError
	}

	// 获取object，如果存在的话，则判断是不是error
	if object, ok := res["object"]; ok {
		if object == "error" {
			// 获取message，如果存在的话，则打日志的时候加上，否则为空
			errMsg := ""
			if message, ok := res["message"]; ok {
				errMsg = message.(string)
			}
			common.BaseLogger.WithField("res", string(body)).Error("GPT：请求失败:" + errMsg)
			return defaultError
		}
	}
	// 获取error，如果存在的话，则说明也有错误，一般是azure的错误
	if resError, ok := res["error"]; ok {
		// 尝试获取message
		errMsg := ""
		if message, ok := resError.(map[string]interface{})["message"]; ok {
			errMsg = message.(string)
		}
		// 打印错误日志
		common.BaseLogger.WithField("res", string(body)).Error("GPT：请求失败:" + errMsg)
		return defaultError
	}
	chains.Set("res", res)

	return nil
}
