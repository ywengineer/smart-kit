package logk

import (
	"log/slog"
	"os"
	"time"

	"github.com/cloudwego/hertz/pkg/common/hlog"
	hertzslog "github.com/hertz-contrib/logger/slog"
	"go.uber.org/zap/zapcore"
)

func _max[T int | int32 | int64 | uint | time.Duration](a, b T) T {
	if a > b {
		return a
	}
	return b
}

var defSlog = NewSLogger("./logs/default_slog.log")

func GetDefaultSLogger() *slog.Logger {
	return defSlog.Logger()
}

func SetSLoggerLevel(level hlog.Level) {
	defSlog.SetLevel(level)
}

func NewSLogger(logFile string, opts ...Option) *hertzslog.Logger {
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
	var l *hertzslog.Logger
	if o.stdout {
		l = hertzslog.NewLogger(
			hertzslog.WithOutput(zapcore.NewMultiWriteSyncer(zapcore.AddSync(os.Stdout), zapcore.AddSync(newRollingLogger(o)))),
		)
	} else {
		l = hertzslog.NewLogger(
			hertzslog.WithOutput(zapcore.AddSync(newRollingLogger(o))),
		)
	}
	l.SetLevel(hlog.Level(o.level))
	return l
}
