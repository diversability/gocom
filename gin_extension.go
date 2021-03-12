package gocom

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/diversability/gocom/log"
	"github.com/diversability/gocom/trace_id"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
	"time"
)

var AuthAes *AesCbcPKCS7

func InitGinAuth(key string) error {
	if len(key) != 32 {
		return fmt.Errorf("key must by 32bytes length")
	}

	AuthAes = NewAesCbcPKCS7(key)
	return nil
}

func isDebugLog() bool {
	if log.GLog != nil {
		if log.GLog.LogLevel == log.LogLevelDebug {
			return true
		}
	}

	if log.GSizeLog != nil {
		if log.GSizeLog.LogLevel == log.LogLevelDebug {
			return true
		}
	}

	if log.GDailyLog != nil {
		if log.GDailyLog.LogLevel == log.LogLevelDebug {
			return true
		}
	}

	return false
}

func GinLogger(threshold time.Duration) gin.HandlerFunc {
	return func(c *gin.Context) {
		trace_id.SaveTraceId(c.GetHeader(trace_id.TraceIDName))

		if isDebugLog() {
			if c.Request.Method == http.MethodGet {
				log.DebugF("[GIN DEBUG] %s %s URL: %s Header: %+v", c.Request.Method, c.Request.Proto,
					c.Request.URL.String(), c.Request.Header)
			} else {
				contentType := c.ContentType()
				if contentType == gin.MIMEJSON || contentType == gin.MIMEHTML || contentType == gin.MIMEXML ||
					contentType == gin.MIMEXML2 || contentType == gin.MIMEPlain || contentType == gin.MIMEPOSTForm ||
					contentType == gin.MIMEMultipartPOSTForm {
					body, _ := ioutil.ReadAll(c.Request.Body)
					c.Request.Body = ioutil.NopCloser(bytes.NewReader(body))

					if body != nil {
						if contentType == gin.MIMEMultipartPOSTForm || c.Request.ContentLength > 512 {
							log.DebugF("[GIN DEBUG] %s %s URL: %s Header: %+v BodyLen: %d", c.Request.Method, c.Request.Proto,
								c.Request.URL.String(), c.Request.Header, c.Request.ContentLength)
						} else {
							log.DebugF("[GIN DEBUG] %s %s URL: %s Header: %+v Body: %s", c.Request.Method, c.Request.Proto,
								c.Request.URL.String(), c.Request.Header, string(body))
						}
					} else {
						log.DebugF("[GIN DEBUG] %s %s URL: %s Header: %+v. Body err", c.Request.Method, c.Request.Proto,
							c.Request.URL.String(), c.Request.Header)
					}
				} else {
					log.DebugF("[GIN DEBUG] %s %s URL: %s Header: %+v", c.Request.Method, c.Request.Proto,
						c.Request.URL.String(), c.Request.Header)
				}
			}
		}

		// Start timer
		start := time.Now()

		// Process request
		c.Next()

		// Stop timer
		end := time.Now()
		latency := end.Sub(start)

		clientIP := c.ClientIP()
		method := c.Request.Method
		statusCode := c.Writer.Status()
		statusColor := log.ColorForStatus(statusCode)
		methodColor := log.ColorForMethod(method)
		userId := c.Request.Header.Get("selfUserId")

		requestData := getRequestData(c)
		log.InfoF("[GIN] %s%s%s %s%s %s%d%s %.03f [%s] [user_id:%s] %s",
			methodColor, method, log.Reset,
			c.Request.Host, requestData,
			statusColor, statusCode, log.Reset,
			latency.Seconds(),
			clientIP,
			userId,
			c.Errors.String())

		if latency > threshold {
			log.WarnF("[GIN SLOW] %s%s%s %s%s %s%d%s %.03f [%s] [user_id:%s] startAt: %s endAt: %s",
				methodColor, method, log.Reset,
				c.Request.Host, requestData,
				statusColor, statusCode, log.Reset,
				latency.Seconds(),
				clientIP,
				userId,
				start.Format("15:04:05.999999999"),
				end.Format("15:04:05.999999999"))
		}
	}
}

func getRequestData(c *gin.Context) string {
	var requestData string
	method := c.Request.Method
	if method == "GET" || method == "DELETE" {
		requestData = c.Request.RequestURI
	} else {
		c.Request.ParseForm()
		requestData = fmt.Sprintf("%s [%s]", c.Request.RequestURI, c.Request.Form.Encode())
	}

	if len(requestData) > 1024 {
		return requestData[:1024]
	} else {
		return requestData
	}
}

