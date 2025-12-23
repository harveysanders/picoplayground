# LCD + Pot (ADC) Example

Reads an analog value on ADC0 and prints the voltage + raw 16-bit reading to an
HD44780 LCD with a PCF8574 I2C backpack.

## Hardware

- Board: Raspberry Pi Pico 2 W (RP2350)
- LCD backpack: PCF8574/HD44780 (ex: Z-0234)
- I2C pins: SDA=GP4, SCL=GP5 (I2C0)
- ADC input: ADC0 (GP26)

Wiring:

```
Pico 2 W          LCD Backpack
3.3V  ───────────  VCC
GND   ───────────  GND
GP4   ───────────  SDA
GP5   ───────────  SCL

Pico 2 W          Potentiometer
3.3V  ───────────  VCC
GND   ───────────  GND
GP26  ───────────  WIPER
```

## Build + Flash

```bash
tinygo flash -target=pico2-w lcdpot/main.go
```

## What It Does

- Scans common I2C addresses (0x27, 0x3F).
- Initializes the LCD and turns the backlight on.
- Prints voltage, percentage, and the 16-bit ADC value.
- Updates every 500 ms.

## Pico 2 W Hang Note (WIP)

On a Pico 2 W (RP2350), this program can hang after a few minutes. The same
code runs fine on a Pico W (RP2040). Root cause is unknown; leaving this note
here to revisit later.

## Contrast Pot

If the screen is blank (even with backlight on), adjust the contrast pot on the
backpack.

## 5V LCD/I2C Callout

Many LCD backpacks are designed for 5V. If you power the backpack at 5V, its
I2C pull-ups are often 5V too. Pico GPIOs are not 5V tolerant.

Safe options:

- Keep the backpack at 3.3V (dimmer backlight, safe I2C levels).
- Use a proper I2C level shifter if you want to run the backpack at 5V.

## 5V Pot Wiring (Voltage Divider)

If you want to run the potentiometer from 5V but keep the ADC input safe, use a
voltage divider on the wiper. "Top" is the resistor between the wiper and the
ADC input, and "bottom" is the resistor between the ADC input and GND. Example
values that work well:

- 100k (top), 200k (bottom) to ground.

The higher resistance keeps the divider from loading the pot. This brings the
5V wiper range down to about 3.3V max at the ADC input.

Divider math:

```
Vout = Vin * (Rbottom / (Rtop + Rbottom))
Vout = 5V * (200k / (100k + 200k)) = 5V * (200/300) = 3.33V
```

Schematic:

```
VCC (5V) ──[ POT ]── GND
        |
     wiper (AO) ──[ 100k ]── ADC ──[ 200k ]── GND
                    (top)                     (bottom)
```
