# Blinky (TinyGo)

LED blink example for Raspberry Pi Pico boards using TinyGo.

The board will automatically reboot and run your program.

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

You can use any GPIO pin with an external LED + resistor (220 ohm works well):

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
│             │     │R│ 220 ohm
│             │     └┬┘
│             │      │
│      GND ●──┼──────┘
│             │
└─────────────┘
```

**Note:** The resistor can also go before the LED - both work since they're in series.

## Troubleshooting

### LED Doesn't Blink

- **Wrong target?** Make sure you're using `pico2-w` for Pico 2 W boards
- **Wrong board?** Check if your board has WiFi (W model) or not
- **Try external LED** on GPIO 15 with 220 ohm resistor to verify code works

## Resources

- [TinyGo Documentation](https://tinygo.org/docs/)
- [Raspberry Pi Pico Pinout](https://datasheets.raspberrypi.com/pico/Pico-R3-A4-Pinout.pdf)
- [TinyGo Machine Package](https://tinygo.org/docs/reference/microcontrollers/machine/)
