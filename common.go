package PushSDK

import (
	"crypto/md5"
	"crypto/sha256"
	"encoding/hex"
	"strconv"
	"time"
)

// 公共方法
func sha256Encode(str string) string {
	hash := sha256.New()
	//输入数据
	hash.Write([]byte(str))
	//计算哈希值
	bytes := hash.Sum(nil)
	//将字符串编码为16进制格式,返回字符串
	return hex.EncodeToString(bytes)
}

func md5Str(str string) string {
	u := md5.New()
	u.Write([]byte(str))
	return hex.EncodeToString(u.Sum(nil))
}

// 获取36位长度的apnsId
func getApnsId() string {
	apns := md5Str(strconv.FormatInt(time.Now().Unix(), 10))
	return apns[:8] + "-" + apns[8:12] + "-" + apns[12:16] + "-" + apns[16:20] + "-" + apns[20:]
}
