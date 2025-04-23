package llm

import (
	"errors"

	"oukek/crab-ai/app/model"
	_type "oukek/crab-ai/app/service/llm/type"

	"github.com/mitchellh/mapstructure"
)

func MakeResponse(chains *ChatChains) error {
	chains.Next()

	_res, ok := chains.Get("res")
	if !ok {
		return errors.New("GPT：响应体不存在")
	}
	res := _res.(map[string]any)
	var metadata mapstructure.Metadata
	r2 := &_type.Response{}
	var err error
	switch chains.Provider.Type {
	case model.ProviderTypeGoogleAIStudio:
		var r _type.GeminiRes
		err := mapstructure.DecodeMetadata(res, &r, &metadata)
		if err != nil {
			return err
		}
		r2, err = r.ToResponse()

	default:
		err = mapstructure.DecodeMetadata(res, r2, &metadata)

	}
	if err != nil {
		return err
	}
	chains.Res = r2
	return nil
}
