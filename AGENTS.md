# Repository Guidelines

## Project Structure & Module Organization
This repo is a set of TinyGo prototypes for Raspberry Pi Pico/Pico W/Pico 2 boards. Each top-level folder is a self-contained example:
- `blinky/`, `pwmblinky/`, `lcd/`, `lcdpot/`, `potpick/`, `mqttsensor/` hold individual apps (each with `main.go`).
- `mqttsensor/` includes subpackages for WiFi, LCD, MQTT, and NTP (`mqttsensor/cyw43439/`, `mqttsensor/mqtt/`, `mqttsensor/lcd/`, `mqttsensor/ntp/`).
- Project notes live in per-app `README.md` files; wiring diagrams appear in `mqttsensor/sensor-lcd.fzz` and `pico-2-w-pinout.pdf`.

## Build, Test, and Development Commands
- `tinygo flash -target=<target> main.go` (run inside an app folder) builds and flashes to the board.
- `tinygo build -target=<target> -o main.uf2 main.go` builds a UF2 for manual drag-and-drop flashing.
- `make help` lists available Makefile targets.
- `make flash/mqttsensor` flashes the MQTT example; requires `MQTT_ADDR`, `WIFI_SSID`, `WIFI_PASS` env vars (see Makefile).

## Coding Style & Naming Conventions
- Go code follows standard `gofmt` formatting (tabs for indentation).
- Keep package and folder names short and lowercase (`mqttsensor`, `lcdpot`).
- Example apps are single-file (`main.go`) unless shared logic warrants a subpackage.

## Testing Guidelines
- No automated tests are currently present (`*_test.go` files are absent).
- If adding tests, use Go’s standard `testing` package and name tests `TestXxx`.

## Commit & Pull Request Guidelines
- Recent commit subjects are short, sentence-case, and descriptive (e.g., “Sending timestamp in MQTT messages”).
- Keep commits focused on one behavior change.
- PRs should include: a brief summary, target board (`pico`, `pico-w`, `pico2`, `pico2-w`), and any wiring notes or screenshots if UI/LCD output changes.

## Configuration & Local Setup Tips
- Set TinyGo build tags in `.vscode/settings.json` for proper IntelliSense (see `README.md`).
- For WiFi/MQTT builds, keep credentials in env vars rather than hardcoding.
