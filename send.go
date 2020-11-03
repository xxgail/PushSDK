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
	HWAppId                 string
	HWClientSecret          string
	IOSKeyId                string
	IOSTeamId               string
	IOSBundleId             string
	IOSAuthTokenPath        string
	MIAppSecret             string
	MIRestrictedPackageName string
	MZAppId                 string
	MZAppSecret             string
	OPPOAppKey              string
	OPPOMasterSecret        string
}

func InitSend(message MessageBody, channel string, pushId []string, platformParam PlatformParam) *Send {
	return &Send{
		MessageBody:   message,
		Channel:       channel,
		PushId:        pushId,
		PlatformParam: platformParam,
	}
}

func (s *Send) SendMessage() (int,string) {
	code,errReason := 0, ""
	switch s.Channel {
	case "hw":
		code, errReason = hwMessagesSend(s.MessageBody.Title, s.MessageBody.Desc, s.PushId, s.PlatformParam.HWAppId, s.PlatformParam.HWClientSecret)
		break
	case "ios":
		code, errReason = iOSMessagesSend(s.MessageBody.Title, s.MessageBody.Desc, s.PushId, s.PlatformParam.IOSKeyId, s.PlatformParam.IOSTeamId, s.PlatformParam.IOSBundleId, s.PlatformParam.IOSAuthTokenPath)
		break
	case "mi":
		code, errReason = miMessageSend(s.MessageBody.Title, s.MessageBody.Desc, s.PushId, s.PlatformParam.MIAppSecret, s.PlatformParam.MIRestrictedPackageName)
		break
	case "mz":
		code, errReason = mzMessageSend(s.MessageBody.Title, s.MessageBody.Desc, s.PushId, s.PlatformParam.MZAppId, s.PlatformParam.MZAppSecret)
		break
	case "oppo":
		code,errReason = oppoMessageSend(s.MessageBody.Title, s.MessageBody.Desc, s.PushId, s.PlatformParam.OPPOAppKey, s.PlatformParam.OPPOMasterSecret)
		break
	}
	return code,errReason
}
