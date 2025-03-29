package logk

import (
	"github.com/cloudwego/hertz/pkg/common/hlog"
	hertzzap "github.com/hertz-contrib/logger/zap"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
	"os"
	"time"
)

func Max[T int | int32 | int64 | uint | time.Duration](a, b T) T {
	if a > b {
		return a
	}
	return b
}

func NewLogger(logFile string, maxFileMB, maxBackups, maxDays int, level hlog.Level) hlog.FullLogger {
	// 提供压缩和删除
	lumberjackLogger := &lumberjack.Logger{
		Filename:   logFile,
		MaxSize:    Max(20, maxFileMB), // 一个文件最大可达 20M。
		MaxBackups: Max(5, maxBackups), // 最多同时保存 5 个文件。
		MaxAge:     Max(1, maxDays),    // 一个文件最多可以保存 10 天。
		Compress:   true,               // 用 gzip 压缩。
	}
	//
	hlog.SetLogger(hertzzap.NewLogger(
		hertzzap.WithCoreWs(zapcore.NewMultiWriteSyncer(zapcore.AddSync(os.Stdout), zapcore.AddSync(lumberjackLogger))),
		hertzzap.WithZapOptions(
			zap.AddStacktrace(zapcore.ErrorLevel),
		),
	))
	hlog.SetLevel(level)
	return hlog.DefaultLogger()
}
