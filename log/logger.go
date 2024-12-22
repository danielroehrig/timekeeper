package log

import "log"

type Logger interface {
	Debugf(format string, args ...interface{})
	Infof(format string, args ...interface{})
	Warnf(format string, args ...interface{})
	Errorf(format string, args ...interface{})
}

type Impl struct {
	loglevel byte
}

const (
	LevelDebug byte = iota
	LevelInfo
	LevelWarn
	LevelError
)

var std = Impl{loglevel: LevelInfo}

func Debugf(format string, args ...interface{}) {
	std.Debugf(format, args...)
}

func Infof(format string, args ...interface{}) {
	std.Infof(format, args...)
}

func Warnf(format string, args ...interface{}) {
	std.Warnf(format, args...)
}

func Errorf(format string, args ...interface{}) {
	std.Errorf(format, args...)
}

func SetLogLevel(level byte) {
	std.loglevel = level
}

func (l *Impl) Debugf(format string, args ...interface{}) {
	if l.loglevel == LevelDebug {
		log.Printf("DEBUG: "+format, args...)
	}
}

func (l *Impl) Infof(format string, args ...interface{}) {
	if l.loglevel <= LevelInfo {
		log.Printf("INFO: "+format, args...)
	}
}

func (l *Impl) Warnf(format string, args ...interface{}) {
	if l.loglevel <= LevelWarn {
		log.Printf("WARN: "+format, args...)
	}
}

func (l *Impl) Errorf(format string, args ...interface{}) {
	log.Fatalf("ERROF: "+format, args...)
}
