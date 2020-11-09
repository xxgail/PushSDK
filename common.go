package PushSDK

import (
	"crypto/md5"
	"crypto/sha256"
	"encoding/hex"
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

func md5Str(str1, str2 string) string {
	u := md5.New()
	u.Write([]byte(str1 + str2))
	return hex.EncodeToString(u.Sum(nil))
}
