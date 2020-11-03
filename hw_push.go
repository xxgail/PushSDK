package PushSDK

import (
	"encoding/json"
	"fmt"
)

type HWFields struct {
	ValidateOnly bool                `form:"validate_only" json:"validate_only"`
	Message      MessageNotification `json:"message"`
}

type MessageNotification struct {
	Notification Notification `form:"notification" json:"notification"`
	Android      Android      `form:"android" json:"android"`
	Token        []string     `json:"token"`
}

type Notification struct {
	Title string `form:"title" json:"title"`
	Body  string `form:"body" json:"body"`
}

type Android struct {
	Notification AndroidNotification `json:"notification"`
}

type AndroidNotification struct {
	Title       string      `form:"title" json:"title"`
	Body        string      `form:"body" json:"body"`
	ClickAction ClickAction `json:"click_action"`
}

type ClickAction struct {
	Type   int    `json:"type"`
	Intent string `json:"intent"`
}

func initMessageHW(title string, desc string, token []string) *Message {
	fields := HWFields{
		ValidateOnly: false,
		Message: MessageNotification{
			Notification: Notification{
				Title: title,
				Body:  desc,
			},
			Android: Android{
				Notification: AndroidNotification{
					Title: title,
					Body:  desc,
					ClickAction: ClickAction{
						Type:   1,
						Intent: "#Intent;compo=com.rvr/.Activity;S.W=U;end",
					},
				},
			},
			Token: token,
		},
	}
	fieldsStr, _ := json.Marshal(fields)
	return &Message{
		Fields: string(fieldsStr),
	}
}

const (
	HWProductionHost = "https://push-api.cloud.huawei.com/v1/"

	HWMessageURL = "/messages:send"

	HWTokenURL  = "https://oauth-login.cloud.huawei.com/oauth2/v3/token"
	HWGrantType = "client_credentials"
)

const (
	HWSuccess              = "80000000"
	HWPartTokenSendSuccess = "80100000"
)

type HWResult struct {
	Code      string `json:"code"`                //80000000表示成功，非80000000表示失败
	Msg       string `json:"msg"`                 //错误码描述
	RequestId string `json:"requestId,omitempty"` //请求标识。
}

func hwMessagesSend(title string, desc string,token []string, appId, clientSecret string) (int, string) {
	message := initMessageHW(title, desc, token)
	fields := message.Fields.(string)
	requestUrl := HWProductionHost + appId + HWMessageURL
	header := make(map[string]string)
	accessToken := getAccessToken(appId, clientSecret)
	header["Authorization"] = fmt.Sprintf("Bearer %s", accessToken)
	body,err := postReqJson(requestUrl, fields, header)
	var result = &HWResult{}
	err = json.Unmarshal(body, result)
	if err != nil {

	}
	fmt.Println(result)
	if result.Code != HWSuccess {
		return 0, result.Msg
	}
	return 1,""
}

type TokenResult struct {
	AccessToken string `json:"access_token"`
	ExpiresIn   int    `json:"expires_in"`
	TokenType   string `json:"token_type"`
}


func getAccessToken(clientId, clientSecret string) string {
	var accessToken string
	requestPath := HWTokenURL
	forms := make(map[string]string)
	forms["grant_type"] = HWGrantType
	forms["client_id"] = clientId
	forms["client_secret"] = clientSecret
	header := make(map[string]string)
	fmt.Println(requestPath)
	fmt.Println(forms)
	body,err := postReqUrlencoded(requestPath, forms, header)
	var result = &TokenResult{}
	err = json.Unmarshal(body, result)
	if err != nil {
		fmt.Println("getAccessToken Unmarshal", err)
	}
	fmt.Println("请求AccessToken结果：", result)

	accessToken = result.AccessToken

	return accessToken
}