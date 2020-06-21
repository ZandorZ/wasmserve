package main

import (
	"fmt"
	"log"
	"time"

	"github.com/donovanhide/eventsource"
	"github.com/fatih/color"
	"github.com/radovskyb/watcher"
	"github.com/reactivex/rxgo/v2"
)

// var changed = make(chan bool)
var ch = make(chan rxgo.Item)
var observable rxgo.Observable

var w *watcher.Watcher

type timeEvent time.Time

func (t timeEvent) Id() string    { return fmt.Sprint(time.Time(t).UnixNano()) }
func (t timeEvent) Event() string { return "Change" }
func (t timeEvent) Data() string  { return time.Time(t).String() }

func watchFiles(srv *eventsource.Server) {

	observable = rxgo.FromChannel(ch).
		BufferWithTime(rxgo.WithDuration(1 * time.Second))

	go func() {
		for range observable.Observe() {
			start := time.Now()
			srv.Publish([]string{"time"}, timeEvent(start))
			c := color.New(color.FgHiMagenta).Add(color.Underline)
			c.Println("SSE: ", start.Format(time.RFC850))
		}
	}()

	go func() {
		for {
			select {
			case event := <-w.Event:
				c := color.New(color.FgCyan).Add(color.Underline)
				c.Println(event)
				ch <- rxgo.Of(true)
			case err := <-w.Error:
				log.Fatalln(err)
			case <-w.Closed:
				close(ch)
				return
			}
		}
	}()

	w = watcher.New()

	w.FilterOps(watcher.Write)

	// Watch this subfolders for changes.
	if err := w.AddRecursive("./"); err != nil {
		log.Fatalln(err)
	}

	if err := w.Start(time.Millisecond * 200); err != nil {
		log.Fatalln(err)
	}

}
