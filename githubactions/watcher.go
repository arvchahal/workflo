package githubactions

import (
	"log"

	"github.com/fsnotify/fsnotify"
)

// WatchAndUpdate watches for changes in the specified path and regenerates YAML
func WatchAndUpdate(filename string, wf *Workflow) error {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return err
	}
	defer watcher.Close()

	go func() {
		for {
			select {
			case event := <-watcher.Events:
				if event.Op&fsnotify.Write == fsnotify.Write {
					log.Println("File changed:", event.Name)
					wf.GenerateYAML(filename, true)
				}
			case err := <-watcher.Errors:
				log.Println("Watcher error:", err)
			}
		}
	}()

	return watcher.Add("path/to/watch")
}
