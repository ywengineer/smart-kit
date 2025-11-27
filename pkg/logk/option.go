package logk

import (
	"gopkg.in/natefinch/lumberjack.v2"
)

type option struct {
	file       string
	maxFileMB  int
	maxBackups int
	maxDays    int
	level      Level
	stdout     bool
}

type Option func(o *option)

func WithoutStdout() Option {
	return func(o *option) {
		o.stdout = false
	}
}

func WithMaxFileMB(maxFileMB int) Option {
	return func(o *option) {
		o.maxFileMB = maxFileMB
	}
}

func WithMaxBackups(maxBackups int) Option {
	return func(o *option) {
		o.maxBackups = maxBackups
	}
}

func WithMaxDays(maxDays int) Option {
	return func(o *option) {
		o.maxDays = maxDays
	}
}

func WithLevel(level Level) Option {
	return func(o *option) {
		o.level = level
	}
}

func newRollingLogger(o *option) *lumberjack.Logger {
	return &lumberjack.Logger{
		Filename:   o.file,
		MaxSize:    _max(20, o.maxFileMB),  // 一个文件最大可达 20M。
		MaxBackups: _max(20, o.maxBackups), // 最多同时保存 20 个文件。
		MaxAge:     _max(1, o.maxDays),     // 一个文件最多可以保存 10 天。
		Compress:   true,                   // 用 gzip 压缩。
	}
}
