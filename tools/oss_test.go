package tools

import (
	"github.com/chanxuehong/util/math"
	"github.com/diversability/gocom/log"
	"testing"
)

// go test -v oss_test.go oss.go oss_helper.go -test.run TestOssSetOption
func TestOssSetOption(t *testing.T) {
	_, err := log.InitLog("./", "test.log", "debug", 0)
	if nil != err {
		t.Fatal("initLog err :", err)
		return
	}

	PdfBucket, err := InitOssBucket()
	if nil != err {
		t.Fatal("InitOssBucket err : ", err)
		return
	}

	option := OssOption{ClientCache: "max-age=2147483647", Origin: "www.doctool.cn"}
	ret := PdfBucket.SetOption("pdf/943e3acc-b8e7-4fe7-896f-fbde70dbd111.pdf", &option)
	if !ret {
		t.Fatal("fail")
	}
}
