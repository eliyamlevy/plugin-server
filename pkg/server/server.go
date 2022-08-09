package server

import (
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
)

func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Do stuff here
		logrus.Infof("Serving file '%v'", r.RequestURI)
		// Call the next handler, which can be another middleware in the chain, or the final handler.
		next.ServeHTTP(w, r)
	})
}

func New(dir string) *http.Server {
	r := mux.NewRouter()

	// This will serve files under http://localhost:8000/files/<filename>
	logrus.Infof("Serving files from %s\n", dir)
	// fsys := dotFileHidingFileSystem{http.Dir(dir)}

	r.PathPrefix("/files/").Handler(http.StripPrefix("/files/", http.FileServer(http.Dir(dir))))
	r.Use(loggingMiddleware)

	http.Handle("/", r)

	srv := &http.Server{
		Handler: r,
		Addr:    "127.0.0.1:8000",
		// Good practice: enforce timeouts for servers you create!
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	return srv
}
