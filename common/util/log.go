package util

import (
	"fmt"
	"github.com/yudai/pp"
	"sync"
	"time"
)

type level string

const (
	INFO  level = "info"
	WARN  level = "warn"
	ERROR level = "error"
)

type log struct {
	wg    *sync.WaitGroup
	field string
	now   time.Time
}

func Logger() LogImpl {
	return &log{
		wg:    &sync.WaitGroup{},
		field: "common",
		now:   time.Now(),
	}
}

type LogImpl interface {
	Log(level level, format string, args ...interface{}) LogImpl
	SetField(field string) LogImpl
	Cost()
}

func (l log) Log(level level, format string, args ...interface{}) LogImpl {
	l.wg.Add(1)
	go func() {
		_, _ = pp.Printf(fmt.Sprintf("[%s] %s: %s \n", l.field, level, format), args...)
		l.wg.Done()
	}()
	return l
}

func (l log) SetField(field string) LogImpl {
	l.field = field
	return l
}

func (l log) Cost() {
	l.wg.Add(1)
	go func() {
		_, _ = pp.Printf(fmt.Sprintf("[%s] cost: %+v \n", l.field, time.Since(l.now)))
		l.wg.Done()
	}()
	l.wg.Wait()
}
