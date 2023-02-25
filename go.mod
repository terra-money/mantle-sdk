module github.com/terra-project/mantle-sdk

go 1.15

require (
	github.com/cosmos/cosmos-sdk v0.39.1
	github.com/fsnotify/fsnotify v1.4.9 // indirect
	github.com/gertd/go-pluralize v0.1.7
	github.com/go-openapi/errors v0.19.6
	github.com/go-openapi/runtime v0.19.19
	github.com/go-openapi/strfmt v0.19.5
	github.com/go-openapi/swag v0.19.9
	github.com/go-openapi/validate v0.19.10
	github.com/golang/protobuf v1.4.2 // indirect
	github.com/google/uuid v1.1.2
	github.com/gorilla/mux v1.7.4
	github.com/gorilla/websocket v1.4.2
	github.com/graphql-go/graphql v0.7.9
	github.com/graphql-go/handler v0.2.3
	github.com/mitchellh/mapstructure v1.4.1
	github.com/onsi/ginkgo v1.10.1 // indirect
	github.com/onsi/gomega v1.7.0 // indirect
	github.com/pelletier/go-toml v1.8.0 // indirect
	github.com/rs/cors v1.7.0
	github.com/spf13/afero v1.2.2 // indirect
	github.com/spf13/cast v1.3.1 // indirect
	github.com/spf13/viper v1.7.0
	github.com/stretchr/testify v1.6.1
	github.com/syndtr/goleveldb v1.0.1-0.20190923125748-758128399b1d
	github.com/tendermint/tendermint v0.33.7
	github.com/tendermint/tm-db v0.5.2
	github.com/terra-project/core v0.4.1
	github.com/terra-project/mantle-compatibility v1.6.0-columbus-4
	github.com/vmihailenco/msgpack/v5 v5.0.0-beta.1
	golang.org/x/crypto v0.1.0 // indirect
	google.golang.org/appengine v1.6.6 // indirect
	google.golang.org/protobuf v1.24.0 // indirect
	gopkg.in/ini.v1 v1.57.0 // indirect
)

replace github.com/CosmWasm/go-cosmwasm => github.com/terra-project/go-cosmwasm v0.10.3
