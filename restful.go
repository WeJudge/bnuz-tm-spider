package main

type RESTfulAPIResult struct {
	Status bool `json:"status"`
	ErrCode int `json:"errcode"`
	Message string `json:"message"`
	Data interface{} `json:"data"`
}

type RESTfulSignatureParams struct {
	// 签名
	Signature string `json:"signature"`
	// 时间戳
	TimeStamp string `json:"timestamp"`
	// 随机字符串
	Nonce string `json:"nonce"`
}
