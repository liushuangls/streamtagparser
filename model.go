package streamtagparser

type TagStreamType string

const (
	TagStreamTypeText    TagStreamType = "text"
	TagStreamTypeStart   TagStreamType = "start"
	TagStreamTypeContent TagStreamType = "content"
	TagStreamTypeEnd     TagStreamType = "end"
)

type TagAttr struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

type TagStreamData struct {
	Type TagStreamType `json:"type"`

	Text string `json:"text,omitempty"` // text is the content outside the tag

	TagName string    `json:"tag_name,omitempty"`
	Attrs   []TagAttr `json:"attrs,omitempty"`
	Content string    `json:"content,omitempty"` // content is the content of the tag
}

func NewTextTagStreamData(text string) *TagStreamData {
	return &TagStreamData{
		Type: TagStreamTypeText,
		Text: text,
	}
}

func NewStartTagStreamData(tagName string, attrs []TagAttr) *TagStreamData {
	return &TagStreamData{
		Type:    TagStreamTypeStart,
		TagName: tagName,
		Attrs:   attrs,
	}
}

func NewContentTagStreamData(tagName, content string) *TagStreamData {
	return &TagStreamData{
		Type:    TagStreamTypeContent,
		TagName: tagName,
		Content: content,
	}
}

func NewEndTagStreamData(tagName string, attrs []TagAttr, content string) *TagStreamData {
	return &TagStreamData{
		Type:    TagStreamTypeEnd,
		TagName: tagName,
		Attrs:   attrs,
		Content: content,
	}
}
