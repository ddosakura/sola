package router

import (
	"strings"

	"github.com/ddosakura/sola/v2"
)

// Ctx
const (
	CtxMeta        = "xrouter.meta"
	CtxParamPrefix = "xrouter.param."
)

// Meta of Route
type Meta struct {
	Ctx     sola.Context
	Method  string
	Host    string
	Path    []string
	Default bool
}

func newMeta(c sola.Context) *Meta {
	tmp := c.Get(CtxMeta)
	if tmp != nil {
		return tmp.(*Meta)
	}
	r := c.Request()
	m := &Meta{
		Ctx:    c,
		Method: r.Method,
		Host:   r.Host,
		Path:   strings.Split(r.URL.EscapedPath(), "/")[1:],
	}
	c.Set(CtxMeta, m)
	return m
}

func buildMeta(pattern string) *Meta {
	matches := strings.Split(pattern, " ")
	method := ""
	url := matches[0]
	if len(matches) > 1 {
		method = matches[0]
		url = matches[1]
	}
	def := false
	if strings.HasSuffix(url, "/*") {
		def = true
	}
	urls := strings.Split(url, "/")
	if def {
		urls = urls[:len(urls)-1]
	}

	return &Meta{
		Method:  strings.ToUpper(method),
		Host:    urls[0],
		Path:    urls[1:],
		Default: def,
	}
}

// route match pattern
func match(c sola.Context, strict bool, pattern *Meta) func() {
	m := newMeta(c)
	if pattern.Method != "" && pattern.Method != m.Method {
		return nil
	}
	if pattern.Host != "" && pattern.Host != m.Host {
		return nil
	}
	pLen := len(pattern.Path)
	mLen := len(m.Path)
	if mLen < pLen {
		return nil
	}
	if strict && mLen > pLen && !pattern.Default {
		return nil
	}

	params := map[string]string{}
	for i, path := range pattern.Path {
		v := m.Path[i]
		if path == v {
			continue
		}
		if !strings.HasPrefix(path, ":") {
			return nil
		}
		k := strings.Replace(path, ":", "", 1)
		params[k] = v
	}
	for k, v := range params {
		c.Set(CtxParam(k), v)
	}

	return func() {
		mx := &Meta{}
		mx.Ctx = m.Ctx
		mx.Method = m.Method
		mx.Host = m.Host
		mx.Path = m.Path[pLen:]
		c.Set(CtxMeta, mx)
	}
}

// CtxParam Builder
func CtxParam(key string) string {
	return CtxParamPrefix + key
}

// Param in route
func Param(c sola.Context, key string) string {
	return c.Get(CtxParam(key)).(string)
}
