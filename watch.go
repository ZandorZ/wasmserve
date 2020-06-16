package main

import (
	"fmt"
	"log"
	"time"

	"github.com/donovanhide/eventsource"
	"github.com/fatih/color"
	"github.com/radovskyb/watcher"
)

type timeEvent time.Time

func (t timeEvent) Id() string    { return fmt.Sprint(time.Time(t).UnixNano()) }
func (t timeEvent) Event() string { return "Change" }
func (t timeEvent) Data() string  { return time.Time(t).String() }

func watchFiles() {

	w := watcher.New()
	//
	w.FilterOps(watcher.Write)

	go func() {
		for {
			select {
			case event := <-w.Event:
				c := color.New(color.FgCyan).Add(color.Underline)
				c.Println(event)
				changed <- true
			case err := <-w.Error:
				log.Fatalln(err)
			case <-w.Closed:
				return
			}
		}
	}()

	// Watch this folder for changes.
	if err := w.Add("./"); err != nil {
		log.Fatalln(err)
	}

	// Watch this subfolders for changes.
	if err := w.AddRecursive("./"); err != nil {
		log.Fatalln(err)
	}

	if err := w.Start(time.Millisecond * 100); err != nil {
		log.Fatalln(err)
	}
}

func publisher(srv *eventsource.Server) {
	start := time.Now()
	for {
		<-changed
		srv.Publish([]string{"time"}, timeEvent(start))
		start = start.Add(time.Second)
		// log.Println(start)
	}
}
