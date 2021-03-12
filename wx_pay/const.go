package wx_pay

import (
	"encoding/xml"
	"strconv"
)

//地址的定义
const (
	//WX_APP_ID            = "wxb7107d9243749e29" //小程序id
	//WX_MCH_ID            = "1543115441"
	//WX_APP_KEY           = "bt7U56KBDNgXzaVLAOpTJE2eOziTrz2T"
	WX_TRANSFERS_URL     = "https://api.mch.weixin.qq.com/mmpaymkttransfers/promotion/transfers"
	WX_TRANSFERSinfo_URL = "https://api.mch.weixin.qq.com/mmpaymkttransfers/gettransferinfo"
	WX_UNIFIED_ORDER     = "https://api.mch.weixin.qq.com/pay/unifiedorder"
	WX_REFUND_URL        = "https://api.mch.weixin.qq.com/secapi/pay/refund"
	WX_NotifyURL         = "https://api.yjlsj.cn/v1/wx/pay_callback"
)

//
//const WX_KEY = `-----BEGIN PRIVATE KEY-----
//MIIEvQIBADANBgkqhkiG9w0BAQEFAASCBKcwggSjAgEAAoIBAQC8HKro5FSdPg7J
//AAXqXGXnrLb5LIrTR51joFvonDlXuW4zK6BkfB1a8jbLvkElfPZAbA8qSsJ7dFun
//A3J33Y2of/D7jUVK70ih8sAtrN84JKOeKvro182kJiUNhA+PQN9b/YKvh0CQSjSz
//H4MYIbHl0R6ypGByUVcJt3xXOjRH/wqBEsJzl6oEBnJ3OjMUA0vPmkLVocUXtcHK
//AqjujPEhWAQ3YW8mXPKyGobuOftzK9WCrysgPxd+WcV3CPhbIbWNwqBd9eU2w8C2
//eb+Ps7f5JMWz29mxd930R+VernEc19mvlsXKMk5y/8N6fBX5ppuIHDyet04FENnW
//5TOEYxUtAgMBAAECggEAGyLPjNUTV7OSEnDMaah2ktsZcgx44k2caLjDSWTv6LW5
//LeyHMLeuzGXQfceuQigqpdRww5sRPxnj9s3Kf3wYaUw7iS4x5sNp6OLJ0kzzzneK
//mtB8bYZkBd/yzGZWkEW9ctm5NnT+XVI3E/fhw2No9Ewcb4zC1Pri4WX0q+ibjh7/
//A32Z5/o6lXszEVFc4p2aC9eCOfiDhhk56A+na/koE4dEqntn3rUG3aq5G+E8rrNo
//T+Dqvo85nPdZXt7Ox8Sn3d2wcCCrR/Q/kNNYlAsu1EKac5BryV0Tx5ZDhbgpFcXD
//oh8zbAmGSuCCL/S1hEOLjIz1x8H/YeLS2/T4iApbCQKBgQDqaFkDVcPed+QpwUKs
//r/FLCiu2eXqOqLWOabBr8PV+cbqz7vGFZ9wx4iUABJ9WUeJmAQsmENJBh6MX/n6P
//GbcGh64UK/vHbJirrqYb85QnSJroSAwfQKYbN3Dpi5qckz46dv0k44yMi3kcoJxe
//Lkn8YHg5/zF5VxWpeqK0Fet3LwKBgQDNcJsOTi+Tctk2HvcxzF0L8U+fb/Rfd87d
//EcBadkq0ywkihu6vXCEtfJGIQAd8PkImr0MB/MvErraBiyp0V5ULWOX247CyDfxN
//tC8EnxvAJWy1jpxZ1Jt25mCDKcqWlW/9yEGHCNlcL83JAvRrkbXE5lBGTMblhoyi
//ZNI3l4xiYwKBgECGKLp6SUhbyDqWMDxI0irNyeqY1dufJRrmjOGpmmoL9FDDXUhT
//ppE0puqyWwnv0Fozv1XjG31eUM6yBzRs56ysfIag9NWYVw4rLR5UlluZ6Mo3yt5v
//dUnYoQQooY6oGWEOj/Avkui9G8F9lI14QHVwOKf+TygPiK72SwM3ZXGRAoGAXEOb
//T5Rrp4Pn63eCuxm2HBv3D3rfPFT5Ua2cPsRrjsC0zI3e+mCdAem1DoT7F6B6YxdJ
//N8ZJ2X5BtvJCUdfXty3osbXWcFD5pAgtKZ0vgF8OcIeozms+muqiI6YMNw7MKiTa
//0QN3YwCRIhqynPDmupRZLwliNkj0NiajhpYIVVMCgYEAieY0u/S4osIQ3LC5JVfo
//PUetdEHKg1OutbUJPkYoSQnflIoCmYgfTc8xVcZOtxKRj3rrzVX93xapySYwrQXc
//GAdVc6QiVk2dUXqdzYDuCZENm6ITTQX7omHT0V6DCpXgba6Jtm+XVhj9d5rWJrY7
//y9gVPyFDhvwe/HhO/JHb49o=
//-----END PRIVATE KEY-----`
//
//const WX_CERT = `-----BEGIN CERTIFICATE-----
//MIID7DCCAtSgAwIBAgIUWHU8rOaAuer2WhFfxBenWWmPZJ0wDQYJKoZIhvcNAQEL
//BQAwXjELMAkGA1UEBhMCQ04xEzARBgNVBAoTClRlbnBheS5jb20xHTAbBgNVBAsT
//FFRlbnBheS5jb20gQ0EgQ2VudGVyMRswGQYDVQQDExJUZW5wYXkuY29tIFJvb3Qg
//Q0EwHhcNMTkwODA5MDkxODA1WhcNMjQwODA3MDkxODA1WjB+MRMwEQYDVQQDDAox
//NTQzMTE1NDQxMRswGQYDVQQKDBLlvq7kv6HllYbmiLfns7vnu58xKjAoBgNVBAsM
//IeaIkOmDveaYn+i+sOaxh+enkeaKgOaciemZkOWFrOWPuDELMAkGA1UEBgwCQ04x
//ETAPBgNVBAcMCFNoZW5aaGVuMIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKC
//AQEAvByq6ORUnT4OyQAF6lxl56y2+SyK00edY6Bb6Jw5V7luMyugZHwdWvI2y75B
//JXz2QGwPKkrCe3RbpwNyd92NqH/w+41FSu9IofLALazfOCSjnir66NfNpCYlDYQP
//j0DfW/2Cr4dAkEo0sx+DGCGx5dEesqRgclFXCbd8Vzo0R/8KgRLCc5eqBAZydzoz
//FANLz5pC1aHFF7XBygKo7ozxIVgEN2FvJlzyshqG7jn7cyvVgq8rID8XflnFdwj4
//WyG1jcKgXfXlNsPAtnm/j7O3+STFs9vZsXfd9EflXq5xHNfZr5bFyjJOcv/DenwV
//+aabiBw8nrdOBRDZ1uUzhGMVLQIDAQABo4GBMH8wCQYDVR0TBAIwADALBgNVHQ8E
//BAMCBPAwZQYDVR0fBF4wXDBaoFigVoZUaHR0cDovL2V2Y2EuaXRydXMuY29tLmNu
//L3B1YmxpYy9pdHJ1c2NybD9DQT0xQkQ0MjIwRTUwREJDMDRCMDZBRDM5NzU0OTg0
//NkMwMUMzRThFQkQyMA0GCSqGSIb3DQEBCwUAA4IBAQArLnHiNl03zDDbJ11d5g5S
//JXReUoiagrqcuNcXYJ22un2z8amtcCaMGrjLZWDWHbznyKTK7ryaCqMTzjXyVLyW
//ayBJs++oTYZjlEde8vNw+rBsR6w3kvARIUMjezgs26JV3TESV2GQE+uaNEjeo1pv
//vmOz5NXOsxst+sFW9WPKw5upODi+fsZ2JGq+d2wyaFDjGIyyYh8Wdu1LlCiFvcKF
//tF4YXTLQPt1TR+/ruMmegXBoOTHdKOt6rdxtb7WgvA2QgKaPNoFwzEGPpAJZk/As
//YIJKwldPuuLZCGcfeCfukU9/rBR8H0HdFhw//UzSuG4muDICzSm9wV4rk3f1LXOz
//-----END CERTIFICATE-----`

