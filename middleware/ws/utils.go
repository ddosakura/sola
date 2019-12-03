package ws

import (
	"encoding/json"

	uuid "github.com/satori/go.uuid"
)

// XSend with json
type XSend func(uuid.UUID, interface{}) error

// SendWrap Util
func SendWrap(send Send) XSend {
	return func(UUID uuid.UUID, v interface{}) error {
		bs, e := json.Marshal(v)
		if e != nil {
			return e
		}
		send(UUID, bs)
		return nil
	}
}

// NewModel for XHandle
type NewModel func() interface{}

// XHandle with json
type XHandle func(uuid.UUID, interface{}) error

// HandleWrap Util
func HandleWrap(build NewModel, h XHandle) Handler {
	return func(UUID uuid.UUID, bs []byte) error {
		v := build()
		if e := json.Unmarshal(bs, v); e != nil {
			return e
		}
		return h(UUID, v)
	}
}
