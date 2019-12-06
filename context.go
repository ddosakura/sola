package sola

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"
	"os"
)

func newContext() *context {
	return &context{
		store: map[string]interface{}{},
	}
}

// === Set/Get ===

// Origin Context
func (c *context) Origin() Context {
	return c.origin
}

// Shadow Context will search origin context if no match key
func (c *context) Shadow() Context {
	shadow := newContext()
	shadow.origin = c
	return shadow
}

// Set Ctx
func (c *context) Set(key string, value interface{}) {
	c.lock.Lock()
	defer c.lock.Unlock()
	c.store[key] = value
}

// Get Ctx
func (c *context) Get(key string) interface{} {
	c.lock.RLock()
	defer c.lock.RUnlock()
	v := c.store[key]
	if v == nil && c.origin != nil {
		v = c.origin.Get(key)
	}
	return v
}

// === API ===

// Sola Impl
func (c *context) Sola() *Sola {
	return c.Get(CtxSola).(*Sola)
}

// SetCookie proxy http.SetCookie
func (c *context) SetCookie(cookie *http.Cookie) {
	http.SetCookie(c.Response(), cookie)
}

// Request in context
func (c *context) Request() *http.Request {
	return c.Get(CtxRequest).(*http.Request)
}

// Response in context
func (c *context) Response() http.ResponseWriter {
	return c.Get(CtxResponse).(http.ResponseWriter)
}

// === Writer ===

const (
	charsetUTF8 = "charset=UTF-8"
)

// MIME types
const (
	MIMEApplicationJSON                  = "application/json"
	MIMEApplicationJSONCharsetUTF8       = MIMEApplicationJSON + "; " + charsetUTF8
	MIMEApplicationJavaScript            = "application/javascript"
	MIMEApplicationJavaScriptCharsetUTF8 = MIMEApplicationJavaScript + "; " + charsetUTF8
	MIMEApplicationXML                   = "application/xml"
	MIMEApplicationXMLCharsetUTF8        = MIMEApplicationXML + "; " + charsetUTF8
	MIMETextXML                          = "text/xml"
	MIMETextXMLCharsetUTF8               = MIMETextXML + "; " + charsetUTF8
	MIMEApplicationForm                  = "application/x-www-form-urlencoded"
	MIMEApplicationProtobuf              = "application/protobuf"
	MIMEApplicationMsgpack               = "application/msgpack"
	MIMETextHTML                         = "text/html"
	MIMETextHTMLCharsetUTF8              = MIMETextHTML + "; " + charsetUTF8
	MIMETextPlain                        = "text/plain"
	MIMETextPlainCharsetUTF8             = MIMETextPlain + "; " + charsetUTF8
	MIMEMultipartForm                    = "multipart/form-data"
	MIMEOctetStream                      = "application/octet-stream"
)

// Blob Writer
func (c *context) Blob(code int, contentType string, bs []byte) (err error) {
	w := c.Response()
	w.Header().Set("Content-Type", contentType)
	w.WriteHeader(code)
	_, err = w.Write(bs)
	return
}

// HTML Writer
func (c *context) HTML(code int, data string) error {
	return c.Blob(code, MIMETextHTMLCharsetUTF8, []byte(data))
}

// String Writer
func (c *context) String(code int, data string) error {
	return c.Blob(code, MIMETextPlainCharsetUTF8, []byte(data))
}

// JSON Writer
func (c *context) JSON(code int, data interface{}) error {
	bs, err := json.Marshal(data)
	if err != nil {
		return err
	}
	return c.Blob(code, MIMEApplicationJSONCharsetUTF8, bs)
}

// File for Reader
type File interface {
	io.ReadSeeker
	Close() error
	Stat() (os.FileInfo, error)
}

// File Writer
func (c *context) File(f File) (err error) {
	defer f.Close()
	var fi os.FileInfo
	if fi, err = f.Stat(); err != nil {
		return err
	}
	http.ServeContent(c.Response(), c.Request(), fi.Name(), fi.ModTime(), f)
	return
}

// === Reader ===

// JSON Reader
func (c *context) GetJSON(data interface{}) error {
	r := c.Request()
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return err
	}
	if err = json.Unmarshal(body, data); err != nil {
		return err
	}
	return nil
}
