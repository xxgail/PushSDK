package PushSDK

import "time"

type Response struct {
	Code      int
	Reason    string
	ApnsId    string
	TimeStamp time.Time
}
