package logger

import (
	"time"

	"github.com/ddosakura/sola/v2"
)

// Context Key
const (
	CtxLogger = "logger"
)

// Log Message
type Log struct {
	IsAction   bool
	Format     string
	V          []interface{}
	CreateTime time.Time
}

// Handler of Logger
type Handler func(*Log)

// New Logger
func New(bufSize int, h Handler) sola.Middleware {
	ch := make(chan *Log, bufSize)
	go func() {
		for {
			select {
			case l := <-ch:
				h(l)
			}
		}
	}()
	return func(next sola.Handler) sola.Handler {
		return func(c sola.Context) error {
			var tmp chan<- *Log = ch
			c.Set(CtxLogger, tmp)
			return next(c)
		}
	}
}

// Printf to Logger
func Printf(c sola.Context, format string, v ...interface{}) {
	ch := c.Get(CtxLogger).(chan<- *Log)
	ch <- &Log{false, format, v, time.Now()}
}

// Action to Logger
func Action(c sola.Context, v ...interface{}) {
	ch := c.Get(CtxLogger).(chan<- *Log)
	ch <- &Log{true, "", v, time.Now()}
}
