package agent

import (
	"encoding/json"
	"errors"
)

var (
	ErrKeyNotFound = errors.New("key not found")
)

type Parameters map[string]interface{}

func (p Parameters) Get(key string, result interface{}) error {
	log.Debugf("(Parameters).Get => (%+v).%s, %+v", p, key, result)

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
	log.Debugf("(Parameters).Set => (%+v).%s, %+v", p, key, value)
	p[key] = value
}

func (p Parameters) Delete(key string) interface{} {
	log.Debugf("(Parameters).Delete => (%+v).%s", p, key)

	value, ok := p[key]
	if !ok {
		return nil
	}

	delete(p, key)
	return value
}
