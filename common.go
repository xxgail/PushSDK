package PushSDK

import (
	"crypto/md5"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"math/rand"
	"runtime"
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

func isEmpty(s string) error {
	if s == "" {
		return errors.New("platform param can not be empty")
	}
	var m map[string]string
	_ = json.Unmarshal([]byte(s), &m)
	for k, v := range m {
		if v == "" {
			return errors.New(k + "不能为空")
		}
	}
	return nil
}

func getRandomApnsId() string {
	randStr := strconv.Itoa(rand.Int())
	apns := md5Str(strconv.FormatInt(time.Now().Unix(), 10) + randStr)
	return apns[:8] + "-" + apns[8:12] + "-" + apns[12:16] + "-" + apns[16:20] + "-" + apns[20:]
}

func getFileLineNum() string {
	_, file, line, _ := runtime.Caller(1)
	return time.Now().Format("2006-01-02 15:04:05") + "▶ 我走到这里啦！" + file + "--line" + strconv.Itoa(line)
}
