package PushSDK

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strings"
)

type MI struct {
	AppSecret             string `json:"app_secret"`
	RestrictedPackageName string `json:"restricted_package_name"`
}

const (
	MiProductionHost  = "https://api.xmpush.xiaomi.com"
	MiMessageRegIdURL = "/v3/message/regid"
	MiMultiRegIdURL   = "/v2/multi_messages/regids"
)

type MIFields struct {
	Payload               string `json:"payload"`                 //消息内容
	RestrictedPackageName string `json:"restricted_package_name"` //支持多包名
	Title                 string `json:"title"`                   //在通知栏的标题，长度小于16
	Description           string `json:"description"`             //在通知栏的描述，长度小于128
	NotifyType            string `json:"notify_type"`             //通知类型 可组合 (-1 Default_all,1 Default_sound,2 Default_vibrate(震动),4 Default_lights)
	NotifyId              string `json:"notify_id"`               //同一个notifyId在通知栏只会保留一条
	TimeToLive            string `json:"time_to_live"`            //可选项，当用户离线是，消息保留时间，默认两周，单位ms
	TimeToSend            string `json:"time_to_send"`            //可选项，定时发送消息，用自1970年1月1日以来00:00:00.0 UTC时间表示（以毫秒为单位的时间）。
	Extra                 string `json:"extra"`
}

var clickTypeMi = map[string]int{
	"app":       1,
	"url":       3,
	"customize": 2,
}

type Extra struct {
	NotifyEffect int    `json:"notify_effect"` // “1″：通知栏点击后打开app的Launcher Activity。 “2″：通知栏点击后打开app的任一Activity（开发者还需要传入extra.intent_uri）。 “3″：通知栏点击后打开网页（开发者还需要传入extra.web_uri）。
	IntentUri    string `json:"intent_uri"`
	WebUri       string `json:"web_uri"`
}

func (mi *MI) initMessage(m *Message) map[string]string {
	var payload = &Payload{
		PushTitle:    m.Title,
		PushBody:     m.Desc,
		IsShowNotify: "1",
		Ext:          "",
	}
	payloadStr, _ := json.Marshal(payload)
	var extra = Extra{
		NotifyEffect: clickTypeMi[m.ClickType],
		IntentUri:    m.ClickContent,
		WebUri:       m.ClickContent,
	}
	extraStr, _ := json.Marshal(extra)
	fields := MIFields{
		Payload:     string(payloadStr),
		Title:       m.Title,
		Description: m.Desc,
		NotifyType:  "-1",
		NotifyId:    m.ApnsId,
		Extra:       string(extraStr),
	}
	t := reflect.TypeOf(fields)
	v := reflect.ValueOf(fields)
	fieldsMap := make(map[string]string)
	for k := 0; k < t.NumField(); k++ {
		fieldsMap[t.Field(k).Name] = v.Field(k).Interface().(string)
	}
	return fieldsMap
}

type MIResult struct {
	Code        int64       `json:"code"`                  //0表示成功，非0表示失败
	Result      string      `json:"result"`                //"ok" 表示成功,"error" 表示失败
	Description string      `json:"description,omitempty"` //对发送消息失败原因的解释
	Info        string      `json:"info,omitempty"`        //详细信息
	Reason      string      `json:"reason,omitempty"`      //失败原因
	Data        *ResultData `json:"data,omitempty"`        //本身就是一个json字符串
}

type ResultData struct {
	BadRegIds string `json:"bad_regIds"` //推送失败的ids
	Id        string `json:"id"`         //消息的Id
}

const (
	MiSuccess             = 0
	MiNetWorkErrorTimeOut = -1
)

//消息payload，根据业务自定义
type Payload struct {
	PushTitle    string `json:"push_title"`
	PushBody     string `json:"push_body"`
	IsShowNotify string `json:"is_show_notify"`
	Ext          string `json:"ext"`
}

func (mi *MI) SendMessage(m *Message, regIds []string) (*Response, error) {
	response := &Response{}
	forms := mi.initMessage(m)
	forms["registration_id"] = strings.Join(regIds, ",")
	forms["restricted_package_name"] = mi.RestrictedPackageName

	requestUrl := MiProductionHost + MiMessageRegIdURL

	header := make(map[string]string)
	header["Authorization"] = fmt.Sprintf("key=%s", mi.AppSecret)

	body, err := postReqUrlencoded(requestUrl, forms, header)
	if err != nil {
		response.Code = HTTPERROR
		return response, err
	}
	fmt.Println("result-mi", string(body))
	var result = &MIResult{}
	err = json.Unmarshal(body, result)
	if err != nil {

	}
	if result.Code != MiSuccess {
		response.Code = SendError
		response.Reason = result.Reason
	}
	return response, nil
}
