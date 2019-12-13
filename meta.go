package sola

import (
	"errors"
	"html/template"
	"io"
)

const (
	// Version of Sola
	Version = "v2.1.2"
	// Issuer of InternalError
	Issuer = "sola"
)

// error(s)
var (
	ErrM2H = errors.New("Middleware using next() can't be Handler")

	errM2H = IError(Issuer, ErrM2H)
)

// ContextKey
const (
	CtxSola     = "sola"
	CtxRequest  = "sola.request"
	CtxResponse = "sola.response"
)

// TplInternalError for sola
var TplInternalError = template.Must(template.
	New("SolaInternalError").
	Parse(`<html>
<head>
<meta charset="utf-8">
<title>Sola Internal Error</title>
<style>
body {
	margin-top: 50px;
}

th {
    color: #4f6b72;
    border-right: 1px solid #C1DAD7;
    border-bottom: 1px solid #C1DAD7;
    border-top: 1px solid #C1DAD7;
    letter-spacing: 2px;
    text-transform: uppercase;
    text-align: left;
    padding: 6px 6px 6px 12px;
}

th.top {
    border-left: 1px solid #C1DAD7;
    border-right: 1px solid #C1DAD7;
}

th.spec {
    border-left: 1px solid #C1DAD7;
    border-top: 0;
    color: #797268;
}

td {
    border-right: 1px solid #C1DAD7;
    border-bottom: 1px solid #C1DAD7;
    background: #fff;
    font-size: 11px;
    padding: 6px 6px 6px 12px;
    color: #4f6b72;
}

tr:nth-child(odd) td {
    background: #F5FAFA;
    color: #797268;
}
</style>
</head>
<body>
<h1 style="text-align: center;">Sola Internal Error</h1>
<table cellspacing="0" style="width: 70%;margin: auto;">
<tr>
    <th class="top">Meta</th>
    <th>Value</th>
</tr>
{{range $i,$e := .}}
<tr>
    <th class="spec">{{$i}}</td>
    <td>{{$e}}</td>
</tr>
{{end}}
</table>
</body>
</html>`))

// InternalError of sola
type InternalError struct {
	error
	issuer string
}

// IError Internal Error
func IError(issuer string, e error) *InternalError {
	return &InternalError{e, issuer}
}

// Error Internal
func Error(issuer string, e error) error {
	return IError(issuer, e)
}

// Issuer of InternalError
func (e *InternalError) Issuer() string {
	return e.issuer
}

// Handler of InternalError
func (e *InternalError) Handler() Handler {
	return func(c Context) error {
		return e
	}
}

// Write to io.Writer
func (e *InternalError) Write(w io.Writer) error {
	errStr := e.error.Error()
	return TplInternalError.Execute(w, map[string]string{
		"Sola Version": Version,
		"Issuer":       e.Issuer(),
		"Error":        errStr,
	})
}
