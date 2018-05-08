package util

import (
	"github.com/fsnotify/fsnotify"
	"log"
	"net/http"
	"os"
)

type Event struct {
	Name string
	Op   fsnotify.Op
}

func WatchFile(path string) (*fsnotify.Watcher, chan Event) {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatal(err)
	}
	//defer watcher.Close()

	signal := make(chan Event)
	go func() {
		for {
			select {
			case event := <-watcher.Events:
				switch {
				case event.Op&fsnotify.Write == fsnotify.Write:
					signal <- Event{event.Name, fsnotify.Write}
				case event.Op&fsnotify.Create == fsnotify.Create:
					info, e := os.Stat(event.Name)
					if e == nil && info.IsDir() {
						watcher.Add(event.Name)
					}
					signal <- Event{event.Name, fsnotify.Create}
				case event.Op&fsnotify.Remove == fsnotify.Remove:
					info, e := os.Stat(event.Name)
					if e == nil && info.IsDir() {
						watcher.Remove(event.Name)
					}
					signal <- Event{event.Name, fsnotify.Remove}
				}
			case err := <-watcher.Errors:
				log.Println("error:", err)
			}
		}
	}()

	err = watcher.Add(path)
	if err != nil {
		log.Fatal(err)
	}

	return watcher, signal
}

func main() {
	watcher, signal := WatchFile("/Users/liolay/config-repo")
	go func() {
		for {
			a := <-signal
			log.Println("=========",a.Name)
		}
	}()
	log.Fatal(http.ListenAndServe(":8080", nil))
	watcher.Close()
}
