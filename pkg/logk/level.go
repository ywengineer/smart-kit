package logk

import (
	"github.com/cloudwego/hertz/pkg/common/hlog"
	"gopkg.in/yaml.v3"
)

type Level hlog.Level

const InfoLevel = Level(hlog.LevelInfo)
const DebugLevel = Level(hlog.LevelDebug)
const WarnLevel = Level(hlog.LevelWarn)
const ErrorLevel = Level(hlog.LevelError)
const FatalLevel = Level(hlog.LevelFatal)

func (l *Level) String() string {
	switch hlog.Level(*l) {
	case hlog.LevelTrace:
		return "trace"
	case hlog.LevelDebug:
		return "debug"
	case hlog.LevelInfo:
		return "info"
	case hlog.LevelNotice, hlog.LevelWarn:
		return "warn"
	case hlog.LevelError:
		return "error"
	case hlog.LevelFatal:
		return "fatal"
	default:
		return "info"
	}
}

func (l *Level) from(v string) {
	switch v {
	case "trace":
		*l = Level(hlog.LevelTrace)
	case "debug":
		*l = Level(hlog.LevelDebug)
	case "info":
		*l = Level(hlog.LevelInfo)
	case "warn":
		*l = Level(hlog.LevelWarn)
	case "error":
		*l = Level(hlog.LevelError)
	case "fatal":
		*l = Level(hlog.LevelFatal)
	default:
		*l = Level(hlog.LevelDebug)
	}
}

func (l *Level) UnmarshalJSON(bytes []byte) error {
	l.from(string(bytes))
	return nil
}

func (l *Level) UnmarshalText(text []byte) error {
	l.from(string(text))
	return nil
}

func (l *Level) UnmarshalYAML(value *yaml.Node) error {
	var s string
	if err := value.Decode(&s); err != nil {
		return err
	}
	l.from(s)
	return nil
}
