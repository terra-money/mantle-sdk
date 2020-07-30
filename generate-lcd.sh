CORE_VERSION=v0.4.0-rc.2
go get -u github.com/dolmen/yaml2json
go install github.com/dolmen/yaml2json
curl https://raw.githubusercontent.com/terra-project/core/v0.4.0-rc.2/client/lcd/swagger-ui/swagger.yaml > swagger.yaml
swagger generate client -f swagger.yaml -A lcd2