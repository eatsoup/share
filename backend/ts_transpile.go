package main

import (
	"github.com/clarkmcc/go-typescript"
)

// TODO: Load all on startup in some kind of cache
func Transpile(tsString string) (string, error) {
	transpiled, err := typescript.TranspileString(tsString)
	if err != nil {
		return "", err
	}
	return transpiled, nil
}
