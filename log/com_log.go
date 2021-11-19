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

// 按文件大小和天进行日志切割

var GLog *ComLog

const PrefixHeadLen = 48

type ComLog struct {
	LogLevel       int
	LastRotateTime time.Time
	LogMaxSize     int
	LogCurSize     int
	FileFullName   string
	PFile          *os.File
	Log            *log.Logger
	M              sync.Mutex
}

func InitLog(logDir string, logFile string, logStrLevel string, LogMaxSize int) (*ComLog, error) {
	if LogMaxSize == 0 {
		LogMaxSize = 524288000
	}

	logStrLevel = strings.ToLower(logStrLevel)

	if logStrLevel != "debug" && logStrLevel != "info" && logStrLevel != "warn" && logStrLevel != "error" {
		return nil, fmt.Errorf("wrong log level")
	}

	comLog := ComLog{LogMaxSize: LogMaxSize}
	comLog.LogLevel = getLogLevel(logStrLevel)
	if comLog.LogLevel == 0 {
		return nil, fmt.Errorf("wrong log level: %s", logStrLevel)
	}

	comLog.FileFullName = logDir + "/" + logFile

	var err error
	comLog.PFile, err = os.OpenFile(comLog.FileFullName, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		return nil, err
	}

	comLog.Log = log.New(comLog.PFile, "", log.Ldate|log.Lmicroseconds|log.Lshortfile)
	if GLog == nil {
		GLog = &comLog
	}

	ff, err := os.Stat(comLog.FileFullName)
	if err != nil {
		return nil, err
	}

	comLog.LastRotateTime = ff.ModTime()

	// 文件最大2G
	if ff.Size() > math.MaxInt {
		comLog.LogCurSize = math.MaxInt
	} else {
		comLog.LogCurSize = int(ff.Size())
	}

	comLog.rotate()
	return &comLog, nil
}

