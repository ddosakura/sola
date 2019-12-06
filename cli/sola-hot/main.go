package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"path"
	"syscall"
	"time"

	"github.com/fsnotify/fsnotify"
)

var (
	ch = make(chan struct{})
	c  *exec.Cmd
)

func run() {
	cx := exec.Command(
		"go",
		"build",
		"-o",
		"sola-dev",
		".",
	) // TODO
	cx.Stdin = os.Stdin
	cx.Stdout = os.Stdout
	cx.Stderr = os.Stderr
	if e := cx.Run(); e != nil {
		log.Println(e)
	}

	c = exec.Command("./sola-dev") // TODO
	c.Stdin = os.Stdin
	c.Stdout = os.Stdout
	c.Stderr = os.Stderr
	if e := c.Run(); e != nil {
		log.Println(e)
	}
}

func restart() {
	// fmt.Println("up", c.Process.Signal(syscall.SIGHUP))
	fmt.Println("int", c.Process.Signal(syscall.SIGINT))
	fmt.Println("term", c.Process.Signal(syscall.SIGTERM))
	// fmt.Println("kill", c.Process.Kill())
	go run()
}

func throttle() {
	var after <-chan time.Time
	after = time.After(1 * time.Second)
	for {
		select {
		case <-ch:
			return
		case <-after:
			fmt.Println("restart")
			restart()
		}
	}
}

func main() {
	watcher, err := fsnotify.NewWatcher()
	watcher.Add(".") // TODO
	if err != nil {
		log.Fatalln(err)
	}
	defer watcher.Close()

	go run()

	go func() {
		select {
		case <-ch:
			// fmt.Println("first ch")
			return
		}
	}()

	for {
		select {
		case event, ok := <-watcher.Events:
			if !ok {
				return
			}
			eventDispatcher(event)
		case err, ok := <-watcher.Errors:
			if !ok {
				return
			}
			log.Println("error:", err)
		}
	}
}

func eventDispatcher(event fsnotify.Event) {
	ext := path.Ext(event.Name)
	switch event.Op {
	case
		fsnotify.Write,
		fsnotify.Rename:
		if ext == ".go" {
			// log.Println("EVENT", event.Op.String(), event.Name)
			ch <- struct{}{}
			go throttle()
		}
	case fsnotify.Remove:
	case fsnotify.Create:
	}
}
