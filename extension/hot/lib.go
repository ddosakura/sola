package hot

import (
	"errors"
	"log"
	"os"
	"os/exec"
	"path"
	"plugin"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/ddosakura/sola/v2"
	"github.com/fsnotify/fsnotify"
)

// Hot Plugin
type Hot struct {
	option  *Option
	watcher *fsnotify.Watcher

	lock        sync.Mutex
	timers      map[string]*timer
	lh          sync.RWMutex
	handlers    map[string]sola.Handler
	lm          sync.RWMutex
	middlewares map[string]sola.Middleware
}

// Option of Hot Plugin
type Option struct {
	TmpDir string
	Init   []string
	Watch  []string
	Delay  time.Duration
}

// error(s)
var (
	ErrTmpDir = errors.New("TmpDir error")
)

// New Hot Plugin
func New(o *Option) (*Hot, error) {
	if o == nil {
		o = &Option{}
	}
	if o.TmpDir == "" {
		o.TmpDir = "./tmp"
	}
	if fi, err := os.Stat(o.TmpDir); err == nil {
		if !fi.IsDir() {
			return nil, ErrTmpDir
		}
	} else {
		if os.IsNotExist(err) {
			if e := os.Mkdir(o.TmpDir, os.ModePerm); e != nil {
				return nil, e
			}
		} else {
			return nil, ErrTmpDir
		}
	}
	if o.Init == nil {
		o.Init = []string{}
	}
	if o.Watch == nil {
		o.Watch = []string{}
	}
	if o.Delay <= 0 {
		o.Delay = time.Second
	}

	watcher, err := fsnotify.NewWatcher()
	for _, p := range o.Watch {
		watcher.Add(p)
	}
	if err != nil {
		return nil, err
	}

	h := &Hot{
		option:      o,
		watcher:     watcher,
		timers:      make(map[string]*timer),
		handlers:    make(map[string]sola.Handler),
		middlewares: make(map[string]sola.Middleware),
	}
	for _, p := range o.Init {
		load(h, p)
	}

	return h, nil
}

// Used by Sola App
func (h *Hot) Used(app *sola.Sola) {
	app.Use(sola.Handler(func(c sola.Context) error {
		c.Set(CtxHot, h)
		return nil
	}).M())
}

// Modules from context
func Modules(c sola.Context) *Hot {
	return c.Get(CtxHot).(*Hot)
}

// Scan Hot Plugin
func (h *Hot) Scan() {
	defer h.watcher.Close()
	for {
		select {
		case event, ok := <-h.watcher.Events:
			if !ok {
				os.Exit(-1)
			}
			h.eventDispatcher(event)
		case err, ok := <-h.watcher.Errors:
			if !ok {
				log.Fatalln(err)
			}
		}
	}
}

func (h *Hot) eventDispatcher(event fsnotify.Event) {
	ext := path.Ext(event.Name)
	dir := path.Dir(event.Name)
	switch event.Op {
	case
		fsnotify.Write,
		fsnotify.Rename:
		if ext == ".go" {
			log.Println("EVENT", event.Op.String(), event.Name)
			h.lock.Lock()
			t := h.timers[dir]
			if t == nil {
				t = &timer{h, dir, make(chan struct{}), h.option.Delay}
				h.timers[dir] = t
				go t.readch()
			}
			h.lock.Unlock()
			t.ch <- struct{}{}
			go t.throttle()
		}
	case fsnotify.Remove:
	case fsnotify.Create:
	}
}

func (h *Hot) loadHandler(k string, v sola.Handler) {
	h.handlers[k] = v
}

func (h *Hot) loadMiddleware(k string, v sola.Middleware) {
	h.middlewares[k] = v
}

// Handler Getter
func (h *Hot) Handler(k string) sola.Handler {
	return func(c sola.Context) error {
		h.lh.RLock()
		defer h.lh.RUnlock()
		log.Println("use handler:", k)
		if x := h.handlers[k]; x != nil {
			return x(c)
		}
		return nil
	}
}

// Middleware Getter
func (h *Hot) Middleware(k string) sola.Middleware {
	return sola.M(func(c sola.C, next sola.H) error {
		h.lm.RLock()
		defer h.lm.RUnlock()
		log.Println("use middleware:", k)
		if x := h.middlewares[k]; x != nil {
			return x(next)(c)
		}
		return next(c)
	}).Must(NotFound(k))
}

type timer struct {
	h     *Hot
	dir   string
	ch    chan struct{}
	delay time.Duration
}

func (t *timer) throttle() {
	// fmt.Println("do throttle")
	var after <-chan time.Time
	after = time.After(t.delay)
	select {
	case <-t.ch:
		// fmt.Println("cancel ch")
		return
	case <-after:
		log.Println("update plugins")
		t.update()
		go t.readch()
		return
	}
}

func (t *timer) readch() {
	// fmt.Println("do readch")
	select {
	case <-t.ch:
		// fmt.Println("first ch")
		return
	}
}

func (t *timer) update() {
	dir := t.dir
	if !path.IsAbs(dir) {
		dir = "./" + dir
	}
	load(t.h, dir)
}

func load(h *Hot, dir string) {
	log.Println("load", dir)
	var p *plugin.Plugin
	if strings.HasSuffix(dir, ".so") {
		var e error
		if p, e = plugin.Open(dir); e != nil {
			log.Println(e)
			return
		}
	} else {
		hash := strconv.Itoa(int(time.Now().Unix()))
		path := h.option.TmpDir + "/" + hash
		var e error
		if e = run("cp", "-r", dir, path); e != nil {
			log.Println(e)
			return
		}
		dir = path
		path += "/plugin.so"
		if e = run(
			"go",
			"build",
			"-buildmode=plugin",
			"-o",
			path,
			dir,
		); e != nil {
			log.Println(e)
			return
		}
		if p, e = plugin.Open(path); e != nil {
			log.Println(e)
			return
		}
		if e = run("rm", "-r", dir); e != nil {
			log.Println(e)
			return
		}
	}

	go func() {
		if exportH, e := p.Lookup("ExportHandler"); e != nil {
			log.Println(e)
		} else {
			h.lh.Lock()
			defer h.lh.Unlock()
			exportHandler := *exportH.(*map[string]sola.Handler)
			for k := range exportHandler {
				log.Println("load handler:", k)
				h.loadHandler(k, exportHandler[k])
			}
		}
	}()

	go func() {
		if exportM, e := p.Lookup("ExportMiddleware"); e != nil {
			log.Println(e)
		} else {
			h.lm.Lock()
			defer h.lm.Unlock()
			exportMiddleware := *exportM.(*map[string]sola.Middleware)
			for k := range exportMiddleware {
				log.Println("load middleware:", k)
				h.loadMiddleware(k, exportMiddleware[k])
			}
		}
	}()
}

func run(name string, args ...string) error {
	c := exec.Command(name, args...)
	// fmt.Println(c.String())
	c.Stdin = os.Stdin
	c.Stdout = os.Stdout
	c.Stderr = os.Stderr
	return c.Run()
}
