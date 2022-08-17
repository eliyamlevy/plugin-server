package server

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/rancher/plugin-server/pkg/filewatcher"
	"github.com/sirupsen/logrus"
)

type FileServer struct {
	Filewatcher *filewatcher.FileWatcher
	Srv         *http.Server
}

func (fs *FileServer) RefreshHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	fs.Filewatcher.Refresh()
}

func HomeHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Hello!\n")
}

func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		//Log file requests
		logrus.Infof("Request made to '%v'", r.RequestURI)
		// Call the next handler, which can be another middleware in the chain, or the final handler.
		next.ServeHTTP(w, r)
	})
}

func (fs *FileServer) Init(dir string, fw *filewatcher.FileWatcher) {
	fs.Filewatcher = fw
	r := mux.NewRouter()

	r.HandleFunc("/", HomeHandler)
	r.HandleFunc("/refresh", fs.RefreshHandler)

	// This will serve files under http://localhost:8000/files/<filename>
	logrus.Infof("Serving files from /%s", dir)
	r.PathPrefix("/files/").Handler(http.StripPrefix("/files/", http.FileServer(http.Dir(dir))))
	r.Use(loggingMiddleware)

	fs.Srv = &http.Server{
		Handler: r,
		Addr:    "0.0.0.0:8000",
		// Good practice: enforce timeouts for servers you create!
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	logrus.Infof("Created FileServer")
}
