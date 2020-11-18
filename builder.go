package PushSDK

type Message struct {
	Fields interface{} //含有本条消息所有属性的数组
}

type MobileChannel interface {
	SendMessage(body *MessageBody, pushId []string) (*Response, error)
}
