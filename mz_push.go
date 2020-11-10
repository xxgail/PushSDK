package PushSDK

import (
	"encoding/json"
	"fmt"
	"sort"
	"strings"
)

type MessageJson struct {
	NoticeBarInfo    NoticeBarInfo    `json:"noticeBarInfo"`
	NoticeExpandInfo NoticeExpandInfo `json:"noticeExpandInfo"`
	ClickTypeInfo    ClickTypeInfo    `json:"clickTypeInfo"`
	PushTimeInfo     PushTimeInfo     `json:"pushTimeInfo"`
}

type NoticeBarInfo struct {
	NoticeBarType int    `json:"noticeBarType"` // 通知栏样式 0-标准 2-安卓原生 【int 非必填，值为0】
	Title         string `json:"title"`         // 推送标题 【string 必填，字数限制1-32字符】
	Content       string `json:"content"`       // 推送内容 【string 必填，字数限制1-100字符】
}

type NoticeExpandInfo struct {
	NoticeExpandType    int    `json:"noticeExpandType"`    // 0-标准、1-文本 int 非必填
	NoticeExpandContent string `json:"noticeExpandContent"` // 展示内容、为文本时必填
}

type ClickTypeInfo struct {
	ClickType       int               `json:"clickType"`       // 点击动作 0-打开应用（默认）、1-打开应用页面、2-打开URI页面、3-应用客户端自定义
	Url             string            `json:"url"`             // clickType=2时必填
	Parameters      map[string]string `json:"parameters"`      // json 格式 非必填
	Activity        string            `json:"activity"`        // clickType=1时必填。格式：pkg.activity
	CustomAttribute string            `json:"customAttribute"` // clickType=3时必填
}

type PushTimeInfo struct {
	OffLine   int `json:"offLine"`   // 是否是离线消息 0-否、1-是（默认）
	ValidTime int `json:"validTime"` // 有效时长（1-72小时内的正整数，默认24
}

//
//type AdvanceInfo struct {
//	Suspend             int              `json:"suspend"`
//	ClearNoticeBar      int              `json:"clearNoticeBar"`
//	FixDisplay          int              `json:"fixDisplay"`
//	FixStartDisplayTime string           `json:"fixStartDisplayTime"`
//	FixEndDisplayTime   string           `json:"fixEndDisplayTime"`
//	NotificationType    NotificationType `json:"notificationType"`
//	NotifyKey           string           `json:"notifyKey"`
//}
//
//type NotificationType struct {
//	Vibrate int `json:"vibrate"`
//	Lights  int `json:"lights"`
//	Sound   int `json:"sound"`
//}

func initMessageMZ(m MessageBody) *Message {
	var messageJson MessageJson
	messageJson = MessageJson{
		NoticeBarInfo: NoticeBarInfo{
			NoticeBarType: 0,
			Title:         m.Title,
			Content:       m.Desc,
		},
		NoticeExpandInfo: NoticeExpandInfo{
			NoticeExpandType:    0,
			NoticeExpandContent: "",
		},
		ClickTypeInfo: ClickTypeInfo{
			ClickType:  0,
			Parameters: map[string]string{},
			Activity:   "",
		},
		PushTimeInfo: PushTimeInfo{
			OffLine:   0, // 是否进离线消息 否 是 【 非必填，默认值为 】
			ValidTime: 0, // 有效时长(1-72小时内的正整数)【 int offline 值为1时，必填，默认24
		},
	}
	messageJsonStr, _ := json.Marshal(messageJson)
	var fields map[string]string
	fields = map[string]string{
		"messageJson": string(messageJsonStr),
	}
	return &Message{
		Fields: fields,
	}
}

const (
	MZProductionHost = "http://server-api-push.meizu.com"

	MZMessageURL = "/garcia/api/server/push/varnished/pushByPushId" // pushId推送
	MZToAppURL   = "/garcia/api/server/push/pushTask/pushToApp"
)

type MZResult struct {
	Code     string `json:"code"`
	Message  string `json:"message"`
	Value    string `json:"value"`
	Redirect string `json:"redirect"`
	MsgId    string `json:"msgId"`
}

func mzMessageSend(m MessageBody, pushIds []string, appId, appSecret string) (*Response, error) {
	response := &Response{}
	message := initMessageMZ(m)
	fields := message.Fields.(map[string]string)
	fields["appId"] = appId
	fields["pushIds"] = strings.Join(pushIds, ",")
	fields["sign"] = generateSign(fields, appSecret)
	requestUrl := MZProductionHost + MZMessageURL
	header := make(map[string]string)
	body, err := postReqUrlencoded(requestUrl, fields, header)
	if err != nil {
		response.Code = HTTPERROR
		return response, err
	}

	var result = &MZResult{}
	err = json.Unmarshal(body, result)
	if err != nil {

	}
	fmt.Println("rrrrrrrrrrrresult", result)
	if result.Code != MZSuccess && result.Value != "" {
		response.Code = SendError
		response.Reason = result.Message
		response.ApnsId = result.MsgId
	}
	return response, nil
}

const (
	MZSuccess        = "200"
	MZOtherError     = "500"
	MZSystemError    = "1001"
	MZServerBusy     = "1003"
	MZParamError     = "1005"
	MZAuthTokenError = "1006"
	MZAppIdIllegal   = "110000"
	MZAppKeyIllegal  = "110001"
	MZParamEmpty     = "110004"
	MZAppInBlackList = "110009"
)

func generateSign(params map[string]string, appSecret string) string {
	keys := make([]string, 0, len(params))
	for k, _ := range params {
		keys = append(keys, k)
	}

	str := ""
	sort.Strings(keys)
	for _, v := range keys {
		str += fmt.Sprintf("%v=%v", v, params[v])
	}
	return md5Str(str + appSecret)
}
