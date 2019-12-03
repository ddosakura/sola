package ws

import (
	"errors"
	"io"

	"github.com/ddosakura/sola/v2"
	uuid "github.com/satori/go.uuid"
	"golang.org/x/net/websocket"
)

// Handler of WebSocket
type Handler func(uuid.UUID, []byte) error

// ErrorHandler of WebSocket
type ErrorHandler func(uuid.UUID, error)

// Option of WebSocket
type Option struct {
	Handler Handler

	// Optional
	First        func(uuid.UUID)
	ReceiveError ErrorHandler
	SendError    ErrorHandler
	HandlerError ErrorHandler
}

// meta
var (
	ALL = [16]byte{}
)

// error(s)
var (
	ErrOption = errors.New("Must set Handler")
	ErrNoUUID = errors.New("No UUID")
)

func first(uuid.UUID)       {}
func pass(uuid.UUID, error) {}

// Send Action
type Send func(uuid.UUID, []byte) error

// New WebSocket Handler
func New(o *Option) (sola.Handler, Send) {
	if o == nil {
		panic(ErrOption)
	}
	if o.Handler == nil {
		panic(ErrOption)
	}
	if o.First == nil {
		o.First = first
	}
	if o.ReceiveError == nil {
		o.ReceiveError = pass
	}
	if o.SendError == nil {
		o.SendError = pass
	}
	if o.HandlerError == nil {
		o.HandlerError = pass
	}

	m := make(map[uuid.UUID]*websocket.Conn)
	h := func(c sola.Context) error {
		w, r := c.Response(), c.Request()
		h := handle(o, m)
		h.ServeHTTP(w, r)
		return nil
	}
	send := func(UUID uuid.UUID, v []byte) error {
		if UUID == ALL {
			for id, w := range m {
				go sendMessage(o, id, w, v)
			}
			return nil
		}
		w := m[UUID]
		if w == nil {
			return ErrNoUUID
		}
		go sendMessage(o, UUID, m[UUID], v)
		return nil
	}
	return h, send
}

func sendMessage(o *Option, UUID uuid.UUID, w *websocket.Conn, v interface{}) {
	if e := websocket.Message.Send(w, v); e != nil {
		o.SendError(UUID, e)
	}
}

func handle(o *Option, m map[uuid.UUID]*websocket.Conn) websocket.Handler {
	return websocket.Handler(func(w *websocket.Conn) {
		UUID := uuid.NewV4()
		m[UUID] = w
		o.First(UUID)
		for {
			var v []byte
			if e := websocket.Message.Receive(w, &v); e != nil {
				if e == io.EOF {
					break
				}
				o.ReceiveError(UUID, e)
			}
			go func() {
				if e := o.Handler(UUID, v); e != nil {
					o.HandlerError(UUID, e)
				}
			}()
		}
		delete(m, UUID)
	})
}
