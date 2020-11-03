package PushSDK

import (
	"encoding/json"
	"fmt"
	"strconv"
	"time"
)

type MessageFields struct {
	TargetType   int          `json:"target_type"`
	TargetValue  string       `json:"target_value"`
	Notification OPPONotification `json:"notification"`
}

type OPPONotification struct {
	AppMessageId string `json:"app_message_id"`
	Style        int    `json:"style"` // 1-标准样式（默认1）
	Title        string `json:"title"`
	SubTitle     string `json:"sub_title"`
	Content      string `json:"content"`
}

func initMessageOPPO(title string, desc string, registrationIds []string) *Message {
	var messages []MessageFields
	for _, v := range registrationIds {
		message := MessageFields{
			TargetType:  2,
			TargetValue: v,
			Notification: OPPONotification{
				Style:    1,
				Title:    title,
				SubTitle: title,
				Content:  desc,
			},
		}
		messages = append(messages, message)
	}
	messagesStr,_ := json.Marshal(messages)
	fields := map[string]string{
		"messages" : string(messagesStr),
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

func oppoMessageSend(title string, desc string, pushIds []string, appKey, masterSecret string) (int, string) {
	message := initMessageOPPO(title, desc, pushIds)
	fields := message.Fields.(map[string]string)
	requestUrl := OPPOProductionHost + OPPOMessageURL
	header := make(map[string]string)
	header["auth_token"],_ = getAuthTokenOPPO(appKey, masterSecret)
	body, err := postReqUrlencoded(requestUrl, fields, header)
	if err != nil {

	}

	var result = &OPPOResult{}
	err = json.Unmarshal(body, result)
	if err != nil {

	}
	fmt.Println("rrrrrrrrrrrresult",result)
	if result.Code != OPPOSuccess{
		return 0,result.Message
	}
	return 1,""
}

type OPPOTokenResult struct {
	Code    int       `json:"code"`
	Message string    `json:"message"`
	Data    OPPOTokenData `json:"data"`
}

type OPPOTokenData struct {
	AuthToken  string `json:"auth_token"`
	CreateTime int    `json:"create_time"`
}

func getAuthTokenOPPO(appKey string, masterSecret string) (string, error) {
	var authToken string
	// 毫秒级
	currentTimeStr := strconv.FormatInt(time.Now().UnixNano()/(1e6), 10)
	requestUrl := OPPOProductionHost + OPPOTokenURL
	forms := make(map[string]string)
	forms["app_key"] = appKey
	forms["sign"] = sha256Encode(appKey+currentTimeStr+masterSecret)
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