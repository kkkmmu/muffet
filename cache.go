package main

import (
	"github.com/dgraph-io/badger"
	"github.com/dgraph-io/badger/options"
)

type cache struct {
	database *badger.DB
}

func newCache(d string) (cache, error) {
	o := badger.DefaultOptions
	o.Dir = d
	o.ValueDir = d
	o.ValueLogLoadingMode = options.FileIO

	db, err := badger.Open(o)

	if err != nil {
		return cache{}, err
	}

	return cache{db}, nil
}

func (c cache) Close() {
	c.database.Close()
}

func (c cache) Add(u string, x interface{}) error {
	c.database.Update(func(t *badger.Txn) error {
		return t.Set([]byte(u), encodeResult(x))
	})

	return nil
}

func (c cache) Get(u string) (interface{}, error) {
	bs := []byte(nil)

	if err := c.database.View(func(t *badger.Txn) error {
		i, err := t.Get([]byte(u))

		if err != nil {
			return err
		}

		bs, err = i.Value()

		return err
	}); err == badger.ErrKeyNotFound {
		return nil, nil
	} else if err != nil {
		return nil, err
	}

	return decodeResult(bs), nil
}

func encodeResult(x interface{}) []byte {
	switch x := x.(type) {
	case error:
		return []byte(x.Error())
	case fetchResult:
		return r.Encode()
	}

	panic("unreachable")
}

func decodeResult(bs []byte) interface{} {
	switch x := x.(type) {
	case error:
		return []byte(x.Error())
	case fetchResult:

	}

	panic("unreachable")
}
