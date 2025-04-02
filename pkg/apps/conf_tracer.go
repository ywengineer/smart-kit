package apps

import (
	"github.com/cloudwego/hertz/pkg/common/tracer/stats"
	"gopkg.in/yaml.v3"
)

type TraceLevel stats.Level

func (l *TraceLevel) String() string {
	switch stats.Level(*l) {
	case stats.LevelDisabled:
		return "disable"
	case stats.LevelBase:
		return "base"
	case stats.LevelDetailed:
		return "detail"
	default:
		return "base"
	}
}

func (l *TraceLevel) from(v string) {
	switch v {
	case "disable", "none":
		*l = TraceLevel(stats.LevelDisabled)
	case "base":
		*l = TraceLevel(stats.LevelBase)
	case "detail":
		*l = TraceLevel(stats.LevelDetailed)
	default:
		*l = TraceLevel(stats.LevelDisabled)
	}
}

func (l *TraceLevel) UnmarshalJSON(bytes []byte) error {
	l.from(string(bytes))
	return nil
}

func (l *TraceLevel) UnmarshalText(text []byte) error {
	l.from(string(text))
	return nil
}

func (l *TraceLevel) UnmarshalYAML(value *yaml.Node) error {
	var s string
	if err := value.Decode(&s); err != nil {
		return err
	}
	l.from(s)
	return nil
}
