## help: print this help message
.PHONY: help
help:
	@echo 'Usage:'
	@sed -n 's/^##//p' ${MAKEFILE_LIST} | column -t -s ':' |  sed -e 's/^/ /'

## flash/mqtt: flash the MQTT client app to the Pico W. Pass env vars: MQTT_ADDR, WIFI_SSID, WIFI_PASS
.PHONY: flash/mqttsensor
flash/mqttsensor:
	@LDFLAGS="-X 'github.com/harveysanders/picoplayground/mqttsensor/cyw43439.ssid=${WIFI_SSID}' \
		-X 'main.mqttServerAddr=${MQTT_ADDR}' \
		-X 'github.com/harveysanders/picoplayground/mqttsensor/cyw43439.pass=${WIFI_PASS}'"; \
	tinygo flash -target=pico-w -stack-size=16kb -monitor -ldflags="$$LDFLAGS" ./mqttsensor/...
