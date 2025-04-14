package _type

type FunctionDescription struct {
	Name        string  `json:"name"`
	Description *string `json:"description,omitempty"`
	Parameters  any     `json:"parameters"` // JSON Schema object
}

type ToolType string

const FunctionToolType ToolType = "function"

type Tool struct {
	// 目前只支持function
	Type     ToolType            `json:"type"`
	Function FunctionDescription `json:"function"`
}

type ToolChoiceFunctionType struct {
	Type     string `json:"type"`
	Function struct {
		Name string `json:"name"`
	} `json:"function"`
}

const ToolChoiceNone = "none"
const ToolChoiceAuto = "auto"

func ToolChoiceFunction(name string) *ToolChoiceFunctionType {
	return &ToolChoiceFunctionType{
		Type: "function",
		Function: struct {
			Name string `json:"name"`
		}{
			Name: name,
		},
	}

}
