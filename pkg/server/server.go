package server

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/mux"
)

func HomeHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Hello World\n")
}

func FilesHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	w.WriteHeader(http.StatusOK)

	filename := "files/" + string(vars["filename"])

	content, err := os.ReadFile(filename)
	if err != nil {
		log.Print(err)
		content = []byte("File '" + filename + "' not found.\n")
	}

	fmt.Fprintf(w, "%s\n", string(content))
}

func FileListHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)

	content, err := os.ReadFile("files/files.txt")
	if err != nil {
		log.Fatal(err)
	}

	fmt.Fprintf(w, "%s\n", string(content))
}

func New(dir string) *http.Server {
	r := mux.NewRouter()
	r.HandleFunc("/", HomeHandler)
	// r.HandleFunc("/files/{filename}", FilesHandler)
	r.HandleFunc("/files.txt", FileListHandler)
	// This will serve files under http://localhost:8000/static/<filename>
	log.Printf("Serving files from %s\n", dir)
	r.PathPrefix("/files/").Handler(http.StripPrefix("/files/", http.FileServer(http.Dir(dir))))

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
