package PushSDK

import (
	"encoding/json"
	"fmt"
	"strings"
)

const (
	MiProductionHost  = "https://api.xmpush.xiaomi.com"
	MiMessageRegIdURL = "/v3/message/regid"
	MiMultiRegIdURL   = "/v2/multi_messages/regids"
)

type MIFields struct {
	Payload               string `json:"payload"`                 //消息内容
	RestrictedPackageName string `json:"restricted_package_name"` //支持多包名
	PassThrough           string `json:"pass_through"`            //是否透传给app(1 透传 0 通知栏信息)
	Title                 string `json:"title"`                   //在通知栏的标题，长度小于16
	Description           string `json:"description"`             //在通知栏的描述，长度小于128
	NotifyType            string `json:"notify_type"`             //通知类型 可组合 (-1 Default_all,1 Default_sound,2 Default_vibrate(震动),4 Default_lights)
	NotifyId              string `json:"notify_id"`               //同一个notifyId在通知栏只会保留一条
	TimeToLive            string `json:"time_to_live"`            //可选项，当用户离线是，消息保留时间，默认两周，单位ms
	TimeToSend            string `json:"time_to_send"`            //可选项，定时发送消息，用自1970年1月1日以来00:00:00.0 UTC时间表示（以毫秒为单位的时间）。
}

func initMessageMi(title string, desc string) *Message {
	var payload = &Payload{
		PushTitle:    title,
		PushBody:     desc,
		IsShowNotify: "1",
		Ext:          "",
	}
	payloadStr,_ := json.Marshal(payload)
	fields := MIFields{
		Payload:     string(payloadStr),
		PassThrough: "0",
		Title:       title,
		Description: desc,
		NotifyType:  "-1",
	}
	fieldsStr, _ := json.Marshal(fields)
	return &Message{
		Fields: string(fieldsStr),
	}
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

func miMessageSend(title string, desc string, regIds []string, appSecret, restrictedPackageName string) (int, string) {
	message := initMessageMi(title, desc)
	fieldsStr := message.Fields.(string)
	var fields map[string]string
	err := json.Unmarshal([]byte(fieldsStr), &fields)
	if err != nil {
		fmt.Println("field json.Unmarshal error", err)
	}
	fields["registration_id"] = strings.Join(regIds, ",")
	fields["restricted_package_name"] = restrictedPackageName

	requestUrl := MiProductionHost + MiMessageRegIdURL

	header := make(map[string]string)
	header["Authorization"] = fmt.Sprintf("key=%s", appSecret)

	body, err := postReqUrlencoded(requestUrl, fields, header)
	if err != nil {

	}

	var result = &MIResult{}
	err = json.Unmarshal(body, result)
	if err != nil {

	}
	fmt.Println("rrrrrrrrrrrresult",result)
	if result.Code != MiSuccess{
		return 0,result.Reason
	}
	return 1,""
}