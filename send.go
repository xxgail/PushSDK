package PushSDK

type Send struct {
	MessageBody   MessageBody
	Channel       string
	PushId        []string
	PlatformParam PlatformParam
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

func InitSend(message MessageBody, channel string, pushId []string, platformParam PlatformParam) *Send {
	return &Send{
		MessageBody:   message,
		Channel:       channel,
		PushId:        pushId,
		PlatformParam: platformParam,
	}
}

func (s *Send) SendMessage() (int, string) {
	code, errReason := 0, ""
	switch s.Channel {
	case "hw":
		code, errReason = hwMessagesSend(s.MessageBody.Title, s.MessageBody.Desc, s.PushId, s.PlatformParam.HWAppId, s.PlatformParam.HWClientSecret)
		break
	case "ios":
		code, errReason = iOSMessagesSend(s.MessageBody.Title, s.MessageBody.Desc, s.PushId, s.PlatformParam.IOSBundleId, s.PlatformParam.IOSAuthToken)
		break
	case "mi":
		code, errReason = miMessageSend(s.MessageBody.Title, s.MessageBody.Desc, s.PushId, s.PlatformParam.MIAppSecret, s.PlatformParam.MIRestrictedPackageName)
		break
	case "mz":
		code, errReason = mzMessageSend(s.MessageBody.Title, s.MessageBody.Desc, s.PushId, s.PlatformParam.MZAppId, s.PlatformParam.MZAppSecret)
		break
	case "oppo":
		code, errReason = oppoMessageSend(s.MessageBody.Title, s.MessageBody.Desc, s.PushId, s.PlatformParam.OPPOAppKey, s.PlatformParam.OPPOMasterSecret)
		break
	case "vivo":
		code, errReason = vSendMessage(s.MessageBody, s.PushId, s.PlatformParam.VIAuthToken)
		break
	}
	return code, errReason
}
