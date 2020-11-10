package PushSDK

import (
	"encoding/json"
	"fmt"
	"strconv"
	"time"
)

const (
	VIVOProductUrl = "https://api-push.vivo.com.cn"
	VIVOAuthUrl    = "/message/auth"
	VIVOSingleSend = "/message/send"
	VIVOGroupTask  = "/message/saveListPayload"
	VIVOGroupSend  = "/message/pushToList"
)

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

func initMessageSingleV(m MessageBody, pushId string) *Message {
	vFiled := VFieldSingle{
		RegId:       pushId,
		NotifyType:  1,
		Title:       m.Title,
		Content:     m.Desc,
		TimeToLive:  86400,
		SkipType:    2,
		SkipContent: "http://baidu.com",
		RequestId:   m.ApnsId,
	}
	vFiledStr, _ := json.Marshal(vFiled)
	return &Message{
		Fields: string(vFiledStr),
	}
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

func initGroupMessage(m MessageBody) string {
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

func vSendMessage(m MessageBody, pushId []string, authToken string) (*Response, error) {
	response := &Response{}
	header := make(map[string]string)
	header["authToken"] = authToken
	var body []byte
	if len(pushId) == 1 {
		message := initMessageSingleV(m, pushId[0])
		forms := message.Fields.(string)
		var err error
		fmt.Println(VIVOProductUrl + VIVOSingleSend)
		fmt.Println(forms)
		fmt.Println(header)
		body, err = postReqJson(VIVOProductUrl+VIVOSingleSend, forms, header)
		if err != nil {
			response.Code = HTTPERROR
			return response, err
		}
	} else {
		forms := initGroupMessage(m)
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
	fmt.Println("vvvvvv", string(body))
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
	fmt.Println("vivo-token-body", string(body))
	err = json.Unmarshal(body, result)
	if err != nil {

	}
	if result.Result != VIVOAuthSuccess {
		fmt.Println("AuthToken Request Error")
	}
	return result.AuthToken
}
