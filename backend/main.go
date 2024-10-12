package main

import (
	"crypto/md5"
	"errors"
	"fmt"
	"io"
	"log"
	"mime"
	"net/http"
	"os"
	"path/filepath"
)

const (
	Text  string = "text/plain"
	Image string = "image"
)

func main() {
	http.HandleFunc("/upload", uploadHandler)
	http.HandleFunc("/get", getHandler)
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
		tsBytes, err := os.ReadFile(path)
		if err != nil {
			w.WriteHeader(404)
			return
		}
		js, err := Transpile(string(tsBytes))
		if err != nil {
			w.WriteHeader(500)
			return
		}
		w.Header().Set("Content-Type", "application/javascript")
		w.Write([]byte(js))
		return
	}
	http.ServeFile(w, r, path)
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

type InMemoryDB []DatabaseItem

var DB InMemoryDB

func (db *InMemoryDB) Set(blob []byte, mt string) string {
	id := fmt.Sprintf("%x", md5.Sum(blob))
	item := DatabaseItem{
		Id:       id,
		Blob:     blob,
		MimeType: mt,
	}
	*db = append(*db, item)
	return id
}

func (db *InMemoryDB) Get(id string) (DatabaseItem, error) {
	for _, entry := range *db {
		if entry.Id == id {
			return entry, nil
		}
	}
	return DatabaseItem{}, errors.New("not found")
}

type DatabaseItem struct {
	Id       string `json:"id"`
	Blob     []byte `json:"blob"`
	MimeType string `json:"mimeType"`
}
