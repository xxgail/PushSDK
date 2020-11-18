package PushSDK

import (
	"encoding/json"
	"fmt"
	"strconv"
	"time"
)

type OPPO struct {
	AppKey       string `json:"app_key"`
	MasterSecret string `json:"master_secret"`
}

type MessageFields struct {
	TargetType   int              `json:"target_type"`
	TargetValue  string           `json:"target_value"`
	Notification OPPONotification `json:"notification"`
}

var clickTypeOPPO = map[string]int{
	"app":       0,
	"url":       2,
	"customize": 1,
}

type OPPONotification struct {
	AppMessageId        string `json:"app_message_id"`   // 消息tag
	Style               int    `json:"style"`            // 1-标准样式（默认1） 2-长文本 3-大图
	BigPictureId        string `json:"big_picture_id"`   // 大图id
	SmallPictureId      string `json:"small_picture_id"` // 图标
	Title               string `json:"title"`
	SubTitle            string `json:"sub_title"`
	Content             string `json:"content"`
	ClickActionType     int    `json:"click_action_type"`     //点击动作类型0，启动应用；1，打开应用内页（activity的intent action）；2，打开网页；4，打开应用内页（activity）；【非必填，默认值为0】;5,Intent scheme URL
	ClickActionActivity string `json:"click_action_activity"` // 应用内页地址【click_action_type为1/4/ 时必填，长度500】
	ClickActionUrl      string `json:"click_action_url"`      // 网页地址或【click_action_type为2与5时必填
	ActionParameters    string `json:"action_parameters"`     // 传递给网页或应用的参数 json 格式
}

func (o *OPPO) initMessage(m *MessageBody, registrationIds []string) *Message {
	var messages []MessageFields
	for _, v := range registrationIds {
		message := MessageFields{
			TargetType:  2,
			TargetValue: v,
			Notification: OPPONotification{
				AppMessageId:        m.ApnsId,
				Style:               1,
				Title:               m.Title,
				SubTitle:            m.Title,
				Content:             m.Desc,
				ClickActionType:     clickTypeOPPO[m.ClickType],
				ClickActionActivity: m.ClickContent,
				ClickActionUrl:      m.ClickContent,
			},
		}
		messages = append(messages, message)
	}
	messagesStr, _ := json.Marshal(messages)
	fields := map[string]string{
		"messages": string(messagesStr),
	}
	return &Message{
		Fields: fields,
	}
}

const (
	OPPOProductionHost = "https://api.push.oppomobile.com"
	OPPOMessageURL     = "/server/v1/message/notification/unicast_batch" // pushId推送
	OPPOTokenURL       = "/server/v1/auth"
)

const (
	OPPOSuccess                     = 0
	OPPOServiceInFlowControl        = -2
	OPPOServiceCurrentlyUnavailable = -1
	OPPOInvalidAuthToken            = 11
	OPPOHttpActionNotAllowed        = 12
	OPPOAppCallLimited              = 13
)

type OPPOResult struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    []Data `json:"data"`
}

type Data struct {
	MessageId      string `json:"messageId"`
	RegistrationId string `json:"registrationId"`
	ErrorCode      int    `json:"errorCode"`
	ErrorMessage   string `json:"errorMessage"`
}

func (o *OPPO) SendMessage(m *MessageBody, pushIds []string) (*Response, error) {
	response := &Response{}
	message := o.initMessage(m, pushIds)
	fields := message.Fields.(map[string]string)
	requestUrl := OPPOProductionHost + OPPOMessageURL
	header := make(map[string]string)
	header["auth_token"], _ = o.getAuthTokenOPPO()
	body, err := postReqUrlencoded(requestUrl, fields, header)
	if err != nil {
		response.Code = HTTPERROR
		return response, err
	}
	fmt.Println("result-oppo", string(body))
	var result = &OPPOResult{}
	err = json.Unmarshal(body, result)
	if err != nil {

	}
	if result.Code != OPPOSuccess {
		response.Code = SendError
		response.Reason = result.Message
	}
	return response, nil
}

type OPPOTokenResult struct {
	Code    int           `json:"code"`
	Message string        `json:"message"`
	Data    OPPOTokenData `json:"data"`
}

type OPPOTokenData struct {
	AuthToken  string `json:"auth_token"`
	CreateTime int    `json:"create_time"`
}

func (o *OPPO) getAuthTokenOPPO() (string, error) {
	var authToken string
	// 毫秒级
	currentTimeStr := strconv.FormatInt(time.Now().UnixNano()/(1e6), 10)
	requestUrl := OPPOProductionHost + OPPOTokenURL
	forms := make(map[string]string)
	forms["app_key"] = o.AppKey
	forms["sign"] = sha256Encode(o.AppKey + currentTimeStr + o.MasterSecret)
	forms["timestamp"] = currentTimeStr
	header := make(map[string]string)

	body, err := postReqUrlencoded(requestUrl, forms, header)
	if err != nil {

	}

	var result = &OPPOTokenResult{}
	err = json.Unmarshal(body, result)
	if err != nil {
		fmt.Println("getAuthToken Unmarshal", err)
	}

	if result.Code != 0 {
		fmt.Println("请求AuthToken结果错误：", result)
	}

	authToken = result.Data.AuthToken

	return authToken, err
}
