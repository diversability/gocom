package tools

import (
	"github.com/diversability/gocom/log"
	"io/ioutil"
)

func LoadFile(fileFullPath string) []byte {
	if fileFullPath == "" {
		return nil
	}

	f, err := ioutil.ReadFile(fileFullPath)
	if err != nil {
		log.ErrorF("ReadFile err: %s", err.Error())
		return nil
	}

	return f
}
