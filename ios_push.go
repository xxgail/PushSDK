package PushSDK

import (
	"crypto/x509"
	"encoding/json"
	"encoding/pem"
	"errors"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"io/ioutil"
	"time"
)

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

func initMessageIOS(title string, desc string) *Message {
	fields := IOSFields{
		Aps: Aps{
			Alert: Alert{
				Title: title,
				Body:  desc,
			},
			Badge: 10,
		},
	}
	fieldsStr, _ := json.Marshal(fields)
	return &Message{
		Fields: string(fieldsStr),
	}
}

const (
	IOSProductionHost = "https://api.development.push.apple.com"
	IOSMessageURL = "/3/device/"
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

func iOSMessagesSend(title string, desc string, token []string, keyId, teamId, bundleId, authTokenPath string) (int, string) {
	code, reason := 1, ""
	message := initMessageIOS(title, desc)
	fields := message.Fields.(string)
	header := make(map[string]string)
	header["apns-topic"] = bundleId
	authToken,err := getAuthToken(authTokenPath, keyId, teamId)
	if err != nil {

	}
	header["Authorization"] = fmt.Sprintf("bearer %s", authToken)

	for _,v := range token {
		requestUrl := IOSProductionHost + IOSMessageURL + v
		fmt.Println("111111111",requestUrl)
		fmt.Println("222222222",fields)
		fmt.Println("333333333",header)
		body,err := postReqJson(requestUrl, fields, header)
		if err != nil {

		}
		if string(body) != "" {
			code = 0
			var errRes = &ErrResult{}
			err = json.Unmarshal(body, errRes)
			reason = errRes.Reason
			break
		}
	}
	return code,reason
}


func getAuthToken(authTokenPath string, keyID string, teamID string) (string, error) {
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