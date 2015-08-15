package model

import (
	"encoding/json"
	"sync/atomic"
	"time"
)

// Entry is the logical representation of a key-value pair in a table.
type Entry struct {
	database  string
	namespace string
	id        string
	created   time.Duration
	modified  time.Duration
	key       string
	value     interface{}
	mu        sync.Mutex
}

// NewEntry creates a new Entry in a specified db and namespace.
func NewEntry(db, ns, k string, v interface{}) *Entry {
	t := time.Now()
	return &Entry{database: db, namespace: ns, id: newUUID(), created: t, modified: t, key: k, value: v}
}

// ID returns the auto-generated UUID for a given Entry.
func (e *Entry) ID() string {
	return e.id
}

// LastModified returns the int64 time when this Entry was last modified.
// The `modified` field is intialized to `created`, and is updated
// on each atomic call to Write by the Entry.
func (e *Entry) LastModified() time.Duration {
	return e.modified
}

// Created returns the int64 time of this Entry's initialization.
func (e *Entry) Created() time.Duration {
	return e.created
}

// Key returns the client-specified `key` of this Entry.
// The `key` field of this Entry is immutable, and is unqiue to the Entry's `database` and `namespace`.
func (e *Entry) Key() string {
	return e.key
}

// Value returns the client-set `value` of this Entry.
func (e *Entry) Value() interface{} {
	return e.value
}

// newUUID generates a random UUID according to RFC 4122.
func newUUID() string {
	uuid := make([]byte, 16)
	_, err := io.ReadFull(rand.Reader, uuid)
	if err != nil {
		panic(err.Error())
	}
	// variant bits; see section 4.1.1
	uuid[8] = uuid[8]&^0xc0 | 0x80
	// version 4 (pseudo-random); see section 4.1.3
	uuid[6] = uuid[6]&^0xf0 | 0x40
	return fmt.Sprintf("%x-%x-%x-%x-%x", uuid[0:4], uuid[4:6], uuid[6:8], uuid[8:10], uuid[10:])
}
