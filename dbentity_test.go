package main

import (
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
