# TinyGo Pico Prototypes

A collection of small TinyGo projects for prototyping on Raspberry Pi Pico, Pico W,
Pico 2, and Pico 2 W boards.

## Projects

- `blinky` - basic LED blink example (see `blinky/README.md`)
- `pwmblinky` - PWM LED dimming
- `lcd` - LCD display demo
- `lcdpot` - LCD + potentiometer
- `potpick` - potentiometer input demo

## Targets

Choose the TinyGo target that matches your board:

| Board        | Chip          | TinyGo Target |
| ------------ | ------------- | ------------- |
| **Pico**     | RP2040        | `pico`        |
| **Pico W**   | RP2040 + WiFi | `pico-w`      |
| **Pico 2**   | RP2350        | `pico2`       |
| **Pico 2 W** | RP2350 + WiFi | `pico2-w`     |

## Build and Flash (Typical)

1. `cd` into a project directory
2. Flash the project:

```bash
tinygo flash -target=<target> main.go
```

See each project directory for wiring notes and any project-specific details.

### Alternative: Build Then Flash Manually

If you prefer to build first and flash later:

```bash
# Build the .uf2 file
tinygo build -target=pico2-w -o main.uf2 main.go

# Then manually copy main.uf2 to the RPI-RP2 USB drive
```

## VSCode Setup

To get proper IntelliSense and eliminate errors in VSCode, you need to configure gopls with TinyGo build tags.

### Getting Build Tags

To see what build tags TinyGo uses for your target:

```bash
tinygo info pico2-w | grep "build tags:"
```

### VSCode Settings

Create or update `.vscode/settings.json`:

```json
{
  "go.toolsEnvVars": {
    "GOOS": "linux",
    "GOARCH": "arm",
    "GOROOT": "$HOME/Library/Caches/tinygo/goroot-HASH"
  },
  "gopls": {
    "build.buildFlags": [
      "-tags=cortexm,baremetal,linux,arm,rp2350,rp,pico2,pico2-w,cyw43439,tinygo,purego,osusergo,math_big_pure_go,gc.conservative,scheduler.cores,serial.usb"
    ],
    "env": {
      "GOOS": "linux",
      "GOARCH": "arm",
      "GOROOT": "$HOME/Library/Caches/tinygo/goroot-HASH"
    }
  }
}
```

**To find your TinyGo GOROOT path:**

```bash
tinygo info [target, ex: pico2-w]
```

Look for the `cached GOROOT:` line in the output and copy that path into
`GOROOT` in `.vscode/settings.json`.

**Why this is needed:**

- TinyGo uses different build tags than standard Go (arm, baremetal, cortexm, etc.)
- Without these settings, VSCode's gopls will show errors for TinyGo-specific packages like `machine`
- The GOROOT points to TinyGo's modified standard library

### Reload VSCode

After updating settings:

1. Press `Cmd+Shift+P` (Mac) or `Ctrl+Shift+P` (Windows/Linux)
2. Type "Developer: Reload Window"
3. Press Enter

## Troubleshooting

### LED Doesn't Blink

- **Wrong target?** Make sure you're using `pico2-w` for Pico 2 W boards
- **Wrong board?** Check if your board has WiFi (W model) or not
- **Try external LED** on GPIO 15 with 220Ω resistor to verify code works

### VSCode Shows Errors for `machine` Package

- Check that `.vscode/settings.json` has the correct build tags
- Verify GOROOT path exists
- Reload VSCode window after changing settings

### Build Errors

```bash
# Verify TinyGo is installed
tinygo version

# Check available targets
tinygo targets | grep pico
```

## Resources

- [TinyGo Documentation](https://tinygo.org/docs/)
- [TinyGo Machine Package](https://tinygo.org/docs/reference/microcontrollers/machine/)
- [Raspberry Pi Pico Pinout](https://datasheets.raspberrypi.com/pico/Pico-R3-A4-Pinout.pdf)
