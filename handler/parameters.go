package handler

import (
	"encoding/json"
	"errors"
)

var (
	ErrKeyNotFound = errors.New("Key not found")
)

type Parameters map[string]interface{}

func (p Parameters) Get(key string, result interface{}) error {
	if value, ok := p[key]; !ok {
		return ErrKeyNotFound
	} else {
		data, err := json.Marshal(value)
		if err != nil {
			return err
		}

		return json.Unmarshal(data, &result)
	}
}

func (p Parameters) Set(key string, value interface{}) {
	p[key] = value
}

func (p Parameters) Delete(key string) interface{} {
	if value, ok := p[key]; ok {
		delete(p, key)
		return value
	}
	return nil
}
