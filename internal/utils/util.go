package utils

import (
	"encoding/base64"
	"github.com/h2non/filetype"
	"math/rand"
	"time"
	"unsafe"
)

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

const letters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

var src = rand.NewSource(time.Now().UnixNano())

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
