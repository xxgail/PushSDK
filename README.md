# PushSDK
Go-PushSDK
给APNS发推送消息。包括iOS、华为、小米、魅族、OPPO、vivo

▶ 待完善

🍬 Add the library to your $GOPATH/src
`go get github.com/xxgail/PushSDK`

- [builder.go](https://github.com/xxgail/PushSDK/blob/master/builder.go) 构建消息体
- [common.go](https://github.com/xxgail/PushSDK/blob/master/common.go) 包内的公共方法
- [const.go](https://github.com/xxgail/PushSDK/blob/master/const.go) 定义常量
- [http_request.go](https://github.com/xxgail/PushSDK/blob/master/http_request.go) http请求公共方法
- xx_push.go 各个渠道的推送具体请求
- [result_common.go](https://github.com/xxgail/PushSDK/blob/master/result_common.go) 返回格式
- [send.go](https://github.com/xxgail/PushSDK/blob/master/send.go) sendMessage主体步骤

# example
```go
package main

import (
    "fmt"
	"github.com/xxgail/PushSDK"
)

func main() {
    // 1. 推送消息体
    message := PushSDK.NewMessage()
    message.SetTitle("title").SetContent("content")
    // 2. 发送
    send := PushSDK.NewSend()
    send.SetChannel("ios") // 发送渠道，全部小写
    send.SetPushId([]string{"xxxx"}) // 发送用户device_token，数组格式
    send.SetPlatForm("{app_id:xxxx}") // 渠道对应参数，详见下表 channel-param
    response,_ := send.SendMessage(message) // 发送
    fmt.Println(response)
}

```
# channel-param
| **channel** | **plat** |
| ios | {"key_id":"xxx","team_id":"xxx","bundle_id":"xxx","auth_token_path":"xxx.p8"} |
| mi | {"app_secret":"xxx","restricted_package_name":"xxx"} |
| hw | {"app_id":"xxx","client_secret":"xxx"} |
| mz | {"app_id":"xxx","app_secret":"xxx"} |
| oppo | {"app_key":"xxx","master_secret":"xxx"} |
| vivo | {"app_id":"xxx","app_key":"xxx","app_secret":"xxx"} |

# 具体传输方式
| **channel** | **device_token来源** | **请求方式** | **URL** | **Content-Type** | **request-header** | **消息体结构** | **device_token位置** |
| --- | --- | --- | --- | --- | --- | --- | --- |
| iOS | device_token | POST | [https://api.development.push.apple.com/3/device/](https://api.development.push.apple.com/3/device/)+token | application/json | "apns-topic":bundleID<br />"Authorization":"bearer "+ authtoken | json格式 | 连接到URL后面 |
| 小米 | regId | POST | [https://api.xmpush.xiaomi.com](https://api.xmpush.xiaomi.com)/v3/message/regid | application/x-www-form-urlencoded | "Authorization": appSecret | key:value格式 | 在消息体中，作为key传输 |
| 华为 | token | POST | [https://push-api.cloud.huawei.com/v1](https://push-api.cloud.huawei.com/v1/)/messages:send | application/json | "Authorization": "Bearer " + authtoken | json格式 | 在json消息体中，数组格式 |
| 魅族 | pushId | POST | [http://server-api-push.meizu.com](http://server-api-push.meizu.com)/garcia/api/server/push/varnished/pushByPushId | application/x-www-form-urlencoded | 没有token，消息体中有sign的key-value | key:value格式 | 在消息体中，作为key传输 |
| oppo | pushId | POST | [https://api.push.oppomobile.com](https://api.push.oppomobile.com)/server/v1/message/notification/unicast_batch | application/x-www-form-urlencoded | "auth_token": authtoken | key:value格式 | 在messages的json消息体中，分配到每个消息数组中 |
| vivo| regId | POST| [https://api-push.vivo.com.cn](https://api-push.vivo.com.cn)/message/send | application/json | "authToken" = authToken | json格式 | 在json消息体中 |


# 参考来源
> ios- Local and Remote Notification Programming Guide [https://developer.apple.com/library/archive/documentation/NetworkingInternet/Conceptual/RemoteNotificationsPG/CommunicatingwithAPNs.html#//apple_ref/doc/uid/TP40008194-CH11-SW1](https://developer.apple.com/library/archive/documentation/NetworkingInternet/Conceptual/RemoteNotificationsPG/CommunicatingwithAPNs.html#//apple_ref/doc/uid/TP40008194-CH11-SW1)
