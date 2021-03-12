package wx_pay

import (
	"bytes"
	"crypto/md5"
	"crypto/tls"
	"encoding/hex"
	"encoding/xml"
	"fmt"
	"github.com/diversability/gocom/log"
	"io"
	"io/ioutil"
	"net/http"
	"sort"
	"strings"
	"time"
)

/*
微信支付

重要接口：
- PlaceAnWXPayOrder 用户下单
- DoWXRefund 申请退款
*/

var WXPayClient *Client // InitWXPayClient(WX_APP_ID, WX_MCH_ID, WX_APP_KEY)

// 实例化API客户端
func InitWXPayClient(appId, mchId, apiKey, wxKey, wxCert string) *Client {
	WXPayClient = &Client{
		stdClient: &http.Client{},
		AppId:     appId,
		MchId:     mchId,
		ApiKey:    apiKey,
		WXKey:     []byte(wxKey),
		WXCert:    []byte(wxCert),
	}

	return WXPayClient
}

// 用户下单
func PlaceAnWXPayOrder(orderNo, orderBody string, totalFee int, clientIP, notifyURL, openID string) (*WXPayParams, error) {
	log.InfoF("PlaceAnWXPayOrder ExtOrderId: %s totalFee: %d, openID: %s orderBody: %s", orderNo, totalFee, totalFee, openID, orderBody)

	// 附着商户证书
	if err := WXPayClient.WithCertBytes(WXPayClient.WXCert, WXPayClient.WXKey); err != nil {
		log.ErrorF("WithCertBytes Err: %s", err.Error())
		return nil, err
	}

	params := make(WXPayParams)
	params.SetString("appid", WXPayClient.AppId)
	params.SetString("mch_id", WXPayClient.MchId)
	params.SetString("nonce_str", genRandomStr()) // 随机字符串
	params.SetString("out_trade_no", orderNo)     // 商户订单号
	params.SetString("body", orderBody)
	params.SetInt64("total_fee", int64(totalFee))
	params.SetString("spbill_create_ip", clientIP)
	params.SetString("notify_url", notifyURL)
	params.SetString("trade_type", "JSAPI")
	params.SetString("openid", openID)

	params.SetString("sign", WXPayClient.Sign(params)) // 签名

	ret, err := WXPayClient.post(WX_UNIFIED_ORDER, params, true)
	if err != nil {
		return nil, err
	}

	if ret.GetString("return_code") != "SUCCESS" {
		return nil, fmt.Errorf(ret.GetString("return_msg"))
	}

	payload := make(WXPayParams)
	payload.SetString("appId", WXPayClient.AppId)
	payload.SetInt64("timeStamp", time.Now().Unix())
	payload.SetString("nonceStr", genRandomStr())
	payload.SetString("package", fmt.Sprintf("prepay_id=%s", ret.GetString("prepay_id")))
	payload.SetString("signType", "MD5")

	// 签名
	payload.SetString("sign", WXPayClient.Sign(payload))
	payload.SetString("prepay_id", ret.GetString("prepay_id"))

	return &payload, nil
}

// 申请退款
func DoWXRefund(transactionID, outRefundNo string, totalFee, refundFee int, refundDesc string) error {
	log.InfoF("DoWXRefund WXOrderId: %s ExtOrderId: %s totalFee: %d, refundFee: %d refundDesc: %s", transactionID, outRefundNo, totalFee, refundFee, refundDesc)

	// 附着商户证书
	if err := WXPayClient.WithCertBytes(WXPayClient.WXCert, WXPayClient.WXKey); err != nil {
		log.ErrorF("WithCertBytes Err: %s", err.Error())
		return err
	}

	params := make(WXPayParams)
	params.SetString("appid", WXPayClient.AppId)
	params.SetString("mch_id", WXPayClient.MchId)
	params.SetString("nonce_str", genRandomStr())     // 随机字符串
	params.SetString("transaction_id", transactionID) // 微信生成的订单号，在支付通知中有返回
	params.SetString("out_refund_no", outRefundNo)
	params.SetInt64("total_fee", int64(totalFee))
	params.SetInt64("refund_fee", int64(refundFee))
	//params.SetString("notify_url", notifyURL)
	params.SetString("refund_desc", refundDesc)
	params.SetString("sign_type", "MD5")

	params.SetString("sign", WXPayClient.Sign(params)) // 签名

	ret, err := WXPayClient.post(WX_REFUND_URL, params, true)
	if err != nil {
		return err
	}

	log.InfoF("DoWXRefund post rsp: %+v", ret)

	if ret.GetString("return_code") != "SUCCESS" {
		log.ErrorF("退款失败. ret: %+v, params: %+v", ret, params)
		return fmt.Errorf(ret.GetString("return_msg"))
	}

	if ret.GetString("result_code") != "SUCCESS" {
		log.ErrorF("退款失败. ret: %+v, params: %+v", ret, params)
		return fmt.Errorf(ret.GetString("err_code_des"))
	}

	return nil
}

