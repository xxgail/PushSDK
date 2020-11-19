package PushSDK

import (
	"encoding/json"
	"errors"
	"strings"
)

type Send struct {
	Channel  string
	PushId   []string
	PlatForm string
	Err      error
}

func NewSend() *Send {
	return &Send{
		Err: nil,
	}
}

var (
	AvailableChannel = []string{"ios", "mi", "hw", "oppo", "vivo", "mz"}
)

func checkChannel(str string) error {
	for _, v := range AvailableChannel {
		if v == str {
			return nil
		}
	}
	return errors.New("channel 参数错误！")
}

// 设置渠道
func (s *Send) SetChannel(channel string) *Send {
	channel = strings.ToLower(channel)
	s.Err = checkChannel(channel)
	s.Channel = channel
	return s
}

// 设置推送用户
func (s *Send) SetPushId(pushId []string) *Send {
	if len(pushId) == 0 {
		s.Err = errors.New("发送用户不能为空")
	}
	s.PushId = pushId
	return s
}

func (s *Send) SetPlatForm(plat string) *Send {
	s.Err = isEmpty(plat)
	s.PlatForm = plat
	return s
}

func (s *Send) SendMessage(message *Message) (*Response, error) {
	var (
		err error
	)
	if message.Title == "" {
		return &Response{}, errors.New("title can not be empty")
	}
	if message.Desc == "" {
		return &Response{}, errors.New("content can not be empty")
	}
	if message.ClickType != "app" && message.ClickContent == "" {
		return &Response{}, errors.New("ClickContent can not be empty")
	}
	channel := s.Channel
	if channel == "" {
		return &Response{}, errors.New("channel can not be empty")
	}
	pushId := s.PushId
	if len(pushId) == 0 {
		return &Response{}, errors.New("pushId can not be empty")
	}
	plat := s.PlatForm
	if plat == "" {
		return &Response{}, errors.New("please set platform param of channel")
	}

	var mc MobileChannel
	switch channel {
	case "hw":
		mc, err = setHWParam(plat)
		break
	case "ios":
		mc, err = setIOSParam(plat)
		break
	case "mi":
		mc, err = setMIParam(plat)
		break
	case "mz":
		mc, err = setMZParam(plat)
		break
	case "oppo":
		mc, err = setOPPOParam(plat)
		break
	case "vivo":
		mc, err = setVIVOParam(plat)
		break
	}
	if err != nil {
		return &Response{}, err
	}
	return mc.SendMessage(message, pushId)
}

func setHWParam(str string) (*HW, error) {
	var err error
	var param *HW
	_ = json.Unmarshal([]byte(str), &param)
	if param.AppId == "" {
		return param, errors.New("AppId" + "can not be empty")
	}
	if param.ClientSecret == "" {
		return param, errors.New("ClientSecret" + "can not be empty")
	}
	return param, err
}

func setIOSParam(str string) (*IOS, error) {
	var err error
	var param *IOS
	_ = json.Unmarshal([]byte(str), &param)
	if param.KeyId == "" {
		return param, errors.New("KeyId" + "can not be empty")
	}
	if param.TeamId == "" {
		return param, errors.New("TeamId" + "can not be empty")
	}
	if param.BundleId == "" {
		return param, errors.New("BundleId" + "can not be empty")
	}
	if param.AuthTokenPath == "" {
		return param, errors.New("AuthTokenPath" + "can not be empty")
	}
	if param.Bearer == "" {
		param.generateIfExpired()
	}
	return param, err
}

func setMIParam(str string) (*MI, error) {
	var err error
	var param *MI
	_ = json.Unmarshal([]byte(str), &param)
	if param.AppSecret == "" {
		return param, errors.New("AppSecret" + "can not be empty")
	}
	if param.RestrictedPackageName == "" {
		return param, errors.New("RestrictedPackageName" + "can not be empty")
	}
	return param, err
}

func setMZParam(str string) (*MZ, error) {
	var err error
	var param *MZ
	_ = json.Unmarshal([]byte(str), &param)
	if param.AppSecret == "" {
		return param, errors.New("AppSecret" + "can not be empty")
	}
	if param.AppId == "" {
		return param, errors.New("AppId" + "can not be empty")
	}
	return param, err
}

func setOPPOParam(str string) (*OPPO, error) {
	var err error
	var param *OPPO
	_ = json.Unmarshal([]byte(str), &param)
	if param.AppKey == "" {
		return param, errors.New("AppKey" + "can not be empty")
	}
	if param.MasterSecret == "" {
		return param, errors.New("MasterSecret" + "can not be empty")
	}
	return param, err
}

func setVIVOParam(str string) (*VIVO, error) {
	err := isEmpty(str)
	var param *VIVO
	_ = json.Unmarshal([]byte(str), &param)
	if param.AppID == "" {
		return param, errors.New("AppId" + "can not be empty")
	}
	if param.AppKey == "" {
		return param, errors.New("AppKey" + "can not be empty")
	}
	if param.AppSecret == "" {
		return param, errors.New("AppSecret" + "can not be empty")
	}
	if param.AuthToken == "" {
		param.generateIfExpired()
	}
	return param, err
}
