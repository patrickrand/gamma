package handler

import (
	"encoding/json"
	"errors"
	log "github.com/patrickrand/gamma/log"
)

var (
	ErrKeyNotFound = errors.New("Key not found")
)

type Parameters map[string]interface{}

func (p Parameters) Get(key string, result interface{}) error {
	log.Debugf("[%s] (Parameters).Get => (%+v).%s, %+v", HANDLER, p, key, result)

	value, ok := p[key]
	if !ok {
		return ErrKeyNotFound
	}

	data, err := json.Marshal(value)
	if err != nil {
		return err
	}

	return json.Unmarshal(data, &result)

}

func (p Parameters) Set(key string, value interface{}) {
	log.Debugf("[%s] (Parameters).Set => (%+v).%s, %+v", HANDLER, p, key, value)
	p[key] = value
}

func (p Parameters) Delete(key string) interface{} {
	log.Debugf("[%s] (Parameters).Delete => (%+v).%s", HANDLER, p, key)

	if value, ok := p[key]; ok {
		delete(p, key)
		return value
	}
	return nil
}
