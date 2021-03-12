package goroutineid

import (
	"sync"
	"unsafe"
)

// 获取近似的GOID。使用该方法后，不能再hack goExit
type HandlerFunc func(goroutineId int64)

var (
	goIdMu           sync.RWMutex
	goIdDataMap            = map[unsafe.Pointer]int64{}
	maxProximateGoID int64 = 1000000
	handlerFuncList  []HandlerFunc
)

func GetGoID() int64 {
	gp := G()

	if gp == nil {
		return 0
	}

	goIdMu.RLock()
	goId, ok := goIdDataMap[gp]
	goIdMu.RUnlock()

	if ok {
		return goId
	}

	// 新的goruntine
	goIdMu.Lock()
	defer goIdMu.Unlock()

	maxProximateGoID = maxProximateGoID + 1
	goIdDataMap[gp] = maxProximateGoID
	if !hack(gp) {
		delete(goIdDataMap, gp)
	}

	return maxProximateGoID
}

func resetAtExit() {
	gp := G()

	if gp == nil {
		return
	}

	goIdMu.RLock()
	goId, ok := goIdDataMap[gp]
	goIdMu.RUnlock()

	if ok {
		for _, f := range handlerFuncList {
			f(goId)
		}

		goIdMu.Lock()
		delete(goIdDataMap, gp)
		goIdMu.Unlock()
	}
}

func HandleWhenExit(f HandlerFunc) {
	handlerFuncList = append(handlerFuncList, f)
}