/**********************************************************************************************************************/

// API客户端
type Client struct {
	stdClient *http.Client
	tlsClient *http.Client

	AppId  string
	MchId  string
	ApiKey string
	WXKey  []byte
	WXCert []byte
}

// 设置请求超时时间
func (c *Client) SetTimeout(d time.Duration) {
	c.stdClient.Timeout = d
	if c.tlsClient != nil {
		c.tlsClient.Timeout = d
	}
}

// 附着商户证书
func (c *Client) WithCert(certFile, keyFile string) error {
	cert, err := ioutil.ReadFile(certFile)
	if err != nil {
		return err
	}

	key, err := ioutil.ReadFile(keyFile)
	if err != nil {
		return err
	}

	return c.WithCertBytes(cert, key)
}

func (c *Client) WithCertBytes(cert, key []byte) error {
	tlsCert, err := tls.X509KeyPair(cert, key)
	if err != nil {
		return err
	}

	conf := &tls.Config{
		Certificates: []tls.Certificate{tlsCert},
	}

	trans := &http.Transport{
		TLSClientConfig: conf,
	}

	c.tlsClient = &http.Client{
		Transport: trans,
	}
	return nil
}

// 发送请求
func (c *Client) post(url string, params WXPayParams, tls bool) (WXPayParams, error) {
	var httpc *http.Client
	if tls {
		if c.tlsClient == nil {
			return nil, fmt.Errorf("tls client is not initialized")
		}

		httpc = c.tlsClient
	} else {
		httpc = c.stdClient
	}

	resp, err := httpc.Post(url, bodyType, c.Encode(params))
	if err != nil {
		return nil, err
	}

	return DecodeWXPayParamsFromXML(resp.Body), nil
}

// XML解码
func DecodeWXPayParamsFromXML(r io.Reader) WXPayParams {
	var (
		d      *xml.Decoder
		start  *xml.StartElement
		params WXPayParams
	)

	d = xml.NewDecoder(r)
	params = make(WXPayParams)
	for {
		tok, err := d.Token()
		if err != nil {
			break
		}

		switch t := tok.(type) {
		case xml.StartElement:
			start = &t
		case xml.CharData:
			if t = bytes.TrimSpace(t); len(t) > 0 {
				params.SetString(start.Name.Local, string(t))
			}
		}
	}
	return params
}

// XML编码
func (c *Client) Encode(params WXPayParams) io.Reader {
	var buf bytes.Buffer
	buf.WriteString(`<xml>`)
	for k, v := range params {
		buf.WriteString(`<`)
		buf.WriteString(k)
		buf.WriteString(`><![CDATA[`)
		buf.WriteString(v)
		buf.WriteString(`]]></`)
		buf.WriteString(k)
		buf.WriteString(`>`)
	}

	buf.WriteString(`</xml>`)
	return &buf
}

// 验证签名
func (c *Client) CheckSign(params WXPayParams) bool {
	return params.GetString("sign") == c.Sign(params)
}

// 生成签名
func (c *Client) Sign(params WXPayParams) string {
	var keys = make([]string, 0, len(params))
	for k := range params {
		if k != "sign" {
			keys = append(keys, k)
		}
	}

	sort.Strings(keys)

	var buf bytes.Buffer
	for _, k := range keys {
		if len(params.GetString(k)) > 0 {
			buf.WriteString(k)
			buf.WriteString(`=`)
			buf.WriteString(params.GetString(k))
			buf.WriteString(`&`)
		}
	}

	buf.WriteString(`key=`)
	buf.WriteString(c.ApiKey)

	sum := md5.Sum(buf.Bytes())
	str := hex.EncodeToString(sum[:])

	return strings.ToUpper(str)
}