// swagger:model
type FailResponse struct {
	// 错误码。0为正常返回，其它为异常，异常时ErrDesc为异常描述。 固定的异常描述请存储于err_code_def.go文件中
	ErrCode int    `json:"err_code"`
	// 错误信息。可展示在页面上
	ErrDesc string `json:"err_desc"`
	// 用于定位某次调用的ID，客户端应该在错误时显示(ErrCode:TraceId)
	TraceId string `json:"trace_id"`
}

// swagger:model
type SimpleResponse struct {
	// 错误码。这里的值为0，主要是为了返回体中有ErrCode字段
	ErrCode int    `json:"err_code"`
	// 返回的内容
	Content string `json:"content"`
}

// swagger:model
type ResponseBase struct {
	// 错误码。这里的值为0，主要是为了返回体中有ErrCode字段
	ErrCode int    `json:"err_code"`
}

func SendSuccSimpleResponse(c *gin.Context, content string) {
	c.Writer.Header().Set("Content-Type", "application/json")
	resp := SimpleResponse{Content: content}
	c.JSON(http.StatusOK, resp)
}

func SendSuccResponse(c *gin.Context, resp interface{}) {
	c.Writer.Header().Set("Content-Type", "application/json")
	c.JSON(http.StatusOK, resp)
}

func SendFailResponse(c *gin.Context, errCode int, errDesc string) {
	c.Writer.Header().Set("Content-Type", "application/json")
	resp := FailResponse{ErrCode:errCode, ErrDesc: errDesc, TraceId: trace_id.GetTraceId()}
	c.JSON(errCode, resp)
}

type UserAgent struct {
	AppVersion        string `json:"app_version"`
	MobilePlatform    string `json:"mobile_platform"`
	MobileSystem      string `json:"mobile_system"`
	MobileDeviceBrand string `json:"mobile_device_brand"`
}

func CheckAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		if isPathInWhiteList(c.Request.URL.Path) {
			// Process request
			c.Next()
			return
		}

		ua := c.GetHeader("Useragent")
		if ua == "" {
			log.WarnF("No UserAgent in the req: %+v", c.Request.Header)
			SendFailResponse(c, http.StatusUnauthorized, "No UserAgent in the req")
			c.Abort()
			return
		}

		var userAgent UserAgent
		err := json.Unmarshal([]byte(ua), &userAgent)
		if err != nil {
			log.ErrorF("Unmarshal UserAgent Fail: %s %s", ua, err.Error())
			SendFailResponse(c, http.StatusUnauthorized, "Error UserAgent")
			c.Abort()
			return
		}

		authToken := c.GetHeader("Authorization")
		tokenPlaintext, err := AuthAes.Decrypt(authToken)
		if err != nil {
			log.ErrorF("Decrypt Authorization error : %s %s", authToken, err.Error())
			SendFailResponse(c, http.StatusUnauthorized, "Wrong Authorization")
			c.Abort()
			return
		}

		items := strings.Split(tokenPlaintext, "|")
		if len(items) != 5 {
			log.ErrorF("Wrong1 Authorization: %s", tokenPlaintext)
			SendFailResponse(c, http.StatusUnauthorized, "Wrong1 Authorization")
			c.Abort()
			return
		}

		if items[0] != userAgent.AppVersion || items[1] != userAgent.MobileSystem || items[2] != userAgent.MobileDeviceBrand {
			log.ErrorF("Wrong2 Authorization: %+v", userAgent)
			SendFailResponse(c, http.StatusUnauthorized, "Wrong2 Authorization")
			c.Abort()
			return
		}

		// 判断过期时间
		timeStamp, err := strconv.ParseInt(items[4], 10, 64)
		if err != nil {
			log.Error("wrong time")
			c.Abort()
			return
		}

		tokenTime := time.Unix(timeStamp, 0)
		if tokenTime.Add(time.Duration(time.Hour * 6)).Before(time.Now()) {
			SendFailResponse(c, http.StatusUnauthorized, "Wrong Authorization, timeout")
			log.InfoF("token timeout. token time: %d", timeStamp)
			c.Abort()
			return
		}

		userId, err := strconv.ParseInt(items[3], 10, 64)
		if err != nil {
			SendFailResponse(c, http.StatusUnauthorized, "Wrong Authorization, Parse userId err")
			log.Error("token wrong userId")
			c.Abort()
			return
		}

		err = GenAuth(c, userId)
		if err != nil {
			log.ErrorF("gen new auth err: %s", err.Error())
		}

		c.Request.Header.Set("selfUserId", items[3])
	}
}

