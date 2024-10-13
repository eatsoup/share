package main

import (
	"crypto/md5"
	"errors"
	"fmt"
)

type InMemoryDB []DatabaseItem

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
