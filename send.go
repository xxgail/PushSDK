package PushSDK

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"reflect"
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
	Err     error
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
		Content: map[string]interface{}{
			"messageBody": &MessageBody{},
			"channel":     "",
			"pushId":      []string{},
			"platform":    "",
			//"platform":    &PlatformParam{},
			"iosParam":  &IOSParam{},
			"hwParam":   &HWParam{},
			"miParam":   &MIParam{},
			"mzParam":   &MZParam{},
			"oppoParam": &OPPOParam{},
			"vParam":    &VParam{},
		},
		Err: nil,
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

func (s *Send) SetPlatForm(plat string) *Send {
	s.Content["platform"] = plat
	return s
}

func (s *Send) message() *MessageBody {
	return s.Content["messageBody"].(*MessageBody)
}

func (s *Send) SetTitle(str string) *Send {
	if str == "" {
		s.Err = errors.New("推送标题不能为空")
	}
	s.message().Title = str
	return s
}

func (s *Send) SetContent(str string) *Send {
	if str == "" {
		s.Err = errors.New("推送内容不能为空")
	}
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

//func (s *Send) platform() *PlatformParam {
//	return s.Content["platform"].(*PlatformParam)
//}

type HWParam struct {
	sync.Mutex
	AppId        string `json:"app_id"`
	ClientSecret string `json:"client_secret"`
}

func (s *Send) hwParam() *HWParam {
	return s.Content["hwParam"].(*HWParam)
}

func (s *Send) SetHWParam(str string) *Send {
	s.Err = isEmpty(str)
	var param HWParam
	_ = json.Unmarshal([]byte(str), &param)
	s.hwParam().AppId = param.AppId
	s.hwParam().ClientSecret = param.ClientSecret
	t := reflect.TypeOf(param)
	v := reflect.ValueOf(param)
	for k := 0; k < t.NumField(); k++ {
		if v.Field(k).Interface() == nil {
			s.Err = errors.New(t.Field(k).Name + "不能为空")
		}
	}
	return s
}

func (s *Send) SetHWAppId(str string) *Send {
	s.hwParam().AppId = str
	return s
}

func (s *Send) SetHWClientSecret(str string) *Send {
	s.hwParam().ClientSecret = str
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

func (s *Send) SetIOSParam(str string) *Send {
	s.Err = isEmpty(str)
	var param IOSParam
	_ = json.Unmarshal([]byte(str), &param)
	if param.KeyId == "" {
		s.Err = errors.New("KeyId" + "can not be empty")
		return s
	}
	if param.TeamId == "" {
		s.Err = errors.New("TeamId" + "can not be empty")
		return s
	}
	if param.BundleId == "" {
		s.Err = errors.New("BundleId" + "can not be empty")
		return s
	}
	if param.AuthTokenPath == "" {
		s.Err = errors.New("AuthTokenPath" + "can not be empty")
		return s
	}
	s.iosParam().KeyId = param.KeyId
	s.iosParam().TeamId = param.TeamId
	s.iosParam().BundleId = param.BundleId
	s.iosParam().AuthTokenPath = param.AuthTokenPath
	if param.Bearer == "" {
		s.iosParam().Bearer = param.generateIfExpired()
	} else {
		s.iosParam().Bearer = param.Bearer
	}
	return s
}

func (s *Send) SetIOSKeyId(str string) *Send {
	s.iosParam().KeyId = str
	return s
}

func (s *Send) SetIOSTeamId(str string) *Send {
	s.iosParam().TeamId = str
	return s
}

func (s *Send) SetIOSBundleId(str string) *Send {
	s.iosParam().BundleId = str
	return s
}

func (s *Send) SetIOSAuthTokenPath(str string) *Send {
	s.iosParam().AuthTokenPath = str
	return s
}

func (s *Send) SetIOSAuthToken(str string) *Send {
	s.iosParam().Bearer = str
	return s
}

type MIParam struct {
	AppSecret             string `json:"app_secret"`
	RestrictedPackageName string `json:"restricted_package_name"`
}

func (s *Send) miParam() *MIParam {
	return s.Content["miParam"].(*MIParam)
}

func (s *Send) SetMIParam(str string) *Send {
	s.Err = isEmpty(str)
	var param MIParam
	_ = json.Unmarshal([]byte(str), &param)
	s.miParam().AppSecret = param.AppSecret
	s.miParam().RestrictedPackageName = param.RestrictedPackageName
	t := reflect.TypeOf(param)
	v := reflect.ValueOf(param)
	for k := 0; k < t.NumField(); k++ {
		if v.Field(k).Interface() == nil {
			s.Err = errors.New(t.Field(k).Name + "不能为空")
		}
	}
	return s
}

func (s *Send) SetMIAppSecret(str string) *Send {
	s.miParam().AppSecret = str
	return s
}

func (s *Send) SetMIRestrictedPackageName(str string) *Send {
	s.miParam().RestrictedPackageName = str
	return s
}

type MZParam struct {
	AppSecret string `json:"app_secret"`
	AppId     string `json:"app_id"`
}

func (s *Send) mzParam() *MZParam {
	return s.Content["mzParam"].(*MZParam)
}

func (s *Send) SetMZParam(str string) *Send {
	s.Err = isEmpty(str)
	var param MZParam
	_ = json.Unmarshal([]byte(str), &param)
	s.mzParam().AppSecret = param.AppSecret
	s.mzParam().AppId = param.AppId
	t := reflect.TypeOf(param)
	v := reflect.ValueOf(param)
	for k := 0; k < t.NumField(); k++ {
		if v.Field(k).Interface() == nil {
			s.Err = errors.New(t.Field(k).Name + "不能为空")
		}
	}
	return s
}

func (s *Send) SetMZAppId(str string) *Send {
	s.mzParam().AppId = str
	return s
}

func (s *Send) SetMZAppSecret(str string) *Send {
	s.mzParam().AppSecret = str
	return s
}

type OPPOParam struct {
	AppKey       string `json:"app_key"`
	MasterSecret string `json:"master_secret"`
}

func (s *Send) oppoParam() *OPPOParam {
	return s.Content["oppoParam"].(*OPPOParam)
}

func (s *Send) SetOPPOParam(str string) *Send {
	s.Err = isEmpty(str)
	var param OPPOParam
	_ = json.Unmarshal([]byte(str), &param)
	s.oppoParam().AppKey = param.AppKey
	s.oppoParam().MasterSecret = param.MasterSecret
	t := reflect.TypeOf(param)
	v := reflect.ValueOf(param)
	for k := 0; k < t.NumField(); k++ {
		if v.Field(k).Interface() == nil {
			s.Err = errors.New(t.Field(k).Name + "不能为空")
		}
	}
	return s
}

func (s *Send) SetOPPOAppKey(str string) *Send {
	s.oppoParam().AppKey = str
	return s
}

func (s *Send) SetOPPOMasterSecret(str string) *Send {
	s.oppoParam().MasterSecret = str
	return s
}

type VParam struct {
	sync.Mutex
	AppID     string `json:"app_id"`
	AppKey    string `json:"app_key"`
	AppSecret string `json:"app_secret"`
	AuthToken string `json:"auth_token"`
	IssuedAt  int64  `json:"issued_at"`
}

func (s *Send) vParam() *VParam {
	return s.Content["vParam"].(*VParam)
}

func (s *Send) SetVIVOParam(str string) *Send {
	s.Err = isEmpty(str)
	var param VParam
	_ = json.Unmarshal([]byte(str), &param)
	if param.AppID == "" {
		s.Err = errors.New("AppId" + "can not be empty")
		return s
	}
	if param.AppKey == "" {
		s.Err = errors.New("AppKey" + "can not be empty")
		return s
	}
	if param.AppSecret == "" {
		s.Err = errors.New("AppSecret" + "can not be empty")
		return s
	}
	s.vParam().AppID = param.AppID
	s.vParam().AppKey = param.AppKey
	s.vParam().AppSecret = param.AppSecret
	if param.AuthToken == "" {
		s.vParam().AuthToken = param.generateIfExpired()
	} else {
		s.vParam().AuthToken = param.AuthToken
	}
	return s
}

func (s *Send) SetVIAppID(str string) *Send {
	s.vParam().AppID = str
	return s
}

func (s *Send) SetVIAppKey(str string) *Send {
	s.vParam().AppKey = str
	return s
}

func (s *Send) SetVIAppSecret(str string) *Send {
	s.vParam().AppSecret = str
	return s
}

func (s *Send) SetVIAuthToken(str string) *Send {
	s.vParam().AuthToken = str
	return s
}

func (s *Send) SendMessage() (*Response, error) {
	var messageBody MessageBody
	mPoint := s.message()
	mJson, _ := json.Marshal(mPoint)
	json.Unmarshal(mJson, &messageBody)
	fmt.Println("messageBody", messageBody)
	if messageBody.ClickType == "" {
		messageBody.ClickType = "app"
	} else if messageBody.ClickType != "app" && messageBody.ClickContent == "" {
		log.Println("点击内容不能为空")
	}

	pushId := s.Content["pushId"].([]string)
	plat := s.Content["platform"].(string)
	switch s.Content["channel"].(string) {
	case "hw":
		if s.hwParam() == nil {
			s.SetHWParam(plat)
		}
		fmt.Println(s.hwParam())
		return hwMessagesSend(messageBody, pushId, s.hwParam())
	case "ios":
		if s.iosParam() == nil {
			s.SetIOSParam(plat)
		}
		return iOSMessagesSend(messageBody, pushId, s.iosParam())
	case "mi":
		if s.miParam() == nil {
			s.SetMIParam(plat)
		}
		return miMessageSend(messageBody, pushId, s.miParam())
	case "mz":
		if s.mzParam() == nil {
			s.SetMZParam(plat)
		}
		return mzMessageSend(messageBody, pushId, s.mzParam())
	case "oppo":
		if s.oppoParam() == nil {
			s.SetOPPOParam(plat)
		}
		return oppoMessageSend(messageBody, pushId, s.oppoParam())
	case "vivo":
		if s.vParam() == nil {
			s.SetVIVOParam(plat)
		}
		return vSendMessage(messageBody, pushId, s.vParam())
	default:
		return &Response{
			Code:   SendError,
			Reason: "No channel",
		}, nil
	}
}
