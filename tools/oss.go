package tools

import (
	"crypto/hmac"
	"crypto/sha1"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/aliyun/aliyun-oss-go-sdk/oss"
	"github.com/diversability/gocom/log"
	"hash"
	"io"
	"net/http"
	"time"
)

type OssBucket struct {
	bucket          *oss.Bucket
	endPoint        string
	accessKeyID     string
	accessKeySecret string
	bucketName      string
	clientCache     string
	callbackUrl     string // 客户端上传成功后，oss的回调地址
	uploadDir       string // 客户端上传文件的目录
	tokenExpireTime int64  // 客户端上时，临时token的过期时间
}

// endPoint, accessKeyID, accessKeySecret, bucketName string, clientCacheTime 为基础配置
// callbackUrl uploadDir tokenExpireTime为客户端上传时的配置，callbackUrl为上传到oss的回调地址，uploadDir为上传的目录，tokenExpireTime为临时token的过期时间
func InitOssBucket(endPoint, accessKeyID, accessKeySecret, bucketName string, clientCacheTime int /*单位秒*/, callbackUrl, uploadDir string, tokenExpireTime int64) (*OssBucket, error) {
	client, err := oss.New(endPoint, accessKeyID, accessKeySecret)
	if err != nil {
		log.ErrorF("init oss client error: %s", err.Error())
		return nil, err
	}

	bucket, err := client.Bucket(bucketName)
	if err != nil {
		log.ErrorF("init oss bucket error: %s", err.Error())
		return nil, err
	}

	return &OssBucket{bucket: bucket, endPoint: endPoint, accessKeyID: accessKeyID, accessKeySecret: accessKeySecret,
		bucketName: bucketName, clientCache: fmt.Sprintf("max-age=%d", clientCacheTime), callbackUrl: callbackUrl, uploadDir: uploadDir, tokenExpireTime: tokenExpireTime}, nil
}

// 服务器端直传
func (bucket *OssBucket) UploadToOss(resourcePath string, contentType string, reader io.Reader) bool {
	options := []oss.Option{
		oss.ContentType(contentType),
		oss.CacheControl(bucket.clientCache), /*缓存365天*/
	}

	signedURL, err := bucket.bucket.SignURL(resourcePath, oss.HTTPPut, 60, options...)
	if err != nil {
		if err != nil {
			log.ErrorF("init oss sign url error: %s", err.Error())
			return false
		}
	}

	err = bucket.bucket.PutObjectWithURL(signedURL, reader, options...)
	if err != nil {
		log.ErrorF("upload res err: %s", err.Error())
		return false
	}

	return true
}

func (bucket *OssBucket) DeleteOssRes(resourcePath string) bool {
	err := bucket.bucket.DeleteObject(resourcePath)
	if err != nil {
		log.ErrorF("delete res err: %s", err.Error())
		return false
	}

	return true
}

// 获取客户端直传签名
// 客户端直传文档：https://help.aliyun.com/document_detail/31925.html， go demo：https://help.aliyun.com/document_detail/91818.html?spm=a2c4g.11186623.2.18.6ff36e28eGmN06#concept-mhj-zzt-2fb
func (bucket *OssBucket) GetPolicyToken(tokenExpireTime int64) string {
	now := time.Now().Unix()
	if tokenExpireTime == 0 {
		tokenExpireTime = bucket.tokenExpireTime
	}

	expireEnd := now + tokenExpireTime
	var tokenExpire = getGmtIso8601(expireEnd)

	//create post policy json
	var config ConfigStruct
	config.Expiration = tokenExpire
	var condition []string
	condition = append(condition, "starts-with")
	condition = append(condition, "$key")
	condition = append(condition, bucket.uploadDir)
	config.Conditions = append(config.Conditions, condition)

	//calculate signature
	result, err := json.Marshal(config)
	deByte := base64.StdEncoding.EncodeToString(result)
	h := hmac.New(func() hash.Hash { return sha1.New() }, []byte(bucket.accessKeySecret))
	io.WriteString(h, deByte)
	signedStr := base64.StdEncoding.EncodeToString(h.Sum(nil))

	var callbackParam CallbackParam
	callbackParam.CallbackUrl = bucket.callbackUrl
	callbackParam.CallbackBody = "filename=${object}&size=${size}&mimeType=${mimeType}&height=${imageInfo.height}&width=${imageInfo.width}"
	callbackParam.CallbackBodyType = "application/x-www-form-urlencoded"
	callbackStr, err := json.Marshal(callbackParam)
	if err != nil {
		fmt.Println("callback json err:", err)
	}

	callbackBase64 := base64.StdEncoding.EncodeToString(callbackStr)

	var policyToken PolicyToken
	policyToken.AccessKeyId = bucket.accessKeyID
	policyToken.Host = bucket.endPoint
	policyToken.Expire = bucket.tokenExpireTime
	policyToken.Signature = signedStr
	policyToken.Directory = bucket.uploadDir
	policyToken.Policy = deByte
	policyToken.Callback = callbackBase64
	response, err := json.Marshal(policyToken)
	if err != nil {
		fmt.Println("json err:", err)
	}

	return string(response)
}

type OssOption struct {
	ContentType string
	ClientCache string /*max-age=300*/
	Origin      string
}

func (bucket *OssBucket) SetOption(resourcePath string, option *OssOption) bool {
	var options []oss.Option
	if option.Origin != "" {
		options = append(options, oss.Origin(option.Origin))
	}
	if option.ContentType != "" {
		options = append(options, oss.ContentType(option.ContentType))
	}
	if option.ClientCache != "" {
		options = append(options, oss.CacheControl(option.ClientCache))
	}

	options = append(options, oss.ACReqMethod("POST"))

	err := bucket.bucket.SetObjectMeta(resourcePath, options...)
	if err != nil {
		log.ErrorF("OptionsMethod res err: %s", err.Error())
		return false
	}

	return true
}

// 客户端直传后，OSS的回调函数中调用该函数进行参数验证。这里都把返回值写入了http.ResponseWriter中
func VerifyCallback(w http.ResponseWriter, r *http.Request) bool {
	log.InfoF("Handle Oss VerifyCallback Request ... ")

	// Get PublicKey bytes
	bytePublicKey, err := getPublicKey(r)
	if err != nil {
		responseOSSFailed(w)
		return false
	}

	// Get Authorization bytes : decode from Base64String
	byteAuthorization, err := getAuthorization(r)
	if err != nil {
		responseOSSFailed(w)
		return false
	}

	// Get MD5 bytes from Newly Constructed Authorization String.
	byteMD5, err := getMD5FromNewAuthString(r)
	if err != nil {
		responseOSSFailed(w)
		return false
	}

	// verifySignature and response to client
	if verifySignature(bytePublicKey, byteMD5, byteAuthorization) {
		// do something you want according to callback_body ...
		responseOSSSuccess(w) // response OK : 200
		return true
	} else {
		responseOSSFailed(w) // response FAILED : 400
		return false
	}
}