var mWhitePathMap = map[string]Empty{
	"/favicon.ico":    empty,
	"/debug/pprof/*":  empty,
	"/api/test/*":     empty,
	"/api/callback/*": empty,
	"/api/login/*":    empty,
}

func AddWhiteList(url string) {
	mWhitePathMap[url] = empty
}

func isPathInWhiteList(path string) bool {
	for k, _ := range mWhitePathMap {
		if k[len(k)-1] == '*' {
			// 进行前缀匹配
			if strings.HasPrefix(path, k[0:len(k)-1]) {
				return true
			}
		} else if k == path {
			return true
		}
	}
	_, ok := mWhitePathMap[path]
	return ok
}

func GenAuth(c *gin.Context, userId int64) error {
	ua := c.GetHeader("Useragent")
	if ua == "" {
		return fmt.Errorf("No UserAgent")
	}

	var userAgent UserAgent
	err := json.Unmarshal([]byte(ua), &userAgent)
	if err != nil {
		return fmt.Errorf("Unmarshal UserAgent Fail: %s %s", ua, err.Error())
	}

	plaintext := fmt.Sprintf("%s|%s|%s|%d|%d", userAgent.AppVersion, userAgent.MobileSystem, userAgent.MobileDeviceBrand, userId, time.Now().Unix())
	out, err := AuthAes.Encrypt([]byte(plaintext))
	if err != nil {
		return err
	}

	c.Writer.Header().Set("Authorization", out)
	return nil
}

func ShowAllRawBody(c *gin.Context) {
	data, _ := ioutil.ReadAll(c.Request.Body)
	log.InfoF("Raw Http Body: %s", string(data))
	c.Request.Body = ioutil.NopCloser(bytes.NewReader(data))
}

func Cors() gin.HandlerFunc {
	return func(c *gin.Context) {
		method := c.Request.Method
		origin := c.Request.Header.Get("Origin")
		if origin == "" {
			origin = "*"
		}
		c.Header("Access-Control-Allow-Origin", origin)
		c.Header("Access-Control-Allow-Headers", "Content-Type,AccessToken,X-CSRF-Token, Authorization, Useragent, selfuserid")
		c.Header("Access-Control-Allow-Methods", "POST, GET, OPTIONS")
		c.Header("Access-Control-Expose-Headers", "Content-Length, Access-Control-Allow-Origin,"+
			" Access-Control-Allow-Headers, Content-Type,Authorization")
		c.Header("Access-Control-Allow-Credentials", "true")
		if method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
		}
		c.Next()
	}
}

func Bind(c *gin.Context, obj interface{}) bool {
	err := c.Bind(obj)
	if err != nil {
		// this check is only needed when your code could produce
		// an invalid value for validation such as interface with nil
		// value most including myself do not usually have code like this.
		// validator 参考： https://github.com/go-playground/validator
		if _, ok := err.(*validator.InvalidValidationError); ok {
			log.ErrorF("bind err. InvalidValidationError: %s", err.Error())
			return false
		}

		errs, ok := err.(validator.ValidationErrors)
		if ok {
			for _, err := range errs {
				log.ErrorF("bind err. ValidationError. StructField: %s, Tag: %s %s, Type: %+v, Value: %+v", err.StructNamespace(), err.ActualTag(), err.Param(), err.Type(), err.Value())
			}
		} else {
			log.ErrorF("bind err. Unknown Err: %s", err.Error())
		}

		return false
	} else {
		return true
	}
}

// 如果绑定失败，则返回错误描述
func Bind2(c *gin.Context, obj interface{}) (bool, string) {
	retErr := false
	if log.GSizeLog != nil && (log.GSizeLog.LogLevel == log.LogLevelDebug) {
		retErr = true
	}

	err := c.Bind(obj)
	if err != nil {
		// this check is only needed when your code could produce
		// an invalid value for validation such as interface with nil
		// value most including myself do not usually have code like this.
		// validator 参考： https://github.com/go-playground/validator
		if _, ok := err.(*validator.InvalidValidationError); ok {
			log.ErrorF("bind err. InvalidValidationError: %s", err.Error())
			if retErr {
				return false, err.Error()
			} else {
				return false, ""
			}
		}

		errs, ok := err.(validator.ValidationErrors)
		if ok {
			for _, err := range errs {
				log.ErrorF("bind err. ValidationError. StructField: %s, Tag: %s %s, Type: %+v, Value: %+v", err.StructNamespace(), err.ActualTag(), err.Param(), err.Type(), err.Value())
			}
		} else {
			log.ErrorF("bind err. Unknown Err: %s", err.Error())
		}

		if retErr {
			return false, err.Error()
		} else {
			return false, ""
		}
	}

	return true, ""
}
