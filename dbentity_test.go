package main

import (
	"github.com/stretchr/testify/assert"
	"reflect"
	"testing"
)

func TestNewDbEntity(t *testing.T) {
	type args struct {
		id   int64
		name string
	}
	tests := []struct {
		name       string
		args       args
		wantEntity DbEntity
		wantErr    bool
	}{
		{"zero id", args{0, ""}, DbEntity{}, true},
		{"negative id", args{-1, ""}, DbEntity{}, true},
		{"overflow id", args{1 << 32, ""}, DbEntity{}, true},
		{"empty name", args{1, ""}, DbEntity{}, true},
		{"bad name", args{1, " \t "}, DbEntity{}, true},
		{"valid", args{3, "Foods"}, DbEntity{EntityID(3), "Foods"}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotEntity, err := NewDbEntity(tt.args.id, tt.args.name)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewDbEntity() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotEntity, tt.wantEntity) {
				t.Errorf("NewDbEntity() gotEntity = %v, want %v", gotEntity, tt.wantEntity)
			}
		})
	}
}

func TestDbEntity_GetId(t *testing.T) {
	e := DbEntity{}
	assert.Equal(t, uint32(0), e.GetId())

	e = DbEntity{ID: 1134, DbName: "Something"}
	assert.Equal(t, uint32(1134), e.GetId())
}

func TestDbEntity_GetName(t *testing.T) {
	e := DbEntity{}
	assert.Equal(t, "", e.GetName())
	e.DbName = "chicken & Biscuits"
	assert.Equal(t, "chicken & Biscuits", e.GetName())
}

func Test_validateEntity(t *testing.T) {
	type args struct {
		id   int64
		name string
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{"zero value", args{0, ""}, "", true},
		{ "empty name", args{1, ""}, "", true},
		{ "whitespace name", args{1, " \t"}, "", true},
		{ "invalid name", args{1, "a"}, "", true},
		{ "invalid name w/space", args{1, " a "}, "", true},
		{ "excessive id", args{int64(1) << 34, "Sol"}, "", true},
		{ "1, Sol", args{1, "Sol"}, "Sol", false},
		{ "1, Sol w/spaces", args{1, " \t Sol \t "}, "Sol", false},
		{ "~0, sOl sYsTeM w/spaces", args{(int64(1) << 32) - 1, " \t sOl sYsTeM\t "}, "sOl sYsTeM", false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := validateEntity(tt.args.id, tt.args.name)
			if (err != nil) != tt.wantErr {
				t.Errorf("validateEntity() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("validateEntity() got = %v, want %v", got, tt.want)
			}
		})
	}
}