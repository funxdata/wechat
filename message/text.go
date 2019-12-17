package message

//Text 文本消息
type Text struct {
	CommonToken
	Content CDATA `xml:"Content"`
}

//NewText 初始化文本消息
func NewText(content string) *Text {
	text := new(Text)
	text.Content = CDATA(content)
	return text
}

// NewTextReply 初始化文本消息消息回复
func NewTextReply(content string) *Reply {
	return &Reply{
		MsgType: MsgTypeText,
		MsgData: NewText(content),
	}
}
