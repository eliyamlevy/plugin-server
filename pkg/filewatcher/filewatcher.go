package filewatcher

import (
	"errors"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"github.com/fsnotify/fsnotify"
	"github.com/sirupsen/logrus"
)

type FileWatcher struct {
	FileRegistry []string
	Dir          string
}

func (fw *FileWatcher) Refresh() {
	fw.FileRegistry = nil
	walkFunc := func(path string, d fs.DirEntry, err error) error {
		//Check if current path is a directory
		if d.IsDir() {
			return nil
		}

		//If it's a file append to list
		fw.FileRegistry = append(fw.FileRegistry, path)
		logrus.Debugf("Adding file, %s, to registry", path)
		return nil
	}

	if err := filepath.WalkDir(fw.Dir, walkFunc); err != nil {
		logrus.Fatal(err)
	}

	if !contains(fw.FileRegistry, "files/files.txt") {
		fw.FileRegistry = append(fw.FileRegistry, "files/files.txt")
	}
}

func (fw *FileWatcher) Init(dir string) {
	fw.Dir = dir
	if _, err := os.Stat(fw.Dir + "/files.txt"); errors.Is(err, os.ErrNotExist) {
		logrus.Debugf("No files.txt found in %s, creating new files.txt.", fw.Dir)
		fw.Refresh()
		f, err := os.OpenFile(fw.Dir+"/files.txt", os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0755)
		if err != nil {
			logrus.Fatal(err)
		}

		newFilesTxt := ""
		for _, s := range fw.FileRegistry {
			newFilesTxt += s + "\n"
		}
		f.WriteString(newFilesTxt)
		f.Close()

	} else {
		rawFiles, err := os.ReadFile(dir + "/files.txt")
		if err != nil {
			logrus.Fatal(err)
		}
		fw.FileRegistry = strings.Split(string(rawFiles), "\n")
	}
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
					if fileInfo, _ := os.Stat(event.Name); fileInfo.IsDir() {
						if err != nil {
							logrus.Fatal(err)
						}
					} else {
						logrus.Infof("File Created: %s, adding to registry", event.Name)
						fw.Refresh()
						f, err := os.OpenFile(fw.Dir+"/files.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
						if err != nil {
							logrus.Fatal(err)
						}
						if _, err := f.WriteString(event.Name + "\n"); err != nil {
							logrus.Fatal(err)
						}
						f.Close()
						logrus.Infof("fw: %v", fw.FileRegistry)
					}

				}
				if event.Op&fsnotify.Remove == fsnotify.Remove {
					if contains(fw.FileRegistry, event.Name) {
						logrus.Infof("File or Directory Deleted: %s", event.Name)
						fw.Refresh()
						f, err := os.OpenFile(fw.Dir+"/files.txt", os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0755)
						if err != nil {
							logrus.Fatal(err)
						}
						newFilesTxt := ""
						for _, s := range fw.FileRegistry {
							newFilesTxt += s + "\n"
						}
						f.WriteString(newFilesTxt)
						f.Close()
						logrus.Infof("fw: %v", fw.FileRegistry)
					} else {
						logrus.Debugf("Remove Event: %s", event.Name)
					}
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
