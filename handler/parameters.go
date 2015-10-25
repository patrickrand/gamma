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
	log.DBUG("parameters", "(Parameters).Get => (%+v).%s, %+v", p, key, result)

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
	log.DBUG("parameters", "(Parameters).Set => (%+v).%s, %+v", p, key, value)

	p[key] = value
}

func (p Parameters) Delete(key string) interface{} {
	log.DBUG("parameters", "(Parameters).Delete => (%+v).%s", p, key)

	if value, ok := p[key]; ok {
		delete(p, key)
		return value
	}
	return nil
}
