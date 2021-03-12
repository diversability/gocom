package wx_pay

import (
	"bytes"
	"crypto/md5"
	"crypto/tls"
	"encoding/hex"
	"encoding/xml"
	"errors"
	"fmt"
	"github.com/diversability/gocom/log"
	"hash"
	"io"
	"io/ioutil"
	"net"
	"net/http"
	"sort"
	"time"
)

/*
微信转账给个人

链接： https://pay.weixin.qq.com/wiki/doc/api/tools/mch_pay.php?chapter=14_2
	如果有异常去网站查询，规则会变的，规则时间（2017-06-06）
接口调用规则：
  ◆ 给同一个实名用户付款，单笔单日限额2W
  ◆ 不支持给非实名用户打款
  ◆ 一个商户同一日付款总额限额100W
  ◆ 单笔最小金额默认为1元
  ◆ 每个用户每天最多可付款10次，可以在商户平台--API安全进行设置
  ◆ 给同一个用户付款时间间隔不得低于15秒
  注意：当返回错误码为“SYSTEMERROR”时，一定要使用原单号重试，否则可能造成重复支付等资金风险。

接口地址： https://api.mch.weixin.qq.com/mmpaymkttransfers/promotion/transfers

重要接口：
- GetTransferInfo  获取订单状态
- WeiXinTransferPay 企业支付
*/

// 微信企业支付
func WXTransfer(amount int, openid, tradeNo, desc, createIP string) (rsp WeiXinTransferPayRsp, err error) {
	//计算签名
	params := make(map[string]string, 16)
	params["mch_appid"] = WXPayClient.AppId
	params["mchid"] = WXPayClient.MchId
	params["nonce_str"] = genRandomStr()
	params["partner_trade_no"] = tradeNo
	params["openid"] = openid
	params["check_name"] = "NO_CHECK"
	params["amount"] = fmt.Sprintf("%d", 1)
	params["desc"] = desc
	params["spbill_create_ip"] = createIP
	sign := sign(params, WXPayClient.ApiKey, nil)
	params["sign"] = sign

	rs, err := formatMapToXMLStr(params)
	if err != nil {
		return rsp, err
	}

	client, err := getWXMpTlsClient()
	if err != nil {
		return rsp, err
	}

	// 一旦发出请求，如果没有得到明确的失败，都认为是三方不能确认，由脚本重新发起
	rs, err = postTlsUrl(client, WX_TRANSFERS_URL, rs)
	if err != nil {
		return rsp, err
	}

	err = xml.Unmarshal([]byte(rs), &rsp)
	if err != nil {
		return rsp, err
	}
	if rsp.PaymentNo == "" {
		log.ErrorF("PaymentNo is empty. tradeNo: %s, rs: %s", tradeNo, rs)
	}

	return rsp, nil
}

// 获取企业微信支付订单状态
func GetTransferInfo(tradeNO string) {
	//计算签名
	params := make(map[string]string, 16)
	params["appid"] = WXPayClient.AppId
	params["mch_id"] = WXPayClient.MchId
	params["nonce_str"] = genRandomStr()
	params["partner_trade_no"] = tradeNO
	sign := sign(params, WXPayClient.ApiKey, nil)
	params["sign"] = sign

	rs, err := formatMapToXMLStr(params)
	if err != nil {
		log.ErrorF("formatMapToXMLStr err: %s", err.Error())
		return
	}

	client, err := getWXMpTlsClient()
	if err != nil {
		return
	}

	// 一旦发出请求，如果没有得到明确的失败，都认为是三方不能确认，由脚本重新发起
	rs, err = postTlsUrl(client, WX_TRANSFERSinfo_URL, rs)
	if err != nil {
		log.ErrorF("GetTransferInfo err: %s", err.Error())
		return
	} else {
		log.DebugF("GetTransferInfo rsp: %s", rs)
	}

	// 结果解析
	rsp := GetTransferInfoRsp{}
	err = xml.Unmarshal([]byte(rs), &rsp)
	if err != nil {
		return
	}

	if rsp.ResultCode == "SUCCESS" && rsp.ReturnCode == "SUCCESS" {
		if rsp.Status == "SUCCESS" {
			log.DebugF("GetTransferInfo rsp: %s", rsp.Status)
		} else {
			log.ErrorF("GetTransferInfo rsp: %s", rsp.Reason)
		}
	} else {
		log.ErrorF("GetTransferInfo rsp: %s", rsp.ReturnMsg)
	}
}

