# Potpick (ADC -> PWM) Example

Reads a potentiometer on ADC0 and maps the value to PWM duty on an LED.

## Hardware

- Board: Raspberry Pi Pico 2 W (RP2350)
- Potentiometer: outer legs to 3.3V and GND, wiper to ADC0 (GP26)
- LED: an external LED + resistor on GP15 to GND

Basic wiring:

```
Pico 2 W          Potentiometer
3.3V  ───────────  VCC (one outer leg)
GND   ───────────  GND (other outer leg)
GP26  ───────────  Wiper (middle leg)

Pico 2 W          LED + resistor
GP15  ──[220Ω]──► LED ── GND
```

## Build + Flash

```bash
tinygo flash -target=pico2-w potpick/main.go
```

## What It Does

- Configures ADC0 to read the potentiometer.
- Configures PWM7 and uses GP15 as the PWM output.
- Maps the ADC reading (0–65535) to PWM duty cycle.
- Prints the raw ADC value over serial.

## 5V LCD/I2C Callout

If you’re also using a 5V LCD backpack (PCF8574/HD44780), the I2C lines may
be pulled up to 5V. The Pico GPIOs are not 5V tolerant. Use a level shifter
or keep the backpack at 3.3V to protect the Pico.
