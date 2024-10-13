package main

import (
	"github.com/clarkmcc/go-typescript"
)

func Transpile(tsString string) (string, error) {
	transpiled, err := typescript.TranspileString(tsString)
	if err != nil {
		return "", err
	}
	return transpiled, nil
}
