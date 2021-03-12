package log

import (
	"fmt"
	"github.com/chanxuehong/util/math"
	"github.com/diversability/gocom/goroutineid"
	"github.com/diversability/gocom/trace_id"
	"log"
	"os"
	"strings"
	"time"
)

// 按文件大小进行日志切割

var GSizeLog *SizeLog

type SizeLog struct {
	LogLevel     int
	LogMaxSize   int
	LogCurSize   int
	FileFullName string
	PFile        *os.File
	Log          *log.Logger
}

func InitSizeLog(logDir string, logFile string, logStrLevel string, LogMaxSize int) (*SizeLog, error) {
	if LogMaxSize == 0 {
		LogMaxSize = 524288000
	}

	logStrLevel = strings.ToLower(logStrLevel)

	if logStrLevel != "debug" && logStrLevel != "info" && logStrLevel != "warn" && logStrLevel != "error" {
		return nil, fmt.Errorf("wrong log level")
	}

	sizeLog := SizeLog{LogMaxSize: LogMaxSize}
	sizeLog.LogLevel = getLogLevel(logStrLevel)
	if sizeLog.LogLevel == 0 {
		return nil, fmt.Errorf("wrong log level: %s", logStrLevel)
	}

	sizeLog.FileFullName = logDir + "/" + logFile

	var err error
	sizeLog.PFile, err = os.OpenFile(sizeLog.FileFullName, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		return nil, err
	}

	sizeLog.Log = log.New(sizeLog.PFile, "", log.Ldate|log.Lmicroseconds|log.Lshortfile)
	if GSizeLog == nil {
		GSizeLog = &sizeLog
	}

	ff, err := os.Stat(sizeLog.FileFullName)
	if err != nil {
		return nil, err
	}

	// 文件最大2G
	if ff.Size() > math.MaxInt {
		sizeLog.LogCurSize = math.MaxInt
	} else {
		sizeLog.LogCurSize = int(ff.Size())
	}

	sizeLog.rotate()
	return &sizeLog, nil
}

func (slog *SizeLog) rotate() {
	if slog.LogCurSize >= slog.LogMaxSize {
		err := slog.PFile.Close()
		if err != nil {
			fmt.Printf("close log file err: %s\n", err.Error())
			return
		}

		err = os.Rename(slog.FileFullName, slog.FileFullName+"."+time.Now().Format("15:04:05.999999999"))
		if err != nil {
			fmt.Printf("Rename log file err: %s\n", err.Error())
			return
		}

		slog.PFile, err = os.OpenFile(slog.FileFullName, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
		if err != nil {
			fmt.Printf("OpenFile log file err: %s\n", err.Error())
			return
		}

		slog.Log.SetOutput(slog.PFile)
		slog.LogCurSize = 0
	}
}

func (slog *SizeLog) Debug(v ...interface{}) {
	if slog.LogLevel == LogLevelDebug {
		out := fmt.Sprintf("[DEBUG][%d][%s] %s", goroutineid.GetGoID(), trace_id.GetTraceId(), fmt.Sprint(v...))
		slog.LogCurSize += PrefixHeadLen + len(out)
		slog.Log.Output(3, out)
		slog.rotate()
	}
}

func (slog *SizeLog) DebugF(format string, args ...interface{}) {
	if slog.LogLevel == LogLevelDebug {
		msg := fmt.Sprintf(format, args...)
		out := fmt.Sprintf("[DEBUG][%d][%s] %s", goroutineid.GetGoID(), trace_id.GetTraceId(), msg)
		slog.LogCurSize += PrefixHeadLen + len(out)
		slog.Log.Output(3, out)
		slog.rotate()
	}
}

func (slog *SizeLog) Info(v ...interface{}) {
	if slog.LogLevel >= LogLevelInfo {
		out := fmt.Sprintf("[INFO][%d][%s] %s", goroutineid.GetGoID(), trace_id.GetTraceId(), fmt.Sprint(v...))
		slog.LogCurSize += PrefixHeadLen + len(out)
		slog.Log.Output(3, out)
		slog.rotate()
	}
}

func (slog *SizeLog) InfoF(format string, args ...interface{}) {
	if slog.LogLevel >= LogLevelInfo {
		msg := fmt.Sprintf(format, args...)
		out := fmt.Sprintf("[INFO][%d][%s] %s", goroutineid.GetGoID(), trace_id.GetTraceId(), msg)
		slog.LogCurSize += PrefixHeadLen + len(out)
		slog.Log.Output(3, out)
		slog.rotate()
	}
}

func (slog *SizeLog) Warn(v ...interface{}) {
	if slog.LogLevel >= LogLevelWarn {
		out := fmt.Sprintf("[WARN][%d][%s] %s", goroutineid.GetGoID(), trace_id.GetTraceId(), fmt.Sprint(v...))
		slog.LogCurSize += PrefixHeadLen + len(out)
		slog.Log.Output(3, out)
		slog.rotate()
	}
}

func (slog *SizeLog) WarnF(format string, args ...interface{}) {
	if slog.LogLevel >= LogLevelWarn {
		msg := fmt.Sprintf(format, args...)
		out := fmt.Sprintf("[WARN][%d][%s] %s", goroutineid.GetGoID(), trace_id.GetTraceId(), msg)
		slog.LogCurSize += PrefixHeadLen + len(out)
		slog.Log.Output(3, out)
		slog.rotate()
	}
}

func (slog *SizeLog) Error(v ...interface{}) {
	if slog.LogLevel >= LogLevelError {
		out := fmt.Sprintf("[ERRO][%d][%s] %s", goroutineid.GetGoID(), trace_id.GetTraceId(), fmt.Sprint(v...))
		slog.LogCurSize += PrefixHeadLen + len(out)
		slog.Log.Output(3, out)
		slog.rotate()
	}
}

func (slog *SizeLog) ErrorF(format string, args ...interface{}) {
	if slog.LogLevel >= LogLevelError {
		msg := fmt.Sprintf(format, args...)
		out := fmt.Sprintf("[ERRO][%d][%s] %s", goroutineid.GetGoID(), trace_id.GetTraceId(), msg)
		slog.LogCurSize += PrefixHeadLen + len(out)
		slog.Log.Output(3, out)
		slog.rotate()
	}
}
