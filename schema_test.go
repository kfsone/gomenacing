package main

import (
	"errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestSchema_Close(t *testing.T) {
	t.Run("Double close check", func(t *testing.T) {
		schema := Schema{}
		assert.Panics(t, func() { failOnError(schema.Close()) })
	})

	testDir := GetTestDir()
	defer testDir.Close()
	database, err := GetDatabase(testDir.Path())
	require.Nil(t, err)
	defer database.Close()

	t.Run("Nominal operation", func(t *testing.T) {
		schema, err := database.GetSchema("nominal")
		require.Nil(t, err)
		require.NotNil(t, schema)
		assert.NotNil(t, schema.store)
		assert.FileExists(t, filepath.Join(schema.db.Path(), "nominal", "main.pix"))
		err = schema.Close()
		assert.Nil(t, err)
		assert.Nil(t, schema.store)
		err = os.RemoveAll(filepath.Join(schema.db.Path(), "nominal"))
		assert.Nil(t, err)
	})

	t.Run("Error handling", func(t *testing.T) {
		schema, err := database.GetSchema("nominal")
		require.Nil(t, err)
		require.NotNil(t, schema)
		store := schema.store
		err = store.Close()
		assert.Nil(t, err)
		err = schema.Close()
		assert.Error(t, err)
		assert.Equal(t, store, schema.store)
	})
}

func TestSchema_PutAndCount(t *testing.T) {
	testDir := GetTestDir()
	defer testDir.Close()

	db, err := GetDatabase(testDir.Path())
	require.Nil(t, err)
	defer db.Close()

	schema, err := db.GetSchema("testing")
	require.Nil(t, err)
	defer func () { failOnError(schema.Close()) }()

	assert.Zero(t, schema.Count())

	assert.Nil(t, schema.Put([]byte("hello"), []byte("world")))
	assert.EqualValues(t, 1, schema.Count())

	assert.Nil(t, schema.Put([]byte("world"), []byte("hello")))
	assert.EqualValues(t, 2, schema.Count())

	assert.Nil(t, schema.Put([]byte("hello"), []byte("hello")))
	assert.EqualValues(t, 2, schema.Count())
}

func TestSchema_LoadData(t *testing.T) {
	testDir := GetTestDir()
	defer testDir.Close()

	db, err := GetDatabase(testDir.Path())
	require.Nil(t, err)
	defer db.Close()

	schema, err := db.GetSchema("schema")
	require.Nil(t, err)

	// With nothing in the database, nothing should happen.
	log := captureLog(t, func(t *testing.T) {
		loaded := 0
		err = schema.LoadData("nothing", func(val []byte) error {
			loaded = 1
			return nil
		})
		assert.Nil(t, err)
		assert.Zero(t, loaded)
	})
	assert.Len(t, log, 1)
	assert.True(t, strings.HasSuffix(log[0], "Loaded 0 nothing."))
	// It should also have closed the schema.
	assert.Panics(t, func() { failOnError(schema.Close()) })

	runTest := func(setupFn func(*Schema)) ([]string, []string, uint32, error) {
		schema, err = db.GetSchema("schema")
		require.Nil(t, err)
		setupFn(schema)

		err = nil
		loaded := make([]string, 0, 4)
		log := captureLog(t, func(t *testing.T) {
			err = schema.LoadData("stuff", func(val []byte) error {
				str := string(val)
				loaded = append(loaded, str)
				if str == "error" {
					return errors.New("got error")
				}
				return nil
			})
			assert.Panics(t, func() { failOnError(schema.Close()) })
		})
		schema, _ = db.GetSchema("schema")
		defer func () { failOnError(schema.Close()) }()
		count := schema.Count()
		return log, loaded, count, err
	}

	log, loaded, count, err := runTest(func(schema *Schema) {
		assert.Nil(t, schema.Put([]byte("hello"), []byte("world")))
		assert.Nil(t, schema.Put([]byte("world"), []byte("hello")))
		assert.Nil(t, schema.Put([]byte("final"), []byte("biscuit")))
	})
	assert.Nil(t, err)
	assert.Equal(t, []string{"world", "hello", "biscuit"}, loaded)
	if assert.Len(t, log, 1) {
		assert.True(t, strings.HasSuffix(log[0], "Loaded 3 stuff."))
	}
	assert.EqualValues(t, 3, count)

	log, loaded, count, err = runTest(func(schema *Schema) {
		assert.Nil(t, schema.Put([]byte("final"), []byte("error")))
	})
	assert.Error(t, err)
	assert.Equal(t, []string{"world", "hello", "error"}, loaded)
	assert.Empty(t, log)
	assert.EqualValues(t, count, 2)
}

func TestSchema_Name(t *testing.T) {
	s := Schema{name: "foible"}
	assert.Equal(t, "foible", s.Name())
}