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

	go func() {
		for {
			select {
			case event, ok := <-watcher.Events:
				if !ok {
					return
				}
				logrus.Debugf("event: %v", event)

				if event.Op&fsnotify.Write == fsnotify.Write {
					logrus.Infof("File Modified: %s", event.Name)
				}
				if event.Op&fsnotify.Create == fsnotify.Create {
					logrus.Infof("File Created: %s", event.Name)
				}
			case err, ok := <-watcher.Errors:
				if !ok {
					return
				}
				logrus.Infof("error:%v", err)
			}
		}
	}()

	logrus.Infof("Registered directory %s with Watcher", dir)
	err = watcher.Add(dir)
	if err != nil {
		logrus.Fatal(err)
	}
}
