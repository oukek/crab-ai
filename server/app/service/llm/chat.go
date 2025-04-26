package llm

import (
	"bytes"
	"context"

	// "encoding/json" // 不再需要直接处理 JSON map
	"errors"
	"fmt"
	"io"
	"net/http"

	"oukek/crab-ai/app/common"
)

func Chat(chains *ChatChains) error {
	defaultError := errors.New("服务器异常，请联系客服")

	// 1. 获取 Provider 实例
	llmProvider, err := GetProvider(chains.Provider.Type)
	if err != nil {
		common.BaseLogger.WithError(err).Errorf("获取 Provider 失败: %s", chains.Provider.Type)
		return defaultError // 或者返回更具体的错误
	}

	// 2. 使用 Provider 构建请求参数
	reqParams, err := llmProvider.MakeChatRequest(context.Background(), &chains.Provider, chains.Req)
	if err != nil {
		common.BaseLogger.WithError(err).Error("构建聊天请求失败")
		// 尝试使用 Provider 的错误处理
		return llmProvider.HandleError(context.Background(), fmt.Errorf("构建聊天请求失败: %w", err), nil, nil)
	}

	// 3. 创建 HTTP 请求
	httpReq, err := http.NewRequest(reqParams.Method, reqParams.URL, bytes.NewReader(reqParams.Body))
	if err != nil {
		common.BaseLogger.WithError(err).Error("创建 HTTP 请求失败")
		return defaultError
	}
	// 设置请求头
	for key, value := range reqParams.Headers {
		httpReq.Header.Set(key, value)
	}

	// 4. 发送 HTTP 请求
	common.BaseLogger.Infof("发送请求到: %s", httpReq.URL)
	httpResp, err := http.DefaultClient.Do(httpReq)
	if err != nil {
		common.BaseLogger.WithError(err).Error("发送 HTTP 请求失败")
		// 尝试使用 Provider 的错误处理
		return llmProvider.HandleError(context.Background(), fmt.Errorf("发送 HTTP 请求失败: %w", err), nil, nil)
	}
	defer httpResp.Body.Close()

	// 5. 读取响应体
	body, err := io.ReadAll(httpResp.Body)
	if err != nil {
		common.BaseLogger.WithError(err).Error("读取响应体失败")
		// 尝试使用 Provider 的错误处理
		return llmProvider.HandleError(context.Background(), fmt.Errorf("读取响应体失败: %w", err), httpResp, nil)
	}

	// 6. 检查 HTTP 状态码 (通用检查)
	if httpResp.StatusCode != http.StatusOK {
		common.BaseLogger.WithField("status", httpResp.StatusCode).WithField("body", string(body)).Errorf("收到非 200 OK 状态码: %d", httpResp.StatusCode)
		// 尝试使用 Provider 的错误处理，传入状态码错误
		err = fmt.Errorf("HTTP 错误: %d %s", httpResp.StatusCode, http.StatusText(httpResp.StatusCode))
		return llmProvider.HandleError(context.Background(), err, httpResp, body)
	}

	// 7. 使用 Provider 解析响应体
	commonResponse, err := llmProvider.ParseChatResponse(context.Background(), body)
	if err != nil {
		common.BaseLogger.WithError(err).WithField("body", string(body)).Error("解析聊天响应失败")
		// 尝试使用 Provider 的错误处理
		return llmProvider.HandleError(context.Background(), fmt.Errorf("解析聊天响应失败: %w", err), httpResp, body)
	}

	// 8. 将通用响应设置到 chains
	chains.Res = commonResponse

	// 不再需要设置原始 res map
	// chains.Set("res", res)

	return nil
}

/*
// 原有的 chat 函数逻辑
func Chat_old(chains *ChatChains) error {
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
*/
