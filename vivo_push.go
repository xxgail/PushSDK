package PushSDK

import (
	"encoding/json"
	"fmt"
	"strconv"
	"sync"
	"time"
)

type VIVO struct {
	sync.Mutex
	AppID     string `json:"app_id"`
	AppKey    string `json:"app_key"`
	AppSecret string `json:"app_secret"`
	AuthToken string `json:"auth_token"`
	IssuedAt  int64  `json:"issued_at"`
}

const (
	VIVOProductUrl   = "https://api-push.vivo.com.cn"
	VIVOAuthUrl      = "/message/auth"
	VIVOSingleSend   = "/message/send"
	VIVOGroupTask    = "/message/saveListPayload"
	VIVOGroupSend    = "/message/pushToList"
	VIVOTokenTimeout = 3600
)

var clickTypeV = map[string]int{
	"app":       1,
	"url":       2,
	"customize": 4,
}

type VFieldSingle struct {
	RegId       string `json:"regId"`
	NotifyType  int    `json:"notifyType"` // 通知类型 1.无、2-响铃、3-振动、4-响铃和振动
	Title       string `json:"title"`
	Content     string `json:"content"`
	TimeToLive  int    `json:"timeToLive"`  // 消息保留时长 单位：秒，取值至少60秒，最长7天。当值为空时，默认一天
	SkipType    int    `json:"skipType"`    // 跳转类型。1-打开APP首页、2-打开链接、3-自定义、4-打开APP内指定页面
	SkipContent string `json:"skipContent"` // 跳转类型为2时，跳转内容最大1000个字符，跳转类型为3或4时，跳转内容最大1024个字符，skipType传3需要在onNotificationMessageClicked回调函数中自己写处理逻辑
	//ClientCustomMap string `json:"clientCustomMap"` // 客户端自定义键值对
	RequestId string `json:"requestId"` // 消息唯一标识
}

const (
	VIVOSuccess = 0
	V10050      = "alias和regId不能为空"
	V10051      = "暂不支持该消息类型"
	V10054      = "notifyType不合法"
	V10055      = "title不能为空"
	V10056      = "title不能超过40个字符串"
)

func (v *VIVO) initMessageSingle(m *Message, pushId string) string {
	if m.ApnsId == "" {
		m.ApnsId = getApnsId()
	}
	vFiled := VFieldSingle{
		RegId:       pushId,
		NotifyType:  1,
		Title:       m.Title,
		Content:     m.Desc,
		TimeToLive:  86400,
		SkipType:    clickTypeV[m.ClickType],
		SkipContent: m.ClickContent,
		RequestId:   m.ApnsId,
	}
	vFiledStr, _ := json.Marshal(vFiled)
	return string(vFiledStr)
}

type GroupMessage struct {
	NotifyType  int    `json:"notifyType"` // 通知类型 1.无、2-响铃、3-振动、4-响铃和振动
	Title       string `json:"title"`
	Content     string `json:"content"`
	TimeToLive  int    `json:"timeToLive"`  // 消息保留时长 单位：秒，取值至少60秒，最长7天。当值为空时，默认一天
	SkipType    int    `json:"skipType"`    // 跳转类型。1-打开APP首页、2-打开链接、3-自定义、4-打开APP内指定页面
	SkipContent string `json:"skipContent"` // 跳转类型为2时，跳转内容最大1000个字符，跳转类型为3或4时，跳转内容最大1024个字符，skipType传3需要在onNotificationMessageClicked回调函数中自己写处理逻辑
	//ClientCustomMap string `json:"clientCustomMap"` // 客户端自定义键值对
	RequestId string `json:"requestId"` // 消息唯一标识
}

func (v *VIVO) initGroupMessage(m *Message) string {
	if m.ApnsId == "" {
		m.ApnsId = getApnsId()
	}
	groupMessage := GroupMessage{
		NotifyType:  1,
		Title:       m.Title,
		Content:     m.Desc,
		TimeToLive:  86400,
		SkipType:    2,
		SkipContent: "http://baidu.com",
		RequestId:   m.ApnsId,
	}
	groupMessageStr, _ := json.Marshal(groupMessage)
	return string(groupMessageStr)
}

