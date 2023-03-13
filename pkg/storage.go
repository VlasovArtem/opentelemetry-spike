package main

import (
	"errors"
)

type data struct {
	name string
}

var inmemory = make(map[string]data)

func InsertData(name string) error {
	if _, ok := inmemory[name]; ok {
		err := errors.New("data already exists")
		return err
	}

	inmemory[name] = data{name: name}
	return nil
}
