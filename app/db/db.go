package db

import "sync"

type Value struct {
	Val string
	Exp int64
}

type DB struct {
	mu  *sync.RWMutex
	kvs map[string]Value
}

func New() *DB {
	return &DB{
		mu:  &sync.RWMutex{},
		kvs: make(map[string]Value),
	}
}

func (d *DB) Set(key string, val Value) {
	d.mu.Lock()
	defer d.mu.Unlock()

	d.kvs[key] = val
}

func (d *DB) Get(key string) (Value, bool) {
	d.mu.RLock()
	defer d.mu.RUnlock()
	val, ok := d.kvs[key]

	return val, ok
}
