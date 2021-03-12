package gocom

import (
	"bytes"
	"encoding/json"
	"github.com/diversability/gocom/log"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
)

// 有reqform时，可能需要设置header：request.Header.Set("Content-Type", "application/x-www-form-urlencoded")
func FormatReqUrl(reqUrl string, reqForm map[string]string) string {
	if reqForm == nil || len(reqForm) == 0 {
		return reqUrl
	}

	form := url.Values{}
	return reqUrl + "?" + form.Encode()
}

// method POST GET
// headers Ext Header
func DoHttpRequest(method string, reqUrl string, headers map[string]string, body io.Reader) (int, []byte, error) {
	return _doHttpRequest(method, reqUrl, headers, body)
}

func DoHttpRequestWithBody(method string, reqUrl string, headers map[string]string, reqBody interface{}) (int, []byte, error) {
	if reqBody == nil {
		return _doHttpRequest(method, reqUrl, headers, nil)
	}

	if headers == nil {
		headers = make(map[string]string)
	}

	headers["Content-Type"] = "application/json"
	b, _ := json.Marshal(reqBody)
	reader := strings.NewReader(string(b))
	return _doHttpRequest(method, reqUrl, headers, reader)
}

func _doHttpRequest(method string, reqUrl string, headers map[string]string, body io.Reader) (int, []byte, error) {
	req, _ := http.NewRequest(method, reqUrl, body)

	if headers != nil && len(headers) > 0 {
		for key, value := range headers {
			req.Header.Set(key, value)
		}
	}

	client := http.Client{}
	response, err := client.Do(req)
	if nil != err {
		log.ErrorF("send request err: %v", err)
		return http.StatusNotFound, nil, err
	}

	defer response.Body.Close()

	rspBody, err := ioutil.ReadAll(response.Body)
	if nil != err {
		rspBody = make([]byte, 0)
	}

	return response.StatusCode, rspBody, err
}

func HttpForward(w http.ResponseWriter, r *http.Request, forwardUrl string) error {
	log.DebugF("HttpForward url: %+v", r.URL)

	cli := &http.Client{}
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.ErrorF("HttpForward read w err: %s", err.Error())
		return err
	}

	req, err := http.NewRequest(r.Method, forwardUrl, strings.NewReader(string(body)))
	if err != nil {
		log.ErrorF("HttpForward http.NewRequest err: %s", err.Error())
		return err
	}

	for k, v := range r.Header {
		req.Header.Set(k, v[0])
	}

	res, err := cli.Do(req)
	if err != nil {
		log.ErrorF("HttpForward Do Request err: %s", err.Error())
		return err
	}

	defer res.Body.Close()
	for k, v := range res.Header {
		w.Header().Set(k, v[0])
	}

	io.Copy(w, res.Body)
	return nil
}

func HttpPostJson(addr string, urlParams url.Values, headers map[string]string, body interface{}) ([]byte, int, error) {
	var request *http.Request
	var err error = nil

	if len(urlParams) > 0 {
		addr += urlParams.Encode()
	}

	var bodyReader io.Reader
	if body != nil {
		req, err := json.Marshal(body)
		if err != nil {
			return nil, 0, err
		}

		bodyReader = bytes.NewReader(req)
	}

	request, err = http.NewRequest("POST", addr, bodyReader)
	if err != nil {
		return nil, 0, err
	}

	for k, v := range headers {
		request.Header.Set(k, v)
	}

	request.Header.Set("Content-Type", "application/json")

	cli := &http.Client{}
	response, err := cli.Do(request)
	if nil != err {
		log.ErrorF("httpRequest: Do request (%+v) error:%v", request, err)
		return nil, 0, err
	}

	if response != nil && response.Body != nil {
		defer response.Body.Close()
	}

	data, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.ErrorF("httpRequest: read response error:%v", err)
		return nil, 0, err
	}

	return data, response.StatusCode, nil
}

func HttpFormRequest(method, addr string, urlParams url.Values, headers map[string]string) ([]byte, int, error) {
	var request *http.Request
	var err error = nil
	if method == "GET" || method == "DELETE" {
		if len(urlParams) > 0 {
			addr = addr + "?" + urlParams.Encode()
		}

		request, err = http.NewRequest(method, addr, nil)
		if err != nil {
			return nil, 0, err
		}
	} else {
		request, err = http.NewRequest(method, addr, strings.NewReader(urlParams.Encode()))
		if err != nil {
			return nil, 0, err
		}
		request.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	}

	for k, v := range headers {
		request.Header.Set(k, v)
	}

	cli := &http.Client{}
	response, err := cli.Do(request)
	if nil != err {
		log.ErrorF("httpRequest: Do request (%+v) error:%v", request, err)
		return nil, 0, err
	}

	if response != nil && response.Body != nil {
		defer response.Body.Close()
	}

	data, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.ErrorF("httpRequest: read response error:%v", err)
		return nil, 0, err
	}

	return data, response.StatusCode, nil
}
