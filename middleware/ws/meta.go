package ws

import (
	"errors"

	"github.com/ddosakura/sola/v2"
)

// Ctx
const (
	Issuer = "sola/ws"
)

// meta
var (
	ALL = [16]byte{}
)

// error(s)
var (
	ErrOption = errors.New("Must set Handler")
	ErrNoUUID = errors.New("No UUID")

	errNoUUID = sola.IError(Issuer, ErrNoUUID)
)
