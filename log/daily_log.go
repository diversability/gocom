package log

import (
	"fmt"
	"github.com/chanxuehong/util/math"
	"github.com/diversability/gocom/goroutineid"
	"github.com/diversability/gocom/trace_id"
	"log"
	"os"
	"strings"
	"sync"
	"time"
)

// 按天进行日志切割

var GDailyLog *DailyLog

type DailyLog struct {
	LastCheckTime time.Time
	LogLevel      int
	LogCurSize    int
	FileFullName  string
	PFile         *os.File
	Log           *log.Logger
	M              sync.Mutex
}

func InitDailyLog(logDir string, logFile string, logStrLevel string) (*DailyLog, error) {
	logStrLevel = strings.ToLower(logStrLevel)

	if logStrLevel != "debug" && logStrLevel != "info" && logStrLevel != "warn" && logStrLevel != "error" {
		return nil, fmt.Errorf("wrong log level")
	}

	dailyLog := DailyLog{}
	dailyLog.LogLevel = getLogLevel(logStrLevel)
	if dailyLog.LogLevel == 0 {
		return nil, fmt.Errorf("wrong log level: %s", logStrLevel)
	}

	dailyLog.FileFullName = logDir + "/" + logFile

	var err error
	dailyLog.PFile, err = os.OpenFile(dailyLog.FileFullName, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		return nil, err
	}

	dailyLog.Log = log.New(dailyLog.PFile, "", log.Ldate|log.Lmicroseconds|log.Lshortfile)

	ff, err := os.Stat(dailyLog.FileFullName)
	if err != nil {
		return nil, err
	}

	if ff.Size() > math.MaxInt {
		dailyLog.LogCurSize = math.MaxInt
	} else {
		dailyLog.LogCurSize = int(ff.Size())
	}

	dailyLog.LastCheckTime = ff.ModTime()
	dailyLog.rotate()

	if GDailyLog == nil {
		GDailyLog = &dailyLog
	}

	return &dailyLog, nil
}

func (dLog *DailyLog) rotate() {
	now := time.Now()
	if now.Year() != dLog.LastCheckTime.Year() || now.Month() != dLog.LastCheckTime.Month()  || now.Day() != dLog.LastCheckTime.Day() {
		dLog.M.Lock()
		defer dLog.M.Unlock()

		if now.Equal(dLog.LastCheckTime) {
			return
		}

		err := dLog.PFile.Close()
		if err != nil {
			fmt.Printf("close log file err: %s\n", err.Error())
			return
		}

		err = os.Rename(dLog.FileFullName, dLog.FileFullName+"."+time.Now().Format("2006-01-02"))
		if err != nil {
			fmt.Printf("Rename log file err: %s\n", err.Error())
			return
		}

		dLog.PFile, err = os.OpenFile(dLog.FileFullName, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
		if err != nil {
			fmt.Printf("OpenFile log file err: %s\n", err.Error())
			return
		}

		dLog.Log.SetOutput(dLog.PFile)
		dLog.LastCheckTime = now
		dLog.LogCurSize = 0
	}
}

func (dLog *DailyLog) Debug(v ...interface{}) {
	if dLog.LogLevel == LogLevelDebug {
		out := fmt.Sprintf("[DEBUG][%d][%s] %s", goroutineid.GetGoID(), trace_id.GetTraceId(), fmt.Sprint(v...))
		dLog.LogCurSize += PrefixHeadLen + len(out)
		dLog.Log.Output(3, out)
		dLog.rotate()
	}
}

func (dLog *DailyLog) DebugF(format string, args ...interface{}) {
	if dLog.LogLevel == LogLevelDebug {
		msg := fmt.Sprintf(format, args...)
		out := fmt.Sprintf("[DEBUG][%d][%s] %s", goroutineid.GetGoID(), trace_id.GetTraceId(), msg)
		dLog.LogCurSize += PrefixHeadLen + len(out)
		dLog.Log.Output(3, out)
		dLog.rotate()
	}
}

func (dLog *DailyLog) Info(v ...interface{}) {
	if dLog.LogLevel >= LogLevelInfo {
		out := fmt.Sprintf("[INFO][%d][%s] %s", goroutineid.GetGoID(), trace_id.GetTraceId(), fmt.Sprint(v...))
		dLog.LogCurSize += PrefixHeadLen + len(out)
		dLog.Log.Output(3, out)
		dLog.rotate()
	}
}

func (dLog *DailyLog) InfoF(format string, args ...interface{}) {
	if dLog.LogLevel >= LogLevelInfo {
		msg := fmt.Sprintf(format, args...)
		out := fmt.Sprintf("[INFO][%d][%s] %s", goroutineid.GetGoID(), trace_id.GetTraceId(), msg)
		dLog.LogCurSize += PrefixHeadLen + len(out)
		dLog.Log.Output(3, out)
		dLog.rotate()
	}
}

func (dLog *DailyLog) Warn(v ...interface{}) {
	if dLog.LogLevel >= LogLevelWarn {
		out := fmt.Sprintf("[WARN][%d][%s] %s", goroutineid.GetGoID(), trace_id.GetTraceId(), fmt.Sprint(v...))
		dLog.LogCurSize += PrefixHeadLen + len(out)
		dLog.Log.Output(3, out)
		dLog.rotate()
	}
}

func (dLog *DailyLog) WarnF(format string, args ...interface{}) {
	if dLog.LogLevel >= LogLevelWarn {
		msg := fmt.Sprintf(format, args...)
		out := fmt.Sprintf("[WARN][%d][%s] %s", goroutineid.GetGoID(), trace_id.GetTraceId(), msg)
		dLog.LogCurSize += PrefixHeadLen + len(out)
		dLog.Log.Output(3, out)
		dLog.rotate()
	}
}

func (dLog *DailyLog) Error(v ...interface{}) {
	if dLog.LogLevel >= LogLevelError {
		out := fmt.Sprintf("[ERRO][%d][%s] %s", goroutineid.GetGoID(), trace_id.GetTraceId(), fmt.Sprint(v...))
		dLog.LogCurSize += PrefixHeadLen + len(out)
		dLog.Log.Output(3, out)
		dLog.rotate()
	}
}

func (dLog *DailyLog) ErrorF(format string, args ...interface{}) {
	if dLog.LogLevel >= LogLevelError {
		msg := fmt.Sprintf(format, args...)
		out := fmt.Sprintf("[ERRO][%d][%s] %s", goroutineid.GetGoID(), trace_id.GetTraceId(), msg)
		dLog.LogCurSize += PrefixHeadLen + len(out)
		dLog.Log.Output(3, out)
		dLog.rotate()
	}
}
