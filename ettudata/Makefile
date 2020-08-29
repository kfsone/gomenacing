# This is an awful but obvious place to record how to generate the pb.go file.
#
datatypes.pb.go: datatypes.proto
	protoc --go_out=paths=source_relative:./ -I. datatypes.proto
	protoc --python_out=./ -I. datatypes.proto
