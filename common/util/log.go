package util

import (
	"fmt"
	"github.com/yudai/pp"
	"time"
)

type level string

const (
	INFO  level = "info"
	WARN  level = "warn"
	ERROR level = "error"
)

type log struct {
	field string
	now   time.Time
}

func Logger() LogImpl {
	return &log{
		field: "common",
		now:   time.Now(),
	}
}

type LogImpl interface {
	Log(level level, format string, args ...interface{}) LogImpl
	SetField(field string) LogImpl
	Cost() LogImpl
	ReSetTimer() LogImpl
}

func (l log) Log(level level, format string, args ...interface{}) LogImpl {
	_, _ = pp.Printf(fmt.Sprintf("[%s] %s: %s \n", l.field, level, format), args...)
	return l
}

func (l log) SetField(field string) LogImpl {
	l.field = field
	return l
}

func (l log) Cost() LogImpl {
	_, _ = pp.Printf(fmt.Sprintf("[%s] cost: %+v \n", l.field, time.Since(l.now)))
	return l.ReSetTimer()
}

func (l log) ReSetTimer() LogImpl {
	l.now = time.Now()
	return l
}
