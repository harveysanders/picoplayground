# Raspberry Pi Pico 2 W - TinyGo LED Blink

Simple LED blink example for Raspberry Pi Pico boards using TinyGo.

## Hardware Variants

Understanding which board you have is critical for choosing the right build target:

| Board        | Chip          | Onboard LED Location | TinyGo Target |
| ------------ | ------------- | -------------------- | ------------- |
| **Pico**     | RP2040        | GPIO 25              | `pico`        |
| **Pico W**   | RP2040 + WiFi | CYW43 WiFi chip      | `pico-w`      |
| **Pico 2**   | RP2350        | GPIO 25              | `pico2`       |
| **Pico 2 W** | RP2350 + WiFi | CYW43 WiFi chip      | `pico2-w`     |

**How to identify:**

- **W models** have a metal WiFi antenna shield on the board
- **Pico 2** boards are marked "Raspberry Pi Pico 2" and have RP2350 chip
- **Original Pico** boards have RP2040 chip

## Building and Flashing

### Quick Method: Build + Flash in One Command

The easiest way is to use `tinygo flash` which builds and flashes automatically:

**For Pico 2 W (RP2350 + WiFi):**
```bash
tinygo flash -target=pico2-w main.go
```

**For Pico W (RP2040 + WiFi):**
```bash
tinygo flash -target=pico-w main.go
```

**For Pico 2 (RP2350, no WiFi):**
```bash
tinygo flash -target=pico2 main.go
```

**For Pico (RP2040, no WiFi):**
```bash
tinygo flash -target=pico main.go
```

**Before running `tinygo flash`:**
1. Hold BOOTSEL button on the Pico board
2. Plug in USB cable while holding BOOTSEL
3. Release BOOTSEL - the board appears as a USB drive (RPI-RP2)
4. Run the `tinygo flash` command

The board will automatically reboot and run your program.

### Alternative: Build Then Flash Manually

If you prefer to build first and flash later:

```bash
# Build the .uf2 file
tinygo build -target=pico2-w -o main.uf2 main.go

# Then manually copy main.uf2 to the RPI-RP2 USB drive
```

## Understanding the LED Differences

### W Models (Pico W, Pico 2 W)

On WiFi-enabled boards, the onboard LED is **not connected to a GPIO pin**. Instead, it's controlled through the CYW43 WiFi chip:

```go
// This works on W models when using pico-w or pico2-w target
led := machine.LED  // machine.LED handles CYW43 chip communication
led.Configure(machine.PinConfig{Mode: machine.PinOutput})
led.High()  // LED on
```

### Non-W Models (Pico, Pico 2)

On non-WiFi boards, the LED is directly connected to GPIO 25:

```go
// This works on non-W models
led := machine.Pin(25)  // or machine.LED
led.Configure(machine.PinConfig{Mode: machine.PinOutput})
led.High()  // LED on
```

### External LEDs (All Models)

You can use any GPIO pin with an external LED + resistor (220Ω works well):

```go
led := machine.Pin(15)  // Any available GPIO
led.Configure(machine.PinConfig{Mode: machine.PinOutput})
led.High()  // LED on
```

**Wiring:**

```
┌─────────────┐
│   Pico 2 W  │
│             │
│  GPIO 15 ●──┼──────┐
│             │      │
│             │     ┌▼┐ LED
│             │     │ │ (long leg)
│             │     └┬┘
│             │      │
│             │     ┌┴┐
│             │     │R│ 220Ω
│             │     └┬┘
│             │      │
│      GND ●──┼──────┘
│             │
└─────────────┘
```

**Note:** The resistor can also go before the LED - both work since they're in series.

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
    "GOROOT": "/Users/YOUR_USERNAME/Library/Caches/tinygo/goroot-HASH"
  },
  "gopls": {
    "build.buildFlags": [
      "-tags=cortexm,baremetal,linux,arm,rp2350,rp,pico2,pico2-w,cyw43439,tinygo,purego,osusergo,math_big_pure_go,gc.conservative,scheduler.cores,serial.usb"
    ],
    "env": {
      "GOOS": "linux",
      "GOARCH": "arm",
      "GOROOT": "/Users/YOUR_USERNAME/Library/Caches/tinygo/goroot-HASH"
    }
  }
}
```

**To find your TinyGo GOROOT path:**

```bash
tinygo info [target, ex: pico2-w]
```

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
- [Raspberry Pi Pico Pinout](https://datasheets.raspberrypi.com/pico/Pico-R3-A4-Pinout.pdf)
- [TinyGo Machine Package](https://tinygo.org/docs/reference/microcontrollers/machine/)