type VFieldGroup struct {
	RegIds    []string `json:"regIds"`
	TaskId    string   `json:"taskId"`
	RequestId string   `json:"requestId"`
}

type VResult struct {
	Result      int    `json:"result"`      // 状态码
	Desc        string `json:"desc"`        // 文字描述接口调用情况
	TaskId      string `json:"taskId"`      // 任务编号
	InvalidUser string `json:"invalidUser"` // 非法用户信息
}

func (v *VIVO) SendMessage(m *Message, pushId []string) (*Response, error) {
	response := &Response{}
	header := make(map[string]string)
	header["authToken"] = v.AuthToken
	var body []byte
	if len(pushId) == 1 {
		forms := v.initMessageSingle(m, pushId[0])
		var err error
		body, err = postReqJson(VIVOProductUrl+VIVOSingleSend, forms, header)
		if err != nil {
			response.Code = HTTPERROR
			return response, err
		}
	} else {
		forms := v.initGroupMessage(m)
		bodyTask, err := postReqJson(VIVOProductUrl+VIVOGroupTask, forms, header)
		var resultTask = &VResult{}
		err = json.Unmarshal(bodyTask, &resultTask)
		if err != nil {

		}
		if resultTask.Result != VIVOSuccess {
			fmt.Println("Task Request Error", resultTask.Desc)
		}
		taskId := resultTask.TaskId
		vFieldGroup := VFieldGroup{
			TaskId: taskId,
			RegIds: pushId,
		}
		vFieldGroupStr, _ := json.Marshal(vFieldGroup)
		body, err = postReqJson(VIVOProductUrl+VIVOGroupSend, string(vFieldGroupStr), header)
		if err != nil {

		}
	}
	var result VResult
	fmt.Println("result-vivo", string(body))
	err := json.Unmarshal(body, &result)
	if err != nil {

	}
	if result.Result != VIVOSuccess {
		fmt.Println("Send Message Request Error", result.Desc)
		response.Code = SendError
		response.Reason = result.Desc
	}
	return response, nil
}

type AuthTokenResult struct {
	Result    int    `json:"result"`
	Desc      string `json:"desc"`
	AuthToken string `json:"authToken"`
}

const (
	VIVOAuthSuccess          = 0
	VIVOAuthAppIdEmpty       = 10200
	VIVOAuthAppKeyEmpty      = 10201
	VIVOAuthAppKeyIllegal    = 10202
	VIVOAuthTimestampEmpty   = 10203
	VIVOAuthSignEmpty        = 10204
	VIVOAuthAppIdNotExist    = 10205
	VIVOAuthSignError        = 10206
	VIVOAuthTimestampIllegal = 10207
	VIVOAuthExceededLimit    = 10250
)

func GetAuthTokenV(appId, appKey, appSecret string) string {
	currentTimeStr := strconv.FormatInt(time.Now().UnixNano()/(1e6), 10)
	sign := md5Str(appId + appKey + currentTimeStr + appSecret)
	header := make(map[string]string)
	param := make(map[string]string)
	param["appId"] = appId
	param["appKey"] = appKey
	param["timestamp"] = currentTimeStr
	param["sign"] = sign
	paramStr, _ := json.Marshal(param)
	body, err := postReqJson(VIVOProductUrl+VIVOAuthUrl, string(paramStr), header)
	if err != nil {

	}

	var result = &AuthTokenResult{}
	err = json.Unmarshal(body, result)
	if err != nil {

	}
	if result.Result != VIVOAuthSuccess {
		fmt.Println("AuthToken Request Error")
	}
	return result.AuthToken
}

func (v *VIVO) generateIfExpired() string {
	v.Lock()
	defer v.Unlock()
	if v.Expired() {
		v.Generate()
	}
	return v.AuthToken
}

func (v *VIVO) Expired() bool {
	return time.Now().Unix() >= (v.IssuedAt + VIVOTokenTimeout)
}

func (v *VIVO) Generate() (bool, error) {
	bearer := GetAuthTokenV(v.AppID, v.AppKey, v.AppSecret)
	if bearer != "" {
		v.AuthToken = bearer
		v.IssuedAt = time.Now().Unix()
	}
	return true, nil
}
