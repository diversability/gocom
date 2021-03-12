package gocom

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"time"
)

type MysqlConfig struct {
	MysqlConn            string
	MysqlConnectPoolSize int
}

type LogConfig struct {
	LogDir   string
	LogFile  string
	LogLevel string
	LogSize  int
}

type RedisConfig struct {
	RedisConn      string
	RedisPassword  string
	ReadTimeout    int
	ConnectTimeout int
	WriteTimeout   int
	IdleTimeout    int
	MinIdleConnNum int
	PoolSize       int
	RedisDb        int
}

type OssConfig struct {
	AccessKeyId     string
	AccessKeySecret string
	Region          string
	Bucket          string
}

type Configure struct {
	Listen        string
	LogSetting    LogConfig //推荐
	MysqlSetting  map[string]MysqlConfig
	RedisSetting  map[string]RedisConfig
	OssSetting    map[string]OssConfig
	External      map[string]string
	ExternalInt64 map[string]int64
	GormDebug     bool   //sql 输出开关
	Environment   string //环境变量区分, ONLINE/TEST
}

var Config *Configure
var g_config_file_last_modify_time time.Time

// example: LoadCfgFromFile("./config.json", cfg)
func LoadCfgFromFile(filename string, config *Configure) error {
	fmt.Println("filename", filename)
	fi, err := os.Stat(filename)
	if err != nil {
		fmt.Println("ReadFile: ", err.Error())
		return err
	}

	if g_config_file_last_modify_time.Equal(fi.ModTime()) {
		return nil
	} else {
		g_config_file_last_modify_time = fi.ModTime()
	}

	bytes, err := ioutil.ReadFile(filename)
	if err != nil {
		fmt.Println("ReadFile: ", err.Error())
		return err
	}

	if err := json.Unmarshal(bytes, config); err != nil {
		err = fmt.Errorf("unmarshal error :%s", string(bytes))
		log.Println(err.Error())
		return err
	}

	fmt.Println("conifg :", *config)
	Config = config
	return nil
}
