package logk

import (
	"os"

	"github.com/cloudwego/hertz/pkg/common/hlog"
	hertzzap "github.com/hertz-contrib/logger/zap"
	"github.com/samber/lo"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func NewZapLogger(logFile string, opts ...Option) *hertzzap.Logger {
	o := &option{
		file:       logFile,
		maxFileMB:  20,
		maxBackups: 20,
		maxDays:    7,
		level:      DebugLevel,
		stdout:     true,
	}
	for _, opt := range opts {
		opt(o)
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
		hertzzap.WithCoreWs(lo.If(o.stdout, zapcore.NewMultiWriteSyncer(zapcore.AddSync(os.Stdout), zapcore.AddSync(newRollingLogger(o)))).Else(zapcore.AddSync(newRollingLogger(o)))),
		hertzzap.WithZapOptions(
			zap.AddStacktrace(zapcore.PanicLevel),
		),
	)
	l.SetLevel(hlog.Level(o.level))
	return l
}
