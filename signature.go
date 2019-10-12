package main

import (
	"crypto/sha256"
	"fmt"
	"github.com/kataras/iris"
	"math"
	"sort"
	"strconv"
	"strings"
	"time"
)

// 从Iris的上下文里获取并校验签名信息
func CheckSignaturesFromIrisContext(ctx iris.Context) bool {
	nonce := ctx.URLParamDefault("nonce", "")
	timeStamp := ctx.URLParamDefault("timestamp", "")
	sign := ctx.URLParamDefault("signature", "")

	timeStampNumber, err := strconv.ParseInt(timeStamp, 10, 64)
	nowTime := time.Now().Unix()
	if err != nil || math.Abs(float64(timeStampNumber - nowTime)) > 600 {
		// 如果时间戳解析失败或者超过正负600秒，则认为失效
		_, _ = ctx.JSON(RESTfulAPIResult{
			Status: false,
			ErrCode: 40301,
			Message: "Permission Denied: Invalid signature.",
		})
		return false
	}

	success := CheckSingatures(sign, RESTfulSignatureParams{
		Nonce: nonce,
		TimeStamp: timeStamp,
	})
	if !success {
		_, _ = ctx.JSON(RESTfulAPIResult{
			Status: false,
			ErrCode: 40301,
			Message: "Permission Denied: Invalid signature.",
		})
	}
	return success
}

// 进行签名运算
func GenerateSignatures(params *RESTfulSignatureParams) {
	token := PortalToken
	s256 := sha256.New()
	msg := []string { token, params.Nonce, params.TimeStamp }
	sort.Strings(msg)
	msgJoined := []byte(strings.Join(msg, ""))
	s256.Write(msgJoined)
	hex := s256.Sum(nil)
	params.Signature = fmt.Sprintf("%x", hex)
	return
}

// 检查签名
func CheckSingatures(signature string, params RESTfulSignatureParams) bool {
	GenerateSignatures(&params)
	return strings.Trim(signature, " ") != "" && signature == params.Signature
}
