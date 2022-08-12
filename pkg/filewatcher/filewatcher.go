package filewatcher

import (
	"log"
	"os"
	"strings"

	"github.com/fsnotify/fsnotify"
	"github.com/sirupsen/logrus"
)

type FileWatcher struct {
	FileRegistry []string
	Dir          string
}

func (fw *FileWatcher) Init(dir string) {
	fw.Dir = dir
	rawFiles, err := os.ReadFile(dir + "/files.txt")
	if err != nil {
		logrus.Fatal(err)
	}
	fw.FileRegistry = strings.Split(string(rawFiles), "\n")
}

func contains(s []string, str string) bool {
	for _, v := range s {
		if v == str {
			return true
		}
	}

	return false
}

func (fw *FileWatcher) Start() {
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
					if event.Name == "files/files.txt" {
						logrus.Debugf("File Modified: %s", event.Name)
					} else {
						logrus.Infof("File Modified: %s", event.Name)
					}
				}
				if event.Op&fsnotify.Create == fsnotify.Create {
					logrus.Infof("File Created: %s, adding to registry", event.Name)
					filename := strings.TrimPrefix(event.Name, "files/")
					if !contains(fw.FileRegistry, filename) {
						fw.FileRegistry = append(fw.FileRegistry, filename)
						f, err := os.OpenFile(fw.Dir+"/files.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
						if err != nil {
							log.Println(err)
						}
						if _, err := f.WriteString(filename + "\n"); err != nil {
							log.Println(err)
						}
						f.Close()
					}
				}
				if event.Op&fsnotify.Remove == fsnotify.Remove {
					logrus.Infof("File Deleted: %s", event.Name)
					f, err := os.OpenFile(fw.Dir+"/files.txt", os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0755)
					if err != nil {
						log.Println(err)
					}
					newFilesTxt := ""
					for _, s := range fw.FileRegistry {
						newFilesTxt += s + "\n"
					}
					f.WriteString(newFilesTxt)
					f.Close()
				}
			case err, ok := <-watcher.Errors:
				if !ok {
					return
				}
				logrus.Infof("FileWatcher error:%v", err)
			}
		}
	}()

	logrus.Infof("Registered directory %s with Watcher", fw.Dir)
	err = watcher.Add(fw.Dir)
	if err != nil {
		logrus.Fatal(err)
	}
	<-done
}
