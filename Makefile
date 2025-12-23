## help: print this help message
.PHONY: help
help:
	@echo 'Usage:'
	@sed -n 's/^##//p' ${MAKEFILE_LIST} | column -t -s ':' |  sed -e 's/^/ /'

## flash/mqtt: flash the MQTT client app to the Pico W
.PHONY: flash/mqtt
flash/mqtt:
	@LDFLAGS="-X 'github.com/harveysanders/picoplayground/mqttsensor/cyw43439.ssid=${WIFI_SSID}' \
		-X 'github.com/harveysanders/picoplayground/mqttsensor/cyw43439.pass=${WIFI_PASS}'"; \
	tinygo flash -target=pico-w -stack-size=16kb -monitor -ldflags="$$LDFLAGS" ./mqttsensor/...
