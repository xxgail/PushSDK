package PushSDK

import (
	"crypto/x509"
	"encoding/json"
	"encoding/pem"
	"errors"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"io/ioutil"
	"sync"
	"time"
)

type IOS struct {
	sync.Mutex
	KeyId         string `json:"key_id"`
	TeamId        string `json:"team_id"`
	BundleId      string `json:"bundle_id"`
	AuthTokenPath string `json:"auth_token_path"`
	Bearer        string `json:"bearer"`
	IssuedAt      int64  `json:"issued_at"`
}

type IOSFields struct {
	Aps Aps `json:"aps"`
}

type Aps struct {
	Alert Alert `json:"alert"`
	Badge int   `json:"badge"` // 右上角小图标
}

type Alert struct {
	Title string `json:"title"`
	Body  string `json:"body"`
}

func (i *IOS) initMessage(m *Message) string {
	fields := IOSFields{
		Aps: Aps{
			Alert: Alert{
				Title: m.Title,
				Body:  m.Desc,
			},
			Badge: 1,
		},
	}
	fieldsStr, _ := json.Marshal(fields)
	return string(fieldsStr)
}

const (
	IOSProductionHost = "https://api.development.push.apple.com"
	IOSMessageURL     = "/3/device/"
	IOSTokenTimeout   = 3000
)

const (
	IOSSuccess                     = 200 // 成功
	IOSErrorRequest                = 400 // 错误的请求
	IOSErrorAuthToken              = 403
	IOSErrorMethod                 = 405
	IOSTimeOutDeviceToken          = 410
	IOSNotificationPayloadTooLarge = 413
	IOSTooManyRequest              = 429
	IOSInternalServerError         = 500
	IOSServerUnavailable           = 503
)

type IOSResult struct {
	ApnsId string `json:"apns-id"`
	Status int    `json:"status"`
	Reason string `json:"reason"`
}

type ErrResult struct {
	Reason string `json:"reason"`
}

func (i *IOS) SendMessage(m *Message, token []string) (*Response, error) {
	response := &Response{}
	forms := i.initMessage(m)
	header := make(map[string]string)
	header["apns-topic"] = i.BundleId
	header["Authorization"] = fmt.Sprintf("bearer %s", i.Bearer)
	header["apns-id"] = m.ApnsId

	for _, v := range token {
		requestUrl := IOSProductionHost + IOSMessageURL + v
		body, err := postReqJson(requestUrl, forms, header)
		if err != nil {
			response.Code = HTTPERROR
			return response, err
		}
		fmt.Println("result-ios", string(body))
		if string(body) != "" {
			var errRes = &ErrResult{}
			err = json.Unmarshal(body, errRes)
			response.Code = SendError
			response.Reason = errRes.Reason
			break
		}
	}
	return response, nil
}

func GetAuthTokenIOS(authTokenPath string, keyID string, teamID string) (string, error) {
	tokenBytes, err := ioutil.ReadFile(authTokenPath)
	if err != nil {
		return "", err
	}
	block, _ := pem.Decode(tokenBytes)
	if block == nil {
		return "", errors.New("Auth token does not seem to be a valid .p8 key file")
	}

	key, err := x509.ParsePKCS8PrivateKey(block.Bytes)
	if err != nil {
		return "", err
	}

	jwtToken := &jwt.Token{
		Header: map[string]interface{}{
			"alg": "ES256",
			"kid": keyID,
		},
		Claims: jwt.MapClaims{
			"iss": teamID,
			"iat": time.Now().Unix(),
		},
		Method: jwt.SigningMethodES256,
	}

	bearer, err := jwtToken.SignedString(key)
	if err != nil {
		return "", err
	}

	return bearer, nil
}

func (i *IOS) generateIfExpired() string {
	i.Lock()
	defer i.Unlock()
	if i.Expired() {
		i.Generate()
	}
	return i.Bearer
}

func (i *IOS) Expired() bool {
	return time.Now().Unix() >= (i.IssuedAt + IOSTokenTimeout)
}

func (i *IOS) Generate() (bool, error) {
	bearer, err := GetAuthTokenIOS(i.AuthTokenPath, i.KeyId, i.TeamId)
	if bearer != "" {
		i.Bearer = bearer
		i.IssuedAt = time.Now().Unix()
	}
	return true, err
}
