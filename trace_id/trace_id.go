package trace_id

import (
	"encoding/hex"
	"fmt"
	"github.com/diversability/gocom/goroutineid"
	"net"
	"strings"
	"sync"
	"sync/atomic"
	"time"
)

const TraceIDName = "traceId"

var mGoId2TraceId sync.Map
var ipHexString string

func init() {
	goroutineid.HandleWhenExit(func(goroutineId int64) {
		mGoId2TraceId.Delete(goroutineId)
	})

	iFaces, err := net.Interfaces()
	if err != nil {
		fmt.Println("Interfaces err: ", err.Error())
		return
	}

	for _, iFace := range iFaces {
		if iFace.Flags&net.FlagUp == 0 {
			continue // interface down
		}

		if iFace.Flags&net.FlagLoopback != 0 {
			continue // loopback interface
		}

		addrs, err := iFace.Addrs()
		if err != nil {
			fmt.Println("Addrs err: ", err.Error())
			return
		}

		for _, addr := range addrs {
			ip := getIpFromAddr(addr)
			if ip == nil {
				continue
			}

			ipHexString = hex.EncodeToString(ip)
			return
		}
	}
}

func getIpFromAddr(addr net.Addr) net.IP {
	var ip net.IP
	switch v := addr.(type) {
	case *net.IPNet:
		ip = v.IP
	case *net.IPAddr:
		ip = v.IP
	}
	if ip == nil || ip.IsLoopback() {
		return nil
	}
	ip = ip.To4()
	if ip == nil {
		return nil // not an ipv4 address
	}

	return ip
}

var gTime time.Time
var gCount int32

func SaveTraceId(traceId string) {
	if traceId == "" {
		now := time.Now()
		var count int32
		if now.Equal(gTime) {
			count = atomic.AddInt32(&gCount, 1)
		} else {
			gTime = now
		}

		traceId = fmt.Sprintf("%s%s%d", ipHexString, gTime.Format("0102150405.999"), count)
		traceId = strings.ReplaceAll(traceId, ".", "")
	}

	mGoId2TraceId.Store(goroutineid.GetGoID(), traceId)
}

func GetTraceId() string {
	value, ok := mGoId2TraceId.Load(goroutineid.GetGoID())
	if ok {
		return value.(string)
	} else {
		return ""
	}
}
