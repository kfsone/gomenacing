module github.com/kfsone/gomenacing

go 1.14

require (
	github.com/akrylysov/pogreb v0.9.1
	github.com/golang/protobuf v1.4.2
	github.com/mattn/go-shellwords v1.0.10
	github.com/spf13/pflag v1.0.5
	github.com/stretchr/testify v1.6.1
	github.com/tidwall/gjson v1.6.1
	golang.org/x/sync v0.0.0-20200625203802-6e8e738ad208
	google.golang.org/protobuf v1.25.0
	gopkg.in/zeromq/goczmq.v4 v4.1.0
)

replace github.com/kfsone/gomenacing/pkg/gomschema => ./pkg/gomschema

replace github.com/kfsone/gomenacing/pkg/parsing => ./pkg/parsing

replace github.com/kfsone/gomenacing/pkg/plugins/eddn => ./pkg/plugins/eddn
