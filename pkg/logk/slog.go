package logk

import (
	"github.com/cloudwego/hertz/pkg/common/hlog"
	hertzslog "github.com/hertz-contrib/logger/slog"
	"gopkg.in/natefinch/lumberjack.v2"
	"log/slog"
	"time"
)

func _max[T int | int32 | int64 | uint | time.Duration](a, b T) T {
	if a > b {
		return a
	}
	return b
}

var defSlog = NewSLogger("./logs/default_slog.log", 10, 10, 7, hlog.LevelDebug)

func GetDefaultSLogger() *slog.Logger {
	return defSlog.Logger()
}

func SetSLoggerLevel(level hlog.Level) {
	defSlog.SetLevel(level)
}

func NewSLogger(logFile string, maxFileMB, maxBackups, maxDays int, level hlog.Level) *hertzslog.Logger {
	// 提供压缩和删除
	lumberjackLogger := &lumberjack.Logger{
		Filename:   logFile,
		MaxSize:    _max(20, maxFileMB), // 一个文件最大可达 20M。
		MaxBackups: _max(5, maxBackups), // 最多同时保存 5 个文件。
		MaxAge:     _max(1, maxDays),    // 一个文件最多可以保存 10 天。
		Compress:   true,                // 用 gzip 压缩。
	}
	//
	l := hertzslog.NewLogger(
		hertzslog.WithOutput(lumberjackLogger),
	)
	l.SetLevel(level)
	return l
}