func postTlsUrl(client *http.Client, url string, data string) (string, error) {
	resp, err := client.Post(url, "text/plain", bytes.NewReader([]byte(data)))
	if err != nil {
		return "", err
	}

	defer resp.Body.Close()
	rtData, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", nil
	}

	str := string(rtData)
	return str, nil
}

func newTLSHttpClient(caCert, caKey []byte) (httpClient *http.Client, err error) {
	cert, err := tls.X509KeyPair(caCert, caKey)
	if err != nil {
		return
	}

	tlsConfig := &tls.Config{
		Certificates: []tls.Certificate{cert},
	}

	httpClient = &http.Client{
		Transport: &http.Transport{
			Proxy: http.ProxyFromEnvironment,
			Dial: (&net.Dialer{
				Timeout:   30 * time.Second,
				KeepAlive: 30 * time.Second,
			}).Dial,
			TLSHandshakeTimeout: 10 * time.Second,
			TLSClientConfig:     tlsConfig,
		},
		Timeout: 60 * time.Second,
	}
	return
}

func sign(parameters map[string]string, apiKey string, fn func() hash.Hash) string {
	ks := make([]string, 0, len(parameters))
	for k := range parameters {
		if k == "sign" || k == "sign_type" || k == "paySign" || parameters[k] == "" {
			continue
		}
		ks = append(ks, k)
	}

	sort.Strings(ks)

	if fn == nil {
		fn = md5.New
	}
	h := fn()
	signature := make([]byte, h.Size()*2)

	for _, k := range ks {
		v := parameters[k]
		if v == "" {
			continue
		}
		h.Write([]byte(k))
		h.Write([]byte{'='})
		h.Write([]byte(v))
		h.Write([]byte{'&'})
	}
	h.Write([]byte("key="))
	h.Write([]byte(apiKey))

	hex.Encode(signature, h.Sum(nil))
	return string(bytes.ToUpper(signature))
}

func getWXMpTlsClient() (httpClient *http.Client, err error) {
	return newTLSHttpClient(WXPayClient.WXCert, WXPayClient.WXKey)
}

func formatMapToXMLStr(m map[string]string) (rs string, err error) {
	buf := bytes.NewBuffer(make([]byte, 0, 1024))
	err = formatMapToXML(buf, m)
	if err == nil {
		rs = buf.String()
	}
	return
}

// FormatMapToXML marshal map[string]string to xmlWriter with xml format, the root node name is xml.
//  NOTE: This function assumes the key of m map[string]string are legitimate xml name string
//  that does not contain the required escape character!
func formatMapToXML(xmlWriter io.Writer, m map[string]string) (err error) {
	if xmlWriter == nil {
		return errors.New("nil xmlWriter")
	}

	if _, err = io.WriteString(xmlWriter, "<xml>"); err != nil {
		return
	}

	for k, v := range m {
		if _, err = io.WriteString(xmlWriter, "<"+k+">"); err != nil {
			return
		}
		if err = xml.EscapeText(xmlWriter, []byte(v)); err != nil {
			return
		}
		if _, err = io.WriteString(xmlWriter, "</"+k+">"); err != nil {
			return
		}
	}

	if _, err = io.WriteString(xmlWriter, "</xml>"); err != nil {
		return
	}
	return
}

func parseXMLStrToMap(str string) (m map[string]string, err error) {
	rd := bytes.NewReader([]byte(str))
	m, err = parseXMLToMap(rd)
	return
}

// ParseXMLToMap parses xml reading from xmlReader and returns the first-level sub-node key-value set,
// if the first-level sub-node contains child nodes, skip it.
func parseXMLToMap(xmlReader io.Reader) (m map[string]string, err error) {
	if xmlReader == nil {
		err = errors.New("nil xmlReader")
		return
	}

	d := xml.NewDecoder(xmlReader)
	m = make(map[string]string)

	var (
		tk    xml.Token
		depth int // xml.Token depth
		key   string
		value bytes.Buffer
	)
	for {
		tk, err = d.Token()
		if err != nil {
			if err != io.EOF {
				return
			}
			err = nil
			return
		}

		switch v := tk.(type) {
		case xml.StartElement:
			depth++
			switch depth {
			case 1:
			case 2:
				key = v.Name.Local
				value.Reset()
			case 3:
				if err = d.Skip(); err != nil {
					return
				}
				depth--
				key = "" // key == "" indicates that the node with depth==2 has children
			default:
				panic("incorrect algorithm")
			}
		case xml.CharData:
			if depth == 2 && key != "" {
				value.Write(v)
			}
		case xml.EndElement:
			if depth == 2 && key != "" {
				m[key] = value.String()
			}
			depth--
		}
	}
}
