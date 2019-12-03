package rest

import (
	"errors"
	"net/http"
	"strconv"
	"strings"

	"github.com/ddosakura/sola/v2"
	"github.com/ddosakura/sola/v2/middleware/router"
)

// Option of RESTful Middleware
type Option struct {
	NewModel func() interface{}

	// Optional
	DefaultPageSize int
	Path            string
	GetFunc         func(id string) interface{}
	ListFunc        func(page, size int) interface{}
	PostFunc        func(interface{}) error
	PutFunc         func(interface{}) error
	DeleteFunc      func(id string) error
}

// error(s)
var (
	ErrOption = errors.New("Must set NewModel when use PostFunc/PutFunc")
)

// Response Common
type Response struct {
	Code int         `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data"`
}

// ErrResponse Common
var ErrResponse = &Response{
	Code: -1,
	Msg:  "FAIL",
}

func success(c sola.Context, v interface{}) error {
	return c.JSON(http.StatusOK, &Response{
		Code: 0,
		Msg:  "SUCCESS",
		Data: v,
	})
}

func fail(c sola.Context, msg ...string) error {
	if msg == nil || len(msg) == 0 {
		return c.JSON(http.StatusOK, ErrResponse)
	}
	return c.JSON(http.StatusOK, &Response{
		Code: -1,
		Msg:  strings.Join(msg, "; "),
	})
}

// New RESTful Router
func New(o *Option) *router.Router {
	if o == nil {
		o = &Option{}
	}
	if o.NewModel == nil &&
		(o.PostFunc != nil || o.PutFunc != nil) {
		panic(ErrOption)
	}
	if o.DefaultPageSize < 1 {
		o.DefaultPageSize = 5
	}

	r := router.New()

	if o.GetFunc != nil {
		r.BindFunc("GET "+o.Path+"/:id", func(c sola.Context) error {
			id := router.Param(c, "id")
			if model := o.GetFunc(id); model != nil {
				return success(c, model)
			}
			return fail(c)
		})
	}
	if o.ListFunc != nil {
		r.BindFunc("GET "+o.Path, func(c sola.Context) error {
			r := c.Request()
			qs := r.URL.Query()
			page, err := strconv.Atoi(qs.Get("page"))
			if err != nil || page < 1 {
				page = 1
			}
			size, err := strconv.Atoi(qs.Get("size"))
			if err != nil || size < 1 {
				size = o.DefaultPageSize
			}
			if models := o.ListFunc(page, size); models != nil {
				return success(c, models)
			}
			return fail(c)
		})
	}

	if o.PostFunc != nil {
		r.BindFunc("POST "+o.Path, func(c sola.Context) error {
			model := o.NewModel()
			if err := c.GetJSON(model); err != nil {
				return err
			}
			if err := o.PostFunc(model); err != nil {
				return fail(c, err.Error())
			}
			return success(c, nil)
		})
	}
	if o.PutFunc != nil {
		r.BindFunc("PUT "+o.Path, func(c sola.Context) error {
			model := o.NewModel()
			if err := c.GetJSON(model); err != nil {
				return err
			}
			if err := o.PutFunc(model); err != nil {
				return fail(c, err.Error())
			}
			return success(c, nil)
		})
	}
	if o.DeleteFunc != nil {
		r.BindFunc("DELETE "+o.Path+"/:id", func(c sola.Context) error {
			id := router.Param(c, "id")
			if err := o.DeleteFunc(id); err != nil {
				return fail(c, err.Error())
			}
			return success(c, nil)
		})
	}

	return r
}
