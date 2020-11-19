package PushSDK

import "errors"

type MobileChannel interface {
	SendMessage(m *Message, pushId []string) (*Response, error)
}

// 消息体
type Message struct {
	Title        string
	Desc         string
	ApnsId       string
	ClickType    string
	ClickContent string
	Err          error
}

func NewMessage() *Message {
	return &Message{
		ClickType: "app",
	}
}

func (m *Message) SetTitle(str string) *Message {
	if str == "" {
		m.Err = errors.New("推送标题不能为空")
	}
	m.Title = str
	return m
}

func (m *Message) SetContent(str string) *Message {
	if str == "" {
		m.Err = errors.New("推送内容不能为空")
	}
	m.Desc = str
	return m
}

func (m *Message) SetApnsId(str string) *Message {
	m.ApnsId = str
	return m
}

func (m *Message) SetClickType(str string) *Message {
	m.ClickType = str
	return m
}

func (m *Message) SetClickContent(str string) *Message {
	m.ClickContent = str
	return m
}
