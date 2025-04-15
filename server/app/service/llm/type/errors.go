package _type

import "errors"

// 429 error
var ErrResourceExhausted = errors.New("系统繁忙，请稍后再试")

// prompt被拦截
var ErrPromptBlocked = errors.New("系统检测到您的输入包含敏感内容，请重新输入")
