package PushSDK

import (
	"crypto/md5"
	"encoding/hex"
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
	NoticeBarType int    `json:"noticeBarType"`
	Title         string `json:"title"`
	Content       string `json:"content"`
}

type NoticeExpandInfo struct {
	NoticeExpandType    int    `json:"noticeExpandType"` // 0-标准、1-文本
	NoticeExpandContent string `json:"noticeExpandContent"`
}

type ClickTypeInfo struct {
	ClickType       int               `json:"clickType"` // 点击动作 0-打开应用（默认）、1-打开应用页面、2-打开URI页面、3-应用客户端自定义
	Url             string            `json:"url"`
	Parameters      map[string]string `json:"parameters"`
	Activity        string            `json:"activity"`
	CustomAttribute string            `json:"customAttribute"`
}

type PushTimeInfo struct {
	OffLine   int `json:"offLine"`
	ValidTime int `json:"validTime"`
}

func initMessageMZ(title string, desc string) *Message {
	var messageJson MessageJson
	messageJson = MessageJson{
		NoticeBarInfo: NoticeBarInfo{
			NoticeBarType: 0,
			Title:         title,
			Content:       desc,
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
			OffLine:   0,
			ValidTime: 0,
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


func mzMessageSend(title string, desc string, pushIds []string,appId,appSecret string) (int, string) {
	message := initMessageMZ(title, desc)
	fields := message.Fields.(map[string]string)
	fields["appId"] = appId
	fields["pushIds"] = strings.Join(pushIds, ",")
	fields["sign"] = generateSign(fields, appSecret)
	requestUrl := MZProductionHost + MZMessageURL
	header := make(map[string]string)
	body, err := postReqUrlencoded(requestUrl, fields, header)
	if err != nil {

	}

	var result = &MZResult{}
	err = json.Unmarshal(body, result)
	if err != nil {

	}
	fmt.Println("rrrrrrrrrrrresult",result)
	if result.Code != MZSuccess{
		return 0,result.Message
	}
	return 1,""
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
	str += appSecret
	u := md5.New()
	u.Write([]byte(str))
	return hex.EncodeToString(u.Sum(nil))
}