const bodyType = "application/xml; charset=utf-8"

type WXPayParams map[string]string

func (p WXPayParams) SetString(k, s string) {
	p[k] = s
}

func (p WXPayParams) GetString(k string) string {
	s, _ := p[k]
	return s
}

func (p WXPayParams) SetInt64(k string, i int64) {
	p[k] = strconv.FormatInt(i, 10)
}

func (p WXPayParams) GetInt64(k string) int64 {
	i, _ := strconv.ParseInt(p.GetString(k), 10, 64)
	return i
}

type WeiXinTransferPayRsp struct {
	XMLName        xml.Name `xml:"xml"`
	ReturnCode     string   `xml:"return_code"`
	ReturnMsg      string   `xml:"return_msg"`
	ResultCode     string   `xml:"result_code"`
	MchAppid       string   `xml:"mch_appid"`
	Mchid          string   `xml:"mchid"`
	NonceStr       string   `xml:"nonce_str"`
	PartnerTradeNo string   `xml:"partner_trade_no"`
	PaymentNo      string   `xml:"payment_no"`
	PaymentTime    string   `xml:"payment_time"`
	ErrCode        string   `xml:"err_code"`
	ErrCodeDes     string   `xml:"err_code_des"`
}

type GetTransferInfoRsp struct {
	XMLName        xml.Name `xml:"xml"`
	ReturnCode     string   `xml:"return_code"`
	ReturnMsg      string   `xml:"return_msg"`
	ResultCode     string   `xml:"result_code"`
	MchID          string   `xml:"mch_id"`
	AppID          string   `xml:"appid"`
	DetailID       string   `xml:"detail_id"`
	PartnerTradeNo string   `xml:"partner_trade_no"`
	Status         string   `xml:"status"`
	PaymentAmount  string   `xml:"payment_amount"`
	Openid         string   `xml:"openid"`
	TransferTime   string   `xml:"transfer_time"`
	TransferName   string   `xml:"transfer_name"`
	Desc           string   `xml:"desc"`
	Reason         string   `xml:"reason"`
}

