package PushSDK

import (
	"encoding/json"
	"errors"
	"fmt"
	"sync"
)

type HW struct {
	sync.Mutex
	AppId        string `json:"app_id"`
	ClientSecret string `json:"client_secret"`
}

type HWFields struct {
	Message MessageNotification `json:"message"`
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
	CollapseKey  int                 `json:"collapse_key"` // 缓存。0-只缓存新的一条，-1对所有离线消息都缓存、1-100分组缓存
	Notification AndroidNotification `json:"notification"`
}

type AndroidNotification struct {
	Title       string      `form:"title" json:"title"`
	Body        string      `form:"body" json:"body"`
	Icon        string      `json:"icon"` // 小图标
	Tag         string      `json:"tag"`  // 消息标签
	ClickAction ClickAction `json:"click_action"`
	//Badge       BadgeNotification `json:"badge"` // 角标
}

var clickTypeHW = map[string]int{
	"app":       3,
	"url":       2,
	"customize": 1,
}

type ClickAction struct {
	Type   int    `json:"type"`   // 1-打开应自定义页面、2-打开URL、3-打开应用APP
	Intent string `json:"intent"` // 自定义页面中intent的实现
	Url    string `json:"url"`
	Action string `json:"action"` // 设置通过action打开应用自定义页面时，本字段填写要打开的页面activity对应的action。
}

type BadgeNotification struct {
	AddNum int    `json:"add_num"`
	Class  string `json:"class"`   // 应用入口Activity类全路径。 样例：com.example.hmstest.MainActivity
	SetNum int    `json:"set_num"` // 角标设置数字，大于等于0小于100的整数。如果set_num与add_num同时存在时，以set_num为准
}

func (h *HW) initMessage(m *MessageBody, token []string) *Message {
	fields := HWFields{
		Message: MessageNotification{
			Notification: Notification{
				Title: m.Title,
				Body:  m.Desc,
			},
			Android: Android{
				Notification: AndroidNotification{
					Title: m.Title,
					Body:  m.Desc,
					ClickAction: ClickAction{
						Type:   clickTypeHW[m.ClickType],
						Url:    m.ClickContent,
						Intent: "#Intent;compo=com.rvr/.Activity;S.W=U;end",
					},
					Tag: m.ApnsId,
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

func (h *HW) SendMessage(m *MessageBody, token []string) (*Response, error) {
	response := &Response{}
	message := h.initMessage(m, token)
	fields := message.Fields.(string)
	if h.AppId == "" {
		return response, errors.New("AppId 不能为空")
	}
	requestUrl := HWProductionHost + h.AppId + HWMessageURL
	header := make(map[string]string)
	accessToken := getAccessToken(h.AppId, h.ClientSecret)
	header["Authorization"] = fmt.Sprintf("Bearer %s", accessToken)
	body, err := postReqJson(requestUrl, fields, header)
	fmt.Println("result-hw", string(body))
	if err != nil {
		response.Code = HTTPERROR
		return response, err
	}
	var result = &HWResult{}
	err = json.Unmarshal(body, result)
	if err != nil {

	}

	if result.Code != HWSuccess {
		response.Code = SendError
		response.Reason = result.Msg
		response.ApnsId = result.RequestId
	}
	return response, err
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
	body, err := postReqUrlencoded(requestPath, forms, header)
	var result = &TokenResult{}
	err = json.Unmarshal(body, result)
	if err != nil {
		fmt.Println("getAccessToken Unmarshal", err)
	}

	accessToken = result.AccessToken

	return accessToken
}
