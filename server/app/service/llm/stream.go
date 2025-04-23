package llm

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"

	"oukek/crab-ai/app/common"
	"oukek/crab-ai/app/common/sse"
	"oukek/crab-ai/app/model"
	"oukek/crab-ai/app/service/chains"
	_type "oukek/crab-ai/app/service/llm/type"

	"github.com/gin-gonic/gin"
	"github.com/mitchellh/mapstructure"
)

func MakeChatWithStream(c *gin.Context) chains.HandlerFunc[*ChatChains] {
	if c == nil {
		return func(chains *ChatChains) error {
			r := make(chan _type.MessageChanItem, 100)
			go common.WithSafeFn(func() {
				ChatWithStream(chains, r)
			})

			buf := ""
			content := ""
			common.BaseLogger.Debugf("开始接收消息")
			res := _type.Response{}
		done:
			for {
				select {
				case m, ok := <-r:
					if !ok {
						common.BaseLogger.Debugf("收到所有的内容为：" + content)
						break done
					}
					switch m.Event {
					case "err":
						return errors.New(m.Data.(string))
					case "finish":
						break done
					case "init":
						initData := m.Data.(_type.StreamItem)
						res.ID = initData.Id
						res.Object = initData.Object
						res.Model = initData.Model
						res.Created = int64(initData.Created)
					case "data":
						common.BaseLogger.Debugf("收到消息:%v", m)
						buf += m.Data.(string)
						// 节省带宽，提升速度，同时又提升首屏速度
						if len(buf) > 8 || len(content) < 20 {
							content += buf
							buf = ""
						}

					}
				}
			}
			if len(buf) > 0 {
				content += buf
			}
			res.Choices = []_type.NonStreamingChoice{
				{
					Index:        0,
					FinishReason: "stop",
					Message: &_type.ResMessage{
						Role:    "assistant",
						Content: content,
					},
				},
			}
			chains.Res = &res

			common.BaseLogger.Debugf("结束")
			return nil
		}
	}
	c.Header("Content-Type", "text/event-stream")
	c.Header("Cache-Control", "no-cache")
	c.Header("Connection", "keep-alive")
	c.Header("X-Accel-Buffering", "no") // Nginx 配置

	return func(chains *ChatChains) error {
		r := make(chan _type.MessageChanItem, 100)
		go common.WithSafeFn(func() {
			ChatWithStream(chains, r)
		})

		buf := ""
		content := ""
		common.BaseLogger.Debugf("开始接收消息")
		isCancel := false
		res := _type.Response{}
	done:
		for {
			select {
			case <-c.Done():
				common.BaseLogger.Debug("context被触发了，暂停请求")
				isCancel = true
			case <-c.Writer.CloseNotify():
				common.BaseLogger.Debug("context被触发了，客户端关闭连接")
				isCancel = true
			case m, ok := <-r:
				if !ok {
					common.BaseLogger.Debugf("收到所有的内容为：" + content)
					break done
				}
				switch m.Event {
				case "err":
					if !isCancel {
						_, err := fmt.Fprintf(c.Writer, "event: %s\ndata: %s\n\n", "err", m.Data)
						if err != nil {
							common.BaseLogger.WithError(err).Error("gpt stream刷到客户端失败")
							isCancel = true
						}
						common.BaseLogger.Debugf("失败:%v", m.Data)
						return err

					}
				case "finish":
					break done
				case "init":
					initData := m.Data.(_type.StreamItem)
					res.ID = initData.Id
					res.Object = initData.Object
					res.Model = initData.Model
					res.Created = int64(initData.Created)
				case "data":
					common.BaseLogger.Debugf("收到消息:%v", m)
					buf += m.Data.(string)
					// 节省带宽，提升速度，同时又提升首屏速度
					if len(buf) > 8 || len(content) < 20 {
						buf = strings.Replace(buf, "\n", "\\n", -1)
						if !isCancel {
							n, err := fmt.Fprintf(c.Writer, "data: %s\n\n", buf)
							if err != nil {
								common.BaseLogger.WithError(err).Debug("gpt stream刷到客户端失败")
								isCancel = true
							} else {
								c.Writer.Flush()
							}
							common.BaseLogger.Debugf("输出，输出长度：%d", n)
						}
						content += buf
						buf = ""

					}

				}
			}
		}
		if len(buf) > 0 {
			buf = strings.Replace(buf, "\n", "\\n", -1)
			if !isCancel {
				n, err := fmt.Fprintf(c.Writer, "data: %s\n\n", buf)
				if err != nil {
					common.BaseLogger.WithError(err).Error("gpt stream刷到客户端失败")
				} else {
					c.Writer.Flush()
					common.BaseLogger.Debugf("输出，输出长度：%d", n)
				}
			}
			content += buf
		}
		res.Choices = []_type.NonStreamingChoice{
			{
				Index:        0,
				FinishReason: "stop",
				Message: &_type.ResMessage{
					Role:    "assistant",
					Content: content,
				},
			},
		}
		chains.Res = &res

		common.BaseLogger.Debugf("结束")
		return nil
	}
}

