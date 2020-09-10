package main

import "testing"

func Test_validateEntityForSerialization(t *testing.T) {
	type args struct {
		kind   string
		entity GomDbEntity
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{ "id 0", args{kind:"x", entity: DbEntity{ID:0, DbName:""}}, "", true},
		{ "empty name", args{kind:"x", entity: DbEntity{ID:47, DbName:""}}, "", true},
		{ "spaces name", args{kind:"x", entity: DbEntity{ID:147, DbName:"   \t \t "}}, "", true},
		{ "clean name", args{kind:"x", entity: DbEntity{ID:247, DbName:"Sol system"}}, "Sol system", false},
		{ "padded name", args{kind:"x", entity: DbEntity{ID:347, DbName:"   SOL system \t "}}, "SOL system", false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := validateEntityForSerialization(tt.args.kind, tt.args.entity)
			if (err != nil) != tt.wantErr {
				t.Errorf("validateEntityForSerialization() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("validateEntityForSerialization() got = %v, want %v", got, tt.want)
			}
		})
	}
}
