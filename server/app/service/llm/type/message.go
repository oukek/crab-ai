package _type

type ContentType string

const (
	TextContentType  ContentType = "text"
	ImageContentType ContentType = "image_url"
)

type ContentPart struct {
	Type     ContentType `json:"type"`
	Text     string      `json:"text,omitempty"`
	ImageURL *struct {
		URL    string `json:"url"`
		Detail string `json:"detail,omitempty"`
	} `json:"image_url,omitempty"`
}

type Content interface {
	string | []ContentPart
}

// Removed StreamItemChoiceDelta, StreamItemChoice, StreamItem, and IsValid method
// These types are moved to response.go for better organization.
