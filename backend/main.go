package main

import (
	"log"
	"mime"
	"net/http"
	"path/filepath"
)

func main() {
	http.HandleFunc("/generate", uploadHandler)
	// http.HandleFunc("/r/", redirectHandler)
	http.HandleFunc("/", mainHandler)

	mime.AddExtensionType(".ts", "application/typescript")
	log.Println("Server started on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func mainHandler(w http.ResponseWriter, r *http.Request) {
	path := "../frontend" + r.URL.Path
	ext := filepath.Ext(path)
	if ext == ".ts" {
		w.Header().Set("Content-Type", "application/typescript")
	}
	http.ServeFile(w, r, path)
}

func uploadHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		w.Write([]byte("Method not allowed"))
		return
	}
	var p []byte
	_, err := r.Body.Read(p)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}
	w.Write(p)
}
