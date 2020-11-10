package PushSDK

import (
	"errors"
	"log"
)

//type Send struct {
//	MessageBody   MessageBody
//	Channel       string
//	PushId        []string
//	PlatformParam PlatformParam
//}

type Send struct {
	content map[string]interface{}
}

type MessageBody struct {
	Title string
	Desc  string
}

type PlatformParam struct {
	HWAppId                 string `json:"hw_appId"`
	HWClientSecret          string `json:"hw_clientSecret"`
	IOSKeyId                string `json:"iOS_keyId"`
	IOSTeamId               string `json:"iOS_teamId"`
	IOSBundleId             string `json:"iOS_bundleId"`
	IOSAuthTokenPath        string `json:"iOS_authTokenPath"`
	IOSAuthToken            string `json:"iOS_authToken"`
	MIAppSecret             string `json:"mi_appSecret"`
	MIRestrictedPackageName string `json:"mi_restrictedPackageName"`
	MZAppId                 string `json:"mz_appId"`
	MZAppSecret             string `json:"mz_appSecret"`
	OPPOAppKey              string `json:"oppo_appKey"`
	OPPOMasterSecret        string `json:"oppo_masterSecret"`
	VIAppID                 string `json:"vi_appId"`
	VIAppKey                string `json:"vi_appKey"`
	VIAppSecret             string `json:"vi_appSecret"`
	VIAuthToken             string `json:"vi_authToken"`
}

func NewSend() *Send {
	return &Send{
		map[string]interface{}{
			"messageBody": &MessageBody{},
			"channel":     "",
			"pushId":      []string{},
			"platform":    &PlatformParam{},
		},
	}
}

func (s *Send) message() *MessageBody {
	return s.content["messageBody"].(*MessageBody)
}

func (s *Send) SetTitle(str string) *Send {
	if len(str) > 40 {
		log.Println(errors.New("标题过长"))
	}
	s.message().Title = str
	return s
}

func (s *Send) SetContent(str string) *Send {
	s.message().Desc = str
	return s
}

func (s *Send) SetChannel(channel string) *Send {
	s.content["channel"] = channel
	return s
}

func (s *Send) SetPushId(pushId []string) *Send {
	s.content["pushId"] = pushId
	return s
}

func (s *Send) platform() *PlatformParam {
	return s.content["platform"].(*PlatformParam)
}

func (s *Send) SetHWAppId(str string) *Send {
	s.platform().HWAppId = str
	return s
}

func (s *Send) SetHWClientSecret(str string) *Send {
	s.platform().HWClientSecret = str
	return s
}

func (s *Send) SetIOSKeyId(str string) *Send {
	s.platform().IOSKeyId = str
	return s
}

func (s *Send) SetIOSTeamId(str string) *Send {
	s.platform().IOSTeamId = str
	return s
}

func (s *Send) SetIOSBundleId(str string) *Send {
	s.platform().IOSBundleId = str
	return s
}

func (s *Send) SetIOSAuthTokenPath(str string) *Send {
	s.platform().IOSAuthTokenPath = str
	return s
}

func (s *Send) SetIOSAuthToken(str string) *Send {
	s.platform().IOSAuthToken = str
	return s
}

func (s *Send) SetMIAppSecret(str string) *Send {
	s.platform().MIAppSecret = str
	return s
}

func (s *Send) SetMIRestrictedPackageName(str string) *Send {
	s.platform().MIRestrictedPackageName = str
	return s
}

func (s *Send) SetMZAppId(str string) *Send {
	s.platform().MZAppId = str
	return s
}

func (s *Send) SetMZAppSecret(str string) *Send {
	s.platform().MZAppSecret = str
	return s
}

func (s *Send) SetOPPOAppKey(str string) *Send {
	s.platform().OPPOAppKey = str
	return s
}

func (s *Send) SetOPPOMasterSecret(str string) *Send {
	s.platform().OPPOMasterSecret = str
	return s
}

func (s *Send) SetVIAppID(str string) *Send {
	s.platform().VIAppID = str
	return s
}

func (s *Send) SetVIAppKey(str string) *Send {
	s.platform().VIAppKey = str
	return s
}

func (s *Send) SetVIAppSecret(str string) *Send {
	s.platform().VIAppSecret = str
	return s
}

func (s *Send) SetVIAuthToken(str string) *Send {
	s.platform().VIAuthToken = str
	return s
}

//func InitSend(message *MessageBody, channel string, pushId []string, platformParam PlatformParam) *Send {
//	return &Send{
//		MessageBody:   *message,
//		Channel:       channel,
//		PushId:        pushId,
//		PlatformParam: platformParam,
//	}
//}

func (s *Send) SendMessage() (int, string) {
	code, errReason := 0, ""
	messageBody := s.content["messageBody"].(MessageBody)
	pushId := s.content["pushId"].([]string)
	platform := s.content["platform"].(PlatformParam)
	switch s.content["channel"].(string) {
	case "hw":
		code, errReason = hwMessagesSend(messageBody, pushId, platform.HWAppId, platform.HWClientSecret)
		break
	case "ios":
		code, errReason = iOSMessagesSend(messageBody, pushId, platform.IOSBundleId, platform.IOSAuthToken)
		break
	case "mi":
		code, errReason = miMessageSend(messageBody, pushId, platform.MIAppSecret, platform.MIRestrictedPackageName)
		break
	case "mz":
		code, errReason = mzMessageSend(messageBody, pushId, platform.MZAppId, platform.MZAppSecret)
		break
	case "oppo":
		code, errReason = oppoMessageSend(messageBody, pushId, platform.OPPOAppKey, platform.OPPOMasterSecret)
		break
	case "vivo":
		code, errReason = vSendMessage(messageBody, pushId, platform.VIAuthToken)
		break
	}
	return code, errReason
}
