package utils

import (
	"encoding/base64"
	"errors"
	"github.com/h2non/filetype"
	"math/rand"
	"time"
	"unsafe"
)

var (
	src           = rand.NewSource(time.Now().UnixNano())
	Constellation = [12]string{
		"水瓶座", // 1月
		"双鱼座", // 2月
		"白羊座", // 3月
		"金牛座", // 4月
		"双子座", // 5月
		"巨蟹座", // 6月
		"狮子座", // 7月
		"处女座", // 8月
		"天秤座", // 9月
		"天蝎座", // 10月
		"射手座", // 11月
		"魔羯座", // 12月
	}
	ErrConstellationNumNotCorrect = errors.New("星座num不正确")
)

const (
	letters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
)
const (
	// 6 bits to represent a letter index
	letterIdBits = 6
	// All 1-bits as many as letterIdBits
	letterIdMask = 1<<letterIdBits - 1
	letterIdMax  = 63 / letterIdBits
)

// RandStr 随机生成长度为n只包含大小写字母的字符串
func RandStr(n int) string {
	b := make([]byte, n)
	// A rand.Int63() generates 63 random bits, enough for letterIdMax letters!
	for i, cache, remain := n-1, src.Int63(), letterIdMax; i >= 0; {
		if remain == 0 {
			cache, remain = src.Int63(), letterIdMax
		}
		if idx := int(cache & letterIdMask); idx < len(letters) {
			b[i] = letters[idx]
			i--
		}
		cache >>= letterIdBits
		remain--
	}
	return *(*string)(unsafe.Pointer(&b))
}

var enc = base64.StdEncoding

// Base64Decode base64Str是base64编码格式中的<data>
// data:[<MIME-type>][;charset=<encoding>][;base64],<data>
func Base64Decode(base64Str string) ([]byte, error) {
	return enc.DecodeString(base64Str)
}

// GetFileType 根据文件获取文件对应的文件类型
func GetFileType(buf []byte) string {
	if filetype.IsImage(buf) == true {
		return "image"
	} else if filetype.IsVideo(buf) == true {
		return "video"
	} else if filetype.IsAudio(buf) == true {
		return "audio"
	} else {
		return "unknown"
	}
}

// GetConstellationCNNameByNum 根据num取出星座对应的中文名
func GetConstellationCNNameByNum(num int) (string, error) {
	if num < 0 || num > 11 {
		return "", ErrConstellationNumNotCorrect
	}
	return Constellation[num], nil
}
