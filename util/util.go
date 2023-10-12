/**
 * 格式化方法
**/

package util

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/http"
)

// 返回格式
type Response struct {
	Code    int
	Message string
}

var (
	Responses Response
)

func AssembleResult(w http.ResponseWriter, code int, message string) {
	Responses.Code = code
	Responses.Message = message
	res, err := json.Marshal(Responses)
	if err == nil {
		fmt.Fprintf(w, string(res))
	}
}

func MD5Encode(str string) string {
	md5Ctx := md5.New()
	md5Ctx.Write([]byte(str))
	encodeStr := md5Ctx.Sum(nil)

	return hex.EncodeToString(encodeStr)
}

func SubStr(str string, start int, strlen int) string {
	end := start + strlen
	s := string([]byte(str)[start:end])
	return s
}
