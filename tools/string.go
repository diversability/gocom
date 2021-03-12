package tools

// 数组到字符串： strings.Replace(strings.Trim(fmt.Sprint(ts.RoomIds), "[]"), " ", ",", -1)

import (
	crand "crypto/rand"
	"encoding/base64"
	"fmt"
	"github.com/diversability/gocom/log"
	"math"
	mrand "math/rand"
	"strconv"
	"time"
)

var GRand = mrand.New(mrand.NewSource(time.Now().Unix()))

// gen random string. 问题：数字部分过多
func RandStr(length int) string {
	if length == 0 {
		return ""
	}

	newLen := math.Ceil(float64(length) / 2)
	buf := make([]byte, int(newLen))
	_, err := crand.Read(buf)
	if err != nil {
		fmt.Printf("gen rand str err: %s\n", err.Error())
		return ""
	}

	out := fmt.Sprintf("%x", buf)
	return out[:int(length)]
}

const LetterArray = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"
const LetterArrayLen = len(LetterArray)

func GenRandomStr(length int) string {
	bytes := make([]byte, length)
	for i := 0; i < length; i++ {
		b := GRand.Intn(LetterArrayLen)
		bytes[i] = LetterArray[b]
	}

	return string(bytes)
}

const LetterArrayCode = "ABCDEFGHJKMNPQRSTWXYZabcdefhijkmnprstwxyz2345678"
const LetterArrayCodeLen = len(LetterArrayCode)

func GenRandomCodeStr(length int) string {
	bytes := make([]byte, length)
	for i := 0; i < length; i++ {
		b := GRand.Intn(LetterArrayCodeLen)
		bytes[i] = LetterArrayCode[b]
	}

	return string(bytes)
}

func Capitalize(str string) string {
	if len(str) <= 0 {
		return str
	}

	if str[0] >= 97 && str[0] <= 122 {
		return string(str[0]-32) + str[1:]
	} else {
		return str
	}
}

func Base64Encode(input string) string {
	inputBytes := []byte(input)
	buf := make([]byte, base64.RawURLEncoding.EncodedLen(len(inputBytes)))
	base64.RawURLEncoding.Encode(buf, inputBytes)

	return string(buf)
}

func Base64Decode(input string) string {
	inputBytes := []byte(input)
	buf := make([]byte, base64.RawURLEncoding.DecodedLen(len(inputBytes)))
	n, err := base64.RawURLEncoding.Decode(buf, inputBytes)
	if n == 0 {
		return ""
	}

	if err != nil {
		log.ErrorF("ase64.URLEncoding.Decode err: %s for: %s", err.Error(), input)
		return ""
	}

	return string(buf)
}

func Int64ToStr(i int64) string {
	return strconv.FormatInt(i,10)
}

func StrToInt64(str string) (int64, error) {
	return strconv.ParseInt(str, 10, 64)
}
