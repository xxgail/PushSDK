# PushSDK
Go-PushSDK
给APNS发推送消息。包括iOS、华为、小米、魅族、OPPO

▶ 待完善

# 具体传输方式
| **channel** | **device_token来源** | **请求方式** | **URL** | **Content-Type** | **request-header** | **消息体结构** | **device_token** |
| --- | --- | --- | --- | --- | --- | --- | --- |
| iOS | device_token | POST | [https://api.development.push.apple.com/3/device/](https://api.development.push.apple.com/3/device/)+token | application/json | "apns-topic":bundleID<br />"Authorization":"bearer "+ authtoken | json格式 | 连接到URL后面 |
| 小米（mi） | regId | POST | [https://api.xmpush.xiaomi.com](https://api.xmpush.xiaomi.com)/v3/message/regid | application/x-www-form-urlencoded | "Authorization": appSecret | key:value格式 | 在消息体中，作为key传输 |
| 华为（hw） | token | POST | [https://push-api.cloud.huawei.com/v1](https://push-api.cloud.huawei.com/v1/)/messages:send | application/json | "Authorization": "Bearer " + authtoken | json格式 | 在json消息体中，数组格式 |
| 魅族（mz） | pushId | POST | [http://server-api-push.meizu.com](http://server-api-push.meizu.com)/garcia/api/server/push/varnished/pushByPushId | application/x-www-form-urlencoded | 没有token，消息体中有sign的key-value | key:value格式 | 在消息体中，作为key传输 |
| oppo | pushId | POST | [https://api.push.oppomobile.com](https://api.push.oppomobile.com)/server/v1/message/notification/unicast_batch | application/x-www-form-urlencoded | "auth_token":  authtoken | json数组格式 | 在messages的json消息体中，分配到每个消息数组中 |

# 推送流程
## iOS
### 方法1 (.p12)

- 提前下好证书（.p12），存储到服务器
- 因为没有用到所以省略
### 方法2（.p8）
#### 流程

1. 根据客户端提供的 .p8文件、keyID、teamID获取authtoken。
```go
import "github.com/dgrijalva/jwt-go"

func GetAuthToken(authTokenPath string, keyID string, teamID string) (string, error) {
	tokenBytes, err := ioutil.ReadFile(authTokenPath)
	if err != nil {
		return "", err
	}
	block, _ := pem.Decode(tokenBytes)
	if block == nil {
		return "", errors.New("Auth token does not seem to be a valid .p8 key file")
	}

	key, err := x509.ParsePKCS8PrivateKey(block.Bytes) // 1. 根据authtokenPath获取到key
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
	} // 2. 进行ES256的jwt加密

	bearer, err := jwtToken.SignedString(key) // 3.标记key
	if err != nil {
		return "", err
	}

	return bearer, nil
}
```

2. 构建消息体。（json字符串的格式）
```go
import "encoding/json"

type Message struct {
	Fields string
}

type Fields struct {
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

func InitMessage(m MessageBody, passThrough string) *Message {
	var fields Fields
	fields = Fields{
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
```

3. 通过http的post的请求方式进行请求
   1. 域名：[https://api.development.push.apple.com/3/device/](https://api.development.push.apple.com/3/device/) + 设备token（device_token）
   1. 请求header
      1. "Content-Type" : "application/json"
      1. "apns-topic" : 包名
      1. "Authorization" : "bearer"+" " + authToken
   3. 请求参数：json格式的消息体
   3. 返回结果：
      1. 请求成功则返回nil
      1. 否则返回json {"reason": }
```go
func postReq(requestPath string, fields string, authToken string, bundleID string) (*Result, error) {
	req, err := http.NewRequest("POST", fmt.Sprintf("%s", baseHost()+requestPath), strings.NewReader(fields))
	if err != nil {
		return nil, err
	} // 创建post请求

	req.Header.Set("Content-Type", "application/json") // 设置header
	req.Header.Add("apns-topic", bundleID)
	req.Header.Add("Authorization", "bearer"+" "+authToken)

	var client = &http.Client{
		Timeout: 5 * time.Second,
	}

	res, err := client.Do(req)
	apnsId := res.Header.Get("apns-id")
	body, err := ioutil.ReadAll(res.Body)
	res.Body.Close()

	if err != nil {
	}

    // type ErrResult struct { Reason string `json:"reason"` }
	var errRes = &ErrResult{}
	err = json.Unmarshal(body, errRes)
	var result = &Result{}
	if string(body) == "" {
		result.Reason = ""
		result.Status = Success
	} else {
		result.Reason = errRes.Reason
		result.Status = ErrorRequest
	}
	result.ApnsId = apnsId
	return result, nil
}
```



## 安卓
### 小米

1. 设置AppSecret
1. 构建消息体，发送方式为regId
```go
fields = map[string]interface{}{
		"registration_id":         "",
		"title":                   title,
		"description":             desc,
		"restricted_package_name": restrictedPackageName,
		"payload":                 payloadStr,
		"notify_type":             "-1",
		"pass_through":            passThrough,
}
```

3. 请求
   1. 请求方式：post
   1. 参数：fields中的key和value
   1. Content-Type："application/x-www-form-urlencoded;charset=UTF-8"
   1. 设置头信息：Add("Authorization", fmt.Sprintf("key=%s", appSecret))`
### 华为

1. 获取AccessToken
1. 构建消息体
```go
type Fields struct {
	ValidateOnly  	bool `form:"validate_only" json:"validate_only"`
	Message			MessageNotification 	`json:"message"`
}

type MessageNotification struct {
	Notification Notification `form:"notification" json:"notification"`
	Android		 Android	`form:"android" json:"android"`
	Token		[]string	`json:"token"`
}

type Notification struct {
	Title		string		`form:"title" json:"title"`
	Body		string		`form:"body" json:"body"`
}

type Android struct {
	Notification	AndroidNotification 	`json:"notification"`
}

type AndroidNotification struct {
	Title		string		`form:"title" json:"title"`
	Body		string		`form:"body" json:"body"`
	ClickAction	ClickAction  `json:"click_action"`
}

type ClickAction struct {
	Type		int 	`json:"type"`
	Intent		string	`json:"intent"`
}
```

3. 初始化消息field之后转为json字符串
3. 请求
   1. 请求方式：post
   1. 参数：field（json格式）
   1. Content-Type："application/json"
   1. 设置头信息：Add("Authorization", "Bearer" + " " + accessToken)
### 魅族
### OPPO


# 参考来源
> ios- Local and Remote Notification Programming Guide [https://developer.apple.com/library/archive/documentation/NetworkingInternet/Conceptual/RemoteNotificationsPG/CommunicatingwithAPNs.html#//apple_ref/doc/uid/TP40008194-CH11-SW1](https://developer.apple.com/library/archive/documentation/NetworkingInternet/Conceptual/RemoteNotificationsPG/CommunicatingwithAPNs.html#//apple_ref/doc/uid/TP40008194-CH11-SW1)
