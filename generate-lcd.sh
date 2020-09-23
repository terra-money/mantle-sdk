CORE_VERSION=v0.4.0
#curl https://raw.githubusercontent.com/terra-project/core/v0.4.0-rc.2/client/lcd/swagger-ui/swagger.yaml > swagger.yaml
rm -rf lcd
mkdir lcd
swagger generate client -f swagger.yaml -A client -t lcd
