package llm

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"time"

	"github.com/mitchellh/mapstructure"
	"oukek/crab-ai/app/common"
	"oukek/crab-ai/app/service/llm/type"
)

func Chat(data *_type.Request) (*_type.Response, error) {
	reqData := data
	defaultError := errors.New("服务器异常，请联系客服")
	// 创建一个json对象
	// 将json对象转换为字节切片
	var jsonData []byte
	var err error
	if data.IsGoogleGemini() {
		reqData2, err := _type.NewGeminiReqFromRequest(reqData)
		if err != nil {
			common.BaseLogger.WithError(err).Error("GPT：body遍码失败")
			return nil, defaultError
		}
		jsonData, err = json.Marshal(reqData2)
	} else {
		jsonData, err = json.Marshal(reqData)
	}
	if err != nil {
		common.BaseLogger.WithError(err).Error("GPT：body遍码失败")
		return nil, defaultError
	}

	// 创建一个http请求，设置请求方法为POST，请求地址为a.com/api，请求体为jsonData
	req, err := getReq(jsonData, data.Provider, false, data.Model)

	var body []byte
	// 发送请求并获取响应
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		common.BaseLogger.WithError(err).Error("GPT：发送请求失败")
		return nil, defaultError
	}
	// 延迟关闭响应体
	defer resp.Body.Close()

	// 读取响应体的内容
	body, err = io.ReadAll(resp.Body)
	if err != nil {
		common.BaseLogger.WithError(err).Error("GPT：读取body失败")
		return nil, defaultError
	}
	// 判断状态码
	if resp.StatusCode != http.StatusOK {
		// 429 代表请求过于频繁
		if resp.StatusCode == http.StatusTooManyRequests {
			// 如果是CF服务商，冻结当前key
			if data.Provider == _type.ProviderCF {
				if key := req.URL.Query().Get("key"); key != "" {
					if err := FreezeCFKey(key, time.Minute); err != nil {
						common.BaseLogger.WithError(err).Error("GPT：冻结 CF key 失败")
					}
				}
			}
			common.BaseLogger.WithField("body", string(body)).WithField("status", resp.StatusCode).Error("GPT：请求过于频繁")
			return nil, _type.ErrResourceExhausted
		}
		common.BaseLogger.WithField("body", string(body)).WithField("status", resp.StatusCode).Error("GPT：请求失败")
		return nil, defaultError
	}

	res := _type.M{}
	err = json.Unmarshal(body, &res)
	if err != nil {
		common.BaseLogger.WithError(err).WithField("res", string(body)).Error("GPT：转换成结构体失败")
		return nil, defaultError
	}

	// 获取object，如果存在的话，则判断是不是error
	if object, ok := res["object"]; ok {
		if object == "error" {
			// 获取message，如果存在的话，则打日志的时候加上，否则为空
			errMsg := ""
			if message, ok := res["message"]; ok {
				errMsg = message.(string)
			}
			common.BaseLogger.WithField("req", string(jsonData)).WithField("res", string(body)).Error("GPT：请求失败:" + errMsg)
			return nil, defaultError
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
		common.BaseLogger.WithField("req", string(jsonData)).WithField("res", string(body)).Error("GPT：请求失败:" + errMsg)
		return nil, defaultError
	}
	// M map[string]interface{}
	// 从res中，解析到_type.Response
	// 如果解析失败，则返回默认错误
	r2 := _type.Response{}
	var metadata mapstructure.Metadata
	if data.IsGoogleGemini() {
		var r3 _type.GeminiRes
		err = mapstructure.DecodeMetadata(res, &r3, &metadata)
		// 打印r3的值
		common.BaseLogger.WithField("r3", r3).Info("GPT：解析结果")
		r4, err := r3.ToGPTResponse()
		if err != nil {
			return nil, err
		}
		r2 = *r4

	} else {
		err = mapstructure.DecodeMetadata(res, &r2, &metadata)
	}

	if err != nil {
		common.BaseLogger.WithError(err).WithField("res", string(body)).Error("GPT：解析结果失败")
		return nil, defaultError
	}
	if len(metadata.Unused) > 0 {
		common.BaseLogger.WithField("unused", metadata.Unused).WithField("res", string(body)).Info("GPT：存在未解析字段")
	}

	return &r2, nil
}
