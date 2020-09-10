package main

import (
	"github.com/kfsone/gomenacing/pkg/gomschema"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/proto"
	"os"
	"path/filepath"
	"reflect"
	"testing"
)

func TestOpenDatabase(t *testing.T) {
	testDir := GetTestDir()
	defer testDir.Close()

	t.Run("Nominal functionality", func(t *testing.T) {
		db, err := OpenDatabase(testDir.Path(), "database")
		assert.Nil(t, err)
		assert.NotNil(t, db)
		defer db.Close()
		expectedPath := filepath.Join(testDir.Path(), "database")
		assert.Equal(t, Database{expectedPath}, *db)
		assert.DirExists(t, expectedPath)
	})

	t.Run("Error handling", func(t *testing.T) {
		filePath := filepath.Join(testDir.Path(), "file")
		file, err := os.Create(filePath)
		require.Nil(t, err)
		assert.NotNil(t, file)
		failOnError(file.Close())

		db, err := OpenDatabase(testDir.Path(), "file")
		if assert.Error(t, err) {
			assert.Nil(t, db)
		}
	})
}

func TestDatabase_GetSchema(t *testing.T) {
	testDir := GetTestDir()
	defer testDir.Close()

	db, err := OpenDatabase(testDir.Path(), "database")
	assert.Nil(t, err)
	assert.NotNil(t, db)
	assert.DirExists(t, db.Path())
	defer db.Close()

	schemaPath := filepath.Join(db.Path(), "schema")
	require.NoDirExists(t, schemaPath)
	schema, err := db.GetSchema("schema")
	require.Nil(t, err)
	assert.NotNil(t, schema)
	assert.DirExists(t, schemaPath)
	failOnError(schema.Close())

	// Should be able to open it again.
	schema, err = db.GetSchema("schema")
	require.Nil(t, err)
	assert.NotNil(t, schema)
	failOnError(schema.Close())

	// Now try opening the schema using a filename to cause an error
	schema, err = db.GetSchema(filepath.Join("schema", "main.pix"))
	assert.Error(t, err)
	assert.Nil(t, schema)
}

func TestDatabase_Schemas(t *testing.T) {
	testDir := GetTestDir()
	defer testDir.Close()

	db, err := OpenDatabase(testDir.Path(), "schemas.db")
	assert.Nil(t, err)
	assert.NotNil(t, db)
	assert.DirExists(t, db.Path())
	defer db.Close()

	schemas := []struct {
		name string
		call func() (*Schema, error)
	}{
		{"commodities", func() (*Schema, error) { return db.Commodities() }},
		{"facilities", func() (*Schema, error) { return db.Facilities() }},
		{"listings", func() (*Schema, error) { return db.Listings() }},
		{"systems", func() (*Schema, error) { return db.Systems() }},
	}
	t.Run("Check schemas", func(t *testing.T) {
		for _, schema := range schemas {
			t.Run("- "+schema.name, func(t *testing.T) {
				schemaPath := filepath.Join(db.Path(), schema.name)
				require.NoDirExists(t, schemaPath)
				db, err := schema.call()
				if assert.Nil(t, err) {
					assert.NotNil(t, db)
					assert.DirExists(t, schemaPath)
					failOnError(db.Close())
				}
			})
		}
	})
}

func TestDatabase_Close(t *testing.T) {
	testDir := GetTestDir()
	defer testDir.Close()
	db, err := OpenDatabase(testDir.Path(), "close.db")
	require.Nil(t, err)
	db.Close()
	// Shouldn't have deleted it.
	require.DirExists(t, db.Path())
}

func Test_getSchemaForMessage(t *testing.T) {
	testDir := GetTestDir()
	defer testDir.Close()

	db, err := OpenDatabase(testDir.Path(), "schema.db")
	require.Nil(t, err)
	defer db.Close()

	type args struct {
		db      *Database
		message proto.Message
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{"invalid", args{db, nil}, true},
		{"commodities", args{db, &gomschema.Commodity{}}, false},
		{"systems", args{db, &gomschema.System{}}, false},
		{"facilities", args{db, &gomschema.Facility{}}, false},
		{"listings", args{db, &gomschema.FacilityListing{}}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := getSchemaForMessage(tt.args.db, tt.args.message)
			if got != nil {
				defer func() {
					if err = got.Close(); err != nil {
						panic(err)
					}
				}()
			}
			if (err != nil) != tt.wantErr {
				t.Errorf("getSchemaForMessage() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if err == nil && !reflect.DeepEqual(got.name, tt.name) {
				t.Errorf("getSchemaForMessage() got = %s, want %s", got.name, tt.name)
			}
		})
	}
}
