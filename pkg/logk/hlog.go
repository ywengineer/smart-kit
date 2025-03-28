package logk

import (
	"context"
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

// Fatal calls the default logger's Fatal method and then os.Exit(1).
func Fatal(v ...interface{}) {
	DefaultLogger().Fatal(v...)
}

// Error calls the default logger's Error method.
func Error(v ...interface{}) {
	DefaultLogger().Error(v...)
}

// Warn calls the default logger's Warn method.
func Warn(v ...interface{}) {
	DefaultLogger().Warn(v...)
}

// Notice calls the default logger's Notice method.
func Notice(v ...interface{}) {
	DefaultLogger().Notice(v...)
}

// Info calls the default logger's Info method.
func Info(v ...interface{}) {
	DefaultLogger().Info(v...)
}

// Debug calls the default logger's Debug method.
func Debug(v ...interface{}) {
	DefaultLogger().Debug(v...)
}

// Trace calls the default logger's Trace method.
func Trace(v ...interface{}) {
	DefaultLogger().Trace(v...)
}

// Fatalf calls the default logger's Fatalf method and then os.Exit(1).
func Fatalf(format string, v ...interface{}) {
	DefaultLogger().Fatalf(format, v...)
}

// Errorf calls the default logger's Errorf method.
func Errorf(format string, v ...interface{}) {
	DefaultLogger().Errorf(format, v...)
}

// Warnf calls the default logger's Warnf method.
func Warnf(format string, v ...interface{}) {
	DefaultLogger().Warnf(format, v...)
}

// Noticef calls the default logger's Noticef method.
func Noticef(format string, v ...interface{}) {
	DefaultLogger().Noticef(format, v...)
}

// Infof calls the default logger's Infof method.
func Infof(format string, v ...interface{}) {
	DefaultLogger().Infof(format, v...)
}

// Debugf calls the default logger's Debugf method.
func Debugf(format string, v ...interface{}) {
	DefaultLogger().Debugf(format, v...)
}

// Tracef calls the default logger's Tracef method.
func Tracef(format string, v ...interface{}) {
	DefaultLogger().Tracef(format, v...)
}

// CtxFatalf calls the default logger's CtxFatalf method and then os.Exit(1).
func CtxFatalf(ctx context.Context, format string, v ...interface{}) {
	DefaultLogger().CtxFatalf(ctx, format, v...)
}

// CtxErrorf calls the default logger's CtxErrorf method.
func CtxErrorf(ctx context.Context, format string, v ...interface{}) {
	DefaultLogger().CtxErrorf(ctx, format, v...)
}

// CtxWarnf calls the default logger's CtxWarnf method.
func CtxWarnf(ctx context.Context, format string, v ...interface{}) {
	DefaultLogger().CtxWarnf(ctx, format, v...)
}

// CtxNoticef calls the default logger's CtxNoticef method.
func CtxNoticef(ctx context.Context, format string, v ...interface{}) {
	DefaultLogger().CtxNoticef(ctx, format, v...)
}

// CtxInfof calls the default logger's CtxInfof method.
func CtxInfof(ctx context.Context, format string, v ...interface{}) {
	DefaultLogger().CtxInfof(ctx, format, v...)
}

// CtxDebugf calls the default logger's CtxDebugf method.
func CtxDebugf(ctx context.Context, format string, v ...interface{}) {
	DefaultLogger().CtxDebugf(ctx, format, v...)
}

// CtxTracef calls the default logger's CtxTracef method.
func CtxTracef(ctx context.Context, format string, v ...interface{}) {
	DefaultLogger().CtxTracef(ctx, format, v...)
}
