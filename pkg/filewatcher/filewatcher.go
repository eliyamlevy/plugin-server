package filewatcher

import (
	"github.com/fsnotify/fsnotify"
	"github.com/sirupsen/logrus"
)

// type FileWatcher struct {
// 	FileRegistry []string
// }

// func (fw *FileWatcher) New() *FileWatcher {
// 	fw.FileRegistry

// 	return fw
// }

func Start(dir string) {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		logrus.Fatal(err)
	}
	defer watcher.Close()

	done := make(chan bool)
	go func() {
		for {
			select {
			case event, ok := <-watcher.Events:
				if !ok {
					return
				}
				logrus.Debugf("event: %v", event)

				if event.Op&fsnotify.Write == fsnotify.Write {
					logrus.Infof("modified file:%s", event.Name)
				}
			case err, ok := <-watcher.Errors:
				if !ok {
					return
				}
				logrus.Infof("error:%v", err)
			}
		}
	}()

	err = watcher.Add(dir)
	if err != nil {
		logrus.Fatal(err)
	}
	<-done
}