func (cLog *ComLog) rotate() {
	now := time.Now()
	if now.Year() != cLog.LastRotateTime.Year() || now.Month() != cLog.LastRotateTime.Month() || now.Day() != cLog.LastRotateTime.Day() {
		cLog.M.Lock()
		defer cLog.M.Unlock()

		if now.Equal(cLog.LastRotateTime) {
			return
		}

		err := cLog.PFile.Close()
		if err != nil {
			fmt.Printf("close log file err: %s\n", err.Error())
			return
		}

		err = os.Rename(cLog.FileFullName, cLog.FileFullName+"."+time.Now().Format("2006-01-02 15:04:05.999999999"))
		if err != nil {
			fmt.Printf("Rename log file err: %s\n", err.Error())
			return
		}

		cLog.PFile, err = os.OpenFile(cLog.FileFullName, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
		if err != nil {
			fmt.Printf("OpenFile log file err: %s\n", err.Error())
			return
		}

		cLog.Log.SetOutput(cLog.PFile)
		cLog.LastRotateTime = now
		cLog.LogCurSize = 0
		return
	}

	if cLog.LogCurSize > cLog.LogMaxSize {
		cLog.M.Lock()
		defer cLog.M.Unlock()

		if cLog.LogCurSize == 0 {
			return
		}

		err := cLog.PFile.Close()
		if err != nil {
			fmt.Printf("close log file err: %s\n", err.Error())
			return
		}

		err = os.Rename(cLog.FileFullName, cLog.FileFullName+"."+time.Now().Format("2006-01-02 15:04:05.999999999"))
		if err != nil {
			fmt.Printf("Rename log file err: %s\n", err.Error())
			return
		}

		cLog.PFile, err = os.OpenFile(cLog.FileFullName, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
		if err != nil {
			fmt.Printf("OpenFile log file err: %s\n", err.Error())
			return
		}

		cLog.Log.SetOutput(cLog.PFile)
		cLog.LogCurSize = 0
	}
}

func (cLog *ComLog) Debug(v ...interface{}) {
	if cLog.LogLevel == LogLevelDebug {
		out := fmt.Sprintf("[DEBUG][%d][%s] %s", goroutineid.GetGoID(), trace_id.GetTraceId(), fmt.Sprint(v...))
		cLog.LogCurSize += PrefixHeadLen + len(out)
		cLog.Log.Output(3, out)
		cLog.rotate()
	}
}

func (cLog *ComLog) DebugF(format string, args ...interface{}) {
	if cLog.LogLevel == LogLevelDebug {
		msg := fmt.Sprintf(format, args...)
		out := fmt.Sprintf("[DEBUG][%d][%s] %s", goroutineid.GetGoID(), trace_id.GetTraceId(), msg)
		cLog.LogCurSize += PrefixHeadLen + len(out)
		cLog.Log.Output(3, out)
		cLog.rotate()
	}
}

func (cLog *ComLog) Info(v ...interface{}) {
	if cLog.LogLevel >= LogLevelInfo {
		out := fmt.Sprintf("[INFO][%d][%s] %s", goroutineid.GetGoID(), trace_id.GetTraceId(), fmt.Sprint(v...))
		cLog.LogCurSize += PrefixHeadLen + len(out)
		cLog.Log.Output(3, out)
		cLog.rotate()
	}
}

func (cLog *ComLog) InfoF(format string, args ...interface{}) {
	if cLog.LogLevel >= LogLevelInfo {
		msg := fmt.Sprintf(format, args...)
		out := fmt.Sprintf("[INFO][%d][%s] %s", goroutineid.GetGoID(), trace_id.GetTraceId(), msg)
		cLog.LogCurSize += PrefixHeadLen + len(out)
		cLog.Log.Output(3, out)
		cLog.rotate()
	}
}

func (cLog *ComLog) Warn(v ...interface{}) {
	if cLog.LogLevel >= LogLevelWarn {
		out := fmt.Sprintf("[WARN][%d][%s] %s", goroutineid.GetGoID(), trace_id.GetTraceId(), fmt.Sprint(v...))
		cLog.LogCurSize += PrefixHeadLen + len(out)
		cLog.Log.Output(3, out)
		cLog.rotate()
	}
}

func (cLog *ComLog) WarnF(format string, args ...interface{}) {
	if cLog.LogLevel >= LogLevelWarn {
		msg := fmt.Sprintf(format, args...)
		out := fmt.Sprintf("[WARN][%d][%s] %s", goroutineid.GetGoID(), trace_id.GetTraceId(), msg)
		cLog.LogCurSize += PrefixHeadLen + len(out)
		cLog.Log.Output(3, out)
		cLog.rotate()
	}
}

func (cLog *ComLog) Error(v ...interface{}) {
	if cLog.LogLevel >= LogLevelError {
		out := fmt.Sprintf("[ERRO][%d][%s] %s", goroutineid.GetGoID(), trace_id.GetTraceId(), fmt.Sprint(v...))
		cLog.LogCurSize += PrefixHeadLen + len(out)
		cLog.Log.Output(3, out)
		cLog.rotate()
	}
}

func (cLog *ComLog) ErrorF(format string, args ...interface{}) {
	if cLog.LogLevel >= LogLevelError {
		msg := fmt.Sprintf(format, args...)
		out := fmt.Sprintf("[ERRO][%d][%s] %s", goroutineid.GetGoID(), trace_id.GetTraceId(), msg)
		cLog.LogCurSize += PrefixHeadLen + len(out)
		cLog.Log.Output(3, out)
		cLog.rotate()
	}
}

func Debug(v ...interface{}) {
	GLog.Debug(v...)
}

func DebugF(format string, args ...interface{}) {
	GLog.DebugF(format, args...)
}

func Info(v ...interface{}) {
	GLog.Info(v...)
}

func InfoF(format string, args ...interface{}) {
	GLog.InfoF(format, args...)
}

func Warn(v ...interface{}) {
	GLog.Warn(v...)
}

func WarnF(format string, args ...interface{}) {
	GLog.WarnF(format, args...)
}

func Error(v ...interface{}) {
	GLog.Error(v...)
}

func ErrorF(format string, args ...interface{}) {
	GLog.ErrorF(format, args...)
}

var (
	green   = string([]byte{27, 91, 57, 55, 59, 52, 50, 109})
	white   = string([]byte{27, 91, 57, 48, 59, 52, 55, 109})
	yellow  = string([]byte{27, 91, 57, 55, 59, 52, 51, 109})
	red     = string([]byte{27, 91, 57, 55, 59, 52, 49, 109})
	blue    = string([]byte{27, 91, 57, 55, 59, 52, 52, 109})
	magenta = string([]byte{27, 91, 57, 55, 59, 52, 53, 109})
	cyan    = string([]byte{27, 91, 57, 55, 59, 52, 54, 109})
	Reset   = string([]byte{27, 91, 48, 109})
)

func ColorForStatus(code int) string {
	switch {
	case code >= 200 && code <= 299:
		return green
	case code >= 300 && code <= 399:
		return white
	case code >= 400 && code <= 499:
		return yellow
	default:
		return red
	}
}

func ColorForMethod(method string) string {
	switch {
	case method == "GET":
		return blue
	case method == "POST":
		return cyan
	case method == "PUT":
		return yellow
	case method == "DELETE":
		return red
	case method == "PATCH":
		return green
	case method == "HEAD":
		return magenta
	case method == "OPTIONS":
		return white
	default:
		return Reset
	}
}

func ColorForReset() string {
	return Reset
}

const (
	LogLevelError = 1 << iota
	LogLevelWarn
	LogLevelInfo
	LogLevelDebug
)

func getLogLevel(logStrLevel string) int {
	if logStrLevel == "debug" {
		return LogLevelDebug
	}

	if logStrLevel == "info" {
		return LogLevelInfo
	}

	if logStrLevel == "warn" {
		return LogLevelWarn
	}

	if logStrLevel == "error" {
		return LogLevelError
	}

	return 0
}