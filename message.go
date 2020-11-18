package PushSDK

import "errors"

type MessageBody struct {
	Title        string
	Desc         string
	ApnsId       string
	ClickType    string
	ClickContent string
	Err          error
}

func NewMessage() *MessageBody {
	return &MessageBody{
		ClickType: "app",
	}
}

func (m *MessageBody) SetTitle(str string) *MessageBody {
	if str == "" {
		m.Err = errors.New("推送标题不能为空")
	}
	m.Title = str
	return m
}

func (m *MessageBody) SetContent(str string) *MessageBody {
	if str == "" {
		m.Err = errors.New("推送内容不能为空")
	}
	m.Desc = str
	return m
}

func (m *MessageBody) SetApnsId(str string) *MessageBody {
	m.ApnsId = str
	return m
}

func (m *MessageBody) SetClickType(str string) *MessageBody {
	m.ClickType = str
	return m
}

func (m *MessageBody) SetClickContent(str string) *MessageBody {
	m.ClickContent = str
	return m
}
