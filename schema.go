package main

// Schema represents one of the tables in the database, and plays a major
// role in abstracting away database operations.

import (
	"github.com/akrylysov/pogreb"
)

type Schema struct {
	db    *Database
	name  string
	store *pogreb.DB
}

func (s *Schema) Close() error {
	if s.store == nil {
		panic("double close")
	}
	if err := s.store.Close(); err != nil {
		return err
	}
	s.store = nil
	return nil
}

func (s Schema) Name() string {
	return s.name
}

func (s Schema) Count() uint32 {
	return s.store.Count()
}

func (s *Schema) Put(key []byte, value []byte) error {
	return s.store.Put(key, value)
}

func (s *Schema) LoadData(loader *DataLoader) (int, error) {
	defer func() { failOnError(s.Close()) }()

	it := s.store.Items()
	loaded := 0
	for {
		key, val, err := it.Next()
		if err == nil {
			err = loader.Load(val)
		}
		if err != nil {
			if err == pogreb.ErrIterationDone {
				break
			}
			if FilterError(err) != nil {
				failOnError(s.store.Delete(key))
				return loaded, err
			}
		}
		loaded++
	}
	return loaded, nil
}
