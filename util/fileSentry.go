package util

import (
	"github.com/fsnotify/fsnotify"
	"log"
	"os"
	"path/filepath"
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

	signal := make(chan Event)
	go func() {
		for {
			select {
			case event := <-watcher.Events:
				switch {
				case event.Op&fsnotify.Write == fsnotify.Write:
					signal <- Event{event.Name, fsnotify.Write}

				case event.Op&fsnotify.Create == fsnotify.Create:
					if info, err := os.Stat(event.Name); err == nil && info.IsDir() {
						watcher.Add(event.Name)
					} else {
						signal <- Event{event.Name, fsnotify.Create}
					}

				case event.Op&fsnotify.Remove == fsnotify.Remove:
					watcher.Remove(event.Name)
				}

			case err := <-watcher.Errors:
				log.Println("error:", err)
			}
		}
	}()

	filepath.Walk(path, func(p string, info os.FileInfo, err error) error {
		if info != nil && info.IsDir() {
			watcher.Add(path)
		}
		return nil
	})

	if err = watcher.Add(path); err != nil {
		log.Fatal(err)
	}

	return watcher, signal
}
