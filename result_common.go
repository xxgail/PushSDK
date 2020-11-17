package PushSDK

import "time"

type Response struct {
	Code      int
	Reason    string
	ApnsId    string
	TimeStamp time.Time
}

func (r *Response) init() *Response {
	return &Response{
		TimeStamp: time.Now(),
	}
}