// 发送请求并处理响应
func ChatWithStream(chains *ChatChains, messageChan chan<- _type.MessageChanItem) {
	// 自动关闭通道
	defer close(messageChan)

	// 定义错误消息
	const (
		defaultErrorMsg     = "服务异常，请联系客服"
		sensitiveContentMsg = "指令存在敏感内容，已被拦截，请修改重试"
		tooManyRequestsMsg  = "请求过于频繁，请稍后再试"
		generateFailedMsg   = "生成失败，请重试"
		contentFilteredMsg  = "内容存在敏感内容，已被过滤，请不要生成违法内容。如有疑问请咨询客服"
	)

	// 发送错误消息的辅助函数
	sendError := func(msg string) {
		messageChan <- _type.MessageChanItem{Event: "err", Data: msg}
	}

	req, ok := chains.Get("req")
	if !ok {
		sendError(defaultErrorMsg)
		return
	}

	// 发送请求并处理响应
	resp, err := http.DefaultClient.Do(req.(*http.Request))
	if err != nil {
		common.BaseLogger.WithError(err).Error("GPT：发送请求失败")
		sendError(defaultErrorMsg)
		return
	}
	defer resp.Body.Close()

	// 检查响应状态
	if err := checkResponseStatus(resp, chains, messageChan); err != nil {
		return
	}

	// 处理流式响应
	handleStreamResponse(resp, chains, messageChan)
	common.BaseLogger.Debug("结束请求")
}

// 检查响应状态码
func checkResponseStatus(resp *http.Response, chains *ChatChains, messageChan chan<- _type.MessageChanItem) error {
	if resp.StatusCode == http.StatusOK {
		return nil
	}

	body, _ := io.ReadAll(resp.Body)
	common.BaseLogger.WithField("data", chains.Req).
		WithField("status code", resp.StatusCode).
		WithField("body", string(body)).
		Error("GPT：发送请求失败")

	switch resp.StatusCode {
	// 400
	case http.StatusBadRequest:
		var errorResponse struct {
			Error struct {
				Message string `json:"message"`
				Code    string `json:"code"`
			} `json:"error"`
		}
		if err := json.Unmarshal(body, &errorResponse); err == nil && errorResponse.Error.Code == "content_filter" {
			messageChan <- _type.MessageChanItem{Event: "err", Data: "指令存在敏感内容，已被拦截，请修改重试"}
			return fmt.Errorf("content filtered")
		}
	// 429
	case http.StatusTooManyRequests:
		messageChan <- _type.MessageChanItem{Event: "err", Data: "请求过于频繁，请稍后再试"}
		return fmt.Errorf("too many requests")
	}

	messageChan <- _type.MessageChanItem{Event: "err", Data: "服务异常，请联系客服"}
	return fmt.Errorf("unexpected status code: %d", resp.StatusCode)
}

// 处理流式响应
func handleStreamResponse(resp *http.Response, chains *ChatChains, messageChan chan<- _type.MessageChanItem) {
	common.BaseLogger.Debug("开始读取数据")
	scanner := bufio.NewScanner(resp.Body)
	scanner.Split(sse.MakeSplit(string(chains.Provider.Name)))
	isInit := false

	for {
		if !scanner.Scan() {
			if err := scanner.Err(); err != nil {
				common.BaseLogger.WithError(err).Error("chat gpt stream扫描失败")
			}
			common.BaseLogger.Debugf("到达尾部:%v", scanner.Text())
			return
		}

		event := scanner.Text()
		common.BaseLogger.Debug("收到事件：", event)

		// OPENROUTER服务商特殊的msg，它们为了防止等待过长，会返回一个处理中的msg
		if event == ": OPENROUTER PROCESSING" {
			continue
		}

		// 解析SSE事件
		data, _ := sse.ParseSSEEvent(event)
		common.BaseLogger.Debug("事件数据：", data)

		if data == "[DONE]" {
			return
		}

		rec, err := parseStreamResponse(data, chains, messageChan)
		if err != nil {
			messageChan <- _type.MessageChanItem{Event: "err", Data: "生成失败，请重试"}
			return
		}

		if err := rec.IsValid(); err != nil {
			messageChan <- _type.MessageChanItem{Event: "err", Data: err.Error()}
			continue
		}

		// azure现在是在第一条返回对prompt的过滤结果
		if len(rec.Id) == 0 {
			common.BaseLogger.WithField("data", rec).Info("对话过滤")
			continue
		}

		// 如果内容为空，则跳过
		if len(rec.Choices[0].Delta.Content) == 0 {
			continue
		}

		if !isInit {
			isInit = true
			messageChan <- _type.MessageChanItem{Event: "init", Data: *rec}
		}
		messageChan <- _type.MessageChanItem{Event: "data", Data: rec.Choices[0].Delta.Content}
	}
}

// 解析流式响应
func parseStreamResponse(data string, chains *ChatChains, messageChan chan<- _type.MessageChanItem) (*_type.StreamItem, error) {
	var rec *_type.StreamItem

	var metadata mapstructure.Metadata
	// 根据不同的LLM类型解析响应
	switch chains.Provider.Type {
	case model.ProviderTypeGoogleAIStudio:
		var r _type.GeminiStreamItem
		err := mapstructure.DecodeMetadata(data, &r, &metadata)
		if err != nil {
			return nil, err
		}
		rec = r.ToStreamItem(string(chains.Req.Model))

	default:
		err := mapstructure.DecodeMetadata(data, &rec, &metadata)
		if err != nil {
			return nil, err
		}

	}

	return rec, nil
}
