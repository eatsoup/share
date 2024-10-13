package main

import (
	"io"
	"log"
	"net/http"
	"os"
	"path"
	"strings"
)

const (
	Text  string = "text/plain"
	Image string = "image"
)

var jsBag map[string]string = make(map[string]string)

func main() {
	http.HandleFunc("/upload", uploadHandler)
	http.HandleFunc("/get", getHandler)
	http.HandleFunc("/r/", redirectHandler)
	http.HandleFunc("/", mainHandler)
	http.HandleFunc("/script.js", jsHandler)

	tsList := []string{"script.ts"}
	for _, f := range tsList {
		tsContent, err := os.ReadFile("../frontend/" + f)
		if err != nil {
			log.Fatal(err)
		}
		jsBag[strings.Replace(f, ".ts", ".js", 1)], err = Transpile(string(tsContent))
		if err != nil {
			log.Fatal(err)
		}
	}

	log.Println("Server started on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func mainHandler(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "../frontend/index.html")
}

func jsHandler(w http.ResponseWriter, r *http.Request) {
	content := jsBag[path.Base(r.URL.Path)]
	w.Header().Add("content-type", "text/javascript")
	w.Write([]byte(content))
}

func uploadHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		w.Write([]byte("Method not allowed"))
		return
	}
	p, err := io.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}
	mimeType := http.DetectContentType(p)
	if strings.HasPrefix(mimeType, "text/plain") {
		if strings.HasPrefix(string(p), "http://") || strings.HasPrefix(string(p), "https://") {
			mimeType = "redirect"
			w.Header().Add("Content-Type", mimeType)
		}
	}
	id := DB.Set(p, mimeType)
	w.Write([]byte(id))
}

func getHandler(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	data, err := DB.Get(id)
	if err != nil {
		w.WriteHeader(500)
		return
	}
	w.WriteHeader(200)
	w.Header().Add("Content-Type", data.MimeType)
	w.Write(data.Blob)
}

func redirectHandler(w http.ResponseWriter, r *http.Request) {
	id := path.Base(r.URL.Path)
	data, err := DB.Get(id)
	if err != nil {
		w.WriteHeader(500)
		return
	}
	http.Redirect(w, r, string(data.Blob), http.StatusPermanentRedirect)
}

var DB Database = &InMemoryDB{}

type Database interface {
	Set(blob []byte, mt string) string
	Get(id string) (DatabaseItem, error)
}

type DatabaseItem struct {
	Id       string `json:"id"`
	Blob     []byte `json:"blob"`
	MimeType string `json:"mimeType"`
}
