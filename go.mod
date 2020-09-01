module github.com/kfsone/gomenacing

go 1.14

require (
	github.com/akrylysov/pogreb v0.9.1
	github.com/kfsone/gomenacing/pkg/gomschema v0.0.0
	github.com/mattn/go-shellwords v1.0.10
	github.com/spf13/pflag v1.0.5
	github.com/stretchr/testify v1.6.1
	golang.org/x/sync v0.0.0-20190423024810-112230192c58
)

replace github.com/kfsone/gomenacing/pkg/gomschema => ./pkg/gomschema
