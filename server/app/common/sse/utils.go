package sse

import (
	"strings"
)

// 解析SSE事件的内容
func ParseSSEEvent(event string) (data, eventType string) {
	lines := strings.Split(event, "\n")
	for _, line := range lines {
		trimmedLine := strings.TrimSpace(line)
		if strings.HasPrefix(trimmedLine, "data:") {
			data = strings.TrimSpace(trimmedLine[5:])
		} else if strings.HasPrefix(trimmedLine, "event:") {
			eventType = strings.TrimSpace(trimmedLine[6:])
		}
	}
	return data, eventType
}

// 自定义的Scanner分隔函数，用于按照SSE事件的格式进行分隔
func MakeSplit(provider string) func(data []byte, atEOF bool) (advance int, token []byte, err error) {
	return func(data []byte, atEOF bool) (advance int, token []byte, err error) {
		if atEOF && len(data) == 0 {
			return 0, nil, nil
		}
		if provider == "cohere" {
			// \n就返回
			if i := strings.Index(string(data), "\n"); i >= 0 {
				return i + 1, data[0:i], nil
			}
		}

		if i := strings.Index(string(data), "\n\n"); i >= 0 {
			return i + 2, data[0:i], nil
		}
		// 有可能是\r\n\r\n
		if i := strings.Index(string(data), "\r\n\r\n"); i >= 0 {
			return i + 4, data[0:i], nil
		}

		if atEOF {
			return len(data), data, nil
		}
		// azure中，done后面只有一个换行。这里兼容这种情况
		if provider == "azure" {
			if i := strings.Index(string(data), "[DONE]\n"); i >= 0 {
				return i + 7, data[0 : i+6], nil
			}
		}
		if i := strings.Index(string(data), "[DONE]\n\n\n"); i >= 0 {
			return i + 9, data[0 : i+8], nil
		}

		return 0, nil, nil
	}
}
