package easylogger

import (
	"fmt"
	"io"
)

// LoggerWriter ...
// log writer for library
type LoggerWriter struct {
	f bool
	w io.Writer
}

func (l *LoggerWriter) SetDebug(f bool) *LoggerWriter {
	l.f = f
	return l
}

func (l *LoggerWriter) SetLogger(w io.Writer) *LoggerWriter {
	l.w = w
	return l
}

func (l *LoggerWriter) LoggerRaw(p string) {
	if l.f && l.w != nil {
		_, _ = l.w.Write([]byte(p))
	}
}

func (l *LoggerWriter) Logger(p string, ex ...interface{}) {
	s := fmt.Sprintf(p, ex...)
	if l.f && l.w != nil {
		_, _ = l.w.Write([]byte(s))
	}
}
