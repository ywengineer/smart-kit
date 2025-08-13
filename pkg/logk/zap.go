package logk

import (
	"github.com/cloudwego/hertz/pkg/common/hlog"
	hertzzap "github.com/hertz-contrib/logger/zap"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
	"os"
)

func NewZapLogger(logFile string, maxFileMB, maxBackups, maxDays int, level hlog.Level) hlog.FullLogger {
	// 提供压缩和删除
	lumberjackLogger := &lumberjack.Logger{
		Filename:   logFile,
		MaxSize:    _max(20, maxFileMB), // 一个文件最大可达 20M。
		MaxBackups: _max(5, maxBackups), // 最多同时保存 5 个文件。
		MaxAge:     _max(1, maxDays),    // 一个文件最多可以保存 10 天。
		Compress:   true,                // 用 gzip 压缩。
	}
	//
	l := hertzzap.NewLogger(
		hertzzap.WithCoreEnc(zapcore.NewJSONEncoder(zapcore.EncoderConfig{
			MessageKey:     "msg",
			LevelKey:       "level",
			NameKey:        "name",
			TimeKey:        "ts",
			CallerKey:      "caller",
			FunctionKey:    "func",
			StacktraceKey:  "stacktrace",
			LineEnding:     "\n",
			EncodeTime:     zapcore.ISO8601TimeEncoder,
			EncodeLevel:    zapcore.CapitalLevelEncoder,
			EncodeDuration: zapcore.StringDurationEncoder,
			EncodeCaller:   zapcore.ShortCallerEncoder,
		})),
		hertzzap.WithExtraKeys([]hertzzap.ExtraKey{"data"}),
		hertzzap.WithExtraKeyAsStr(),
		hertzzap.WithCoreWs(zapcore.NewMultiWriteSyncer(zapcore.AddSync(os.Stdout), zapcore.AddSync(lumberjackLogger))),
		hertzzap.WithZapOptions(
			zap.AddStacktrace(zapcore.ErrorLevel),
		),
	)
	l.SetLevel(level)
	return l
}
