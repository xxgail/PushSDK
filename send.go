package PushSDK

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"sync"
)

//type Send struct {
//	MessageBody   MessageBody
//	Channel       string
//	PushId        []string
//	PlatformParam PlatformParam
//}

type Send struct {
	Content map[string]interface{}
}

type MessageBody struct {
	Title        string
	Desc         string
	ApnsId       string
	ClickType    string
	ClickContent string
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
			"iosParam":    &IOSParam{},
		},
	}
}

// 设置渠道
func (s *Send) SetChannel(channel string) *Send {
	s.Content["channel"] = channel
	return s
}

// 设置推送用户
func (s *Send) SetPushId(pushId []string) *Send {
	s.Content["pushId"] = pushId
	return s
}

func (s *Send) message() *MessageBody {
	return s.Content["messageBody"].(*MessageBody)
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

func (s *Send) SetApnsId(str string) *Send {
	s.message().ApnsId = str
	return s
}

func (s *Send) SetClickType(str string) *Send {
	s.message().ClickType = str
	return s
}

func (s *Send) SetClickContent(str string) *Send {
	s.message().ClickContent = str
	return s
}

func (s *Send) platform() *PlatformParam {
	return s.Content["platform"].(*PlatformParam)
}

func (s *Send) SetHWAppId(str string) *Send {
	s.platform().HWAppId = str
	return s
}

func (s *Send) SetHWClientSecret(str string) *Send {
	s.platform().HWClientSecret = str
	return s
}

type IOSParam struct {
	sync.Mutex
	KeyId         string `json:"key_id"`
	TeamId        string `json:"team_id"`
	BundleId      string `json:"bundle_id"`
	AuthTokenPath string `json:"auth_token_path"`
	Bearer        string `json:"bearer"`
	IssuedAt      int64  `json:"issued_at"`
}

func (s *Send) iosParam() *IOSParam {
	return s.Content["iosParam"].(*IOSParam)
}

func (s *Send) SetIOSParam(str string) (*Send, error) {
	var iosParam IOSParam
	_ = json.Unmarshal([]byte(str), &iosParam)
	s.iosParam().KeyId = iosParam.KeyId
	s.iosParam().TeamId = iosParam.TeamId
	s.iosParam().BundleId = iosParam.BundleId
	s.iosParam().AuthTokenPath = iosParam.AuthTokenPath
	if iosParam.Bearer == "" {
		s.iosParam().Bearer = iosParam.generateIfExpired()
	} else {
		s.iosParam().Bearer = iosParam.Bearer
	}
	return s, nil
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

func (s *Send) SendMessage() (*Response, error) {
	var messageBody MessageBody
	mPoint := s.Content["messageBody"].(*MessageBody)
	mJson, _ := json.Marshal(mPoint)
	json.Unmarshal(mJson, &messageBody)
	fmt.Println("messageBody", messageBody)
	if messageBody.ClickType == "" {
		messageBody.ClickType = "app"
	} else if messageBody.ClickType != "app" && messageBody.ClickContent == "" {
		log.Println("点击内容不能为空")
	}

	pushId := s.Content["pushId"].([]string)

	var platform PlatformParam
	pPoint := s.Content["platform"].(*PlatformParam)
	pJson, _ := json.Marshal(pPoint)
	_ = json.Unmarshal(pJson, &platform)
	fmt.Println("platform", platform)
	switch s.Content["channel"].(string) {
	case "hw":
		return hwMessagesSend(messageBody, pushId, platform.HWAppId, platform.HWClientSecret)
	case "ios":
		iosBundleId := s.Content["iosParam"].(*IOSParam).BundleId
		iosToken := s.Content["iosParam"].(*IOSParam).Bearer
		return iOSMessagesSend(messageBody, pushId, iosBundleId, iosToken)
	case "mi":
		return miMessageSend(messageBody, pushId, platform.MIAppSecret, platform.MIRestrictedPackageName)
	case "mz":
		return mzMessageSend(messageBody, pushId, platform.MZAppId, platform.MZAppSecret)
	case "oppo":
		return oppoMessageSend(messageBody, pushId, platform.OPPOAppKey, platform.OPPOMasterSecret)
	case "vivo":
		return vSendMessage(messageBody, pushId, platform.VIAuthToken)
	default:
		return &Response{
			Code:   SendError,
			Reason: "No channel",
		}, nil
	}
}
