# LCD (I2C HD44780) Example

TinyGo driver + example for an HD44780 LCD with a PCF8574 I2C backpack
(commonly labeled Z-0234).

## Hardware

- Board: Raspberry Pi Pico 2 W (RP2350)
- LCD backpack: [PCF8574](https://www.ti.com/product/PCF8574)/[HD44780](https://cdn.sparkfun.com/assets/9/5/f/7/b/HD44780.pdf) (ex: Z-0234)
- I2C pins: SDA=GP4, SCL=GP5 (I2C0)

Wiring:

```
Pico 2 W          LCD Backpack
3.3V  ───────────  VCC
GND   ───────────  GND
GP4   ───────────  SDA
GP5   ───────────  SCL
```

## Build + Flash

```bash
tinygo flash -target=pico2-w lcd/main.go
```

## What It Does

- Initializes the LCD in 4-bit mode over I2C.
- Turns the backlight on.
- Writes two lines of text.

## Contrast Pot

If the screen is blank (even with backlight on), adjust the contrast pot on the
backpack. This is the most common cause of "no text" issues.

## 5V LCD/I2C Callout

Many LCD backpacks are designed for 5V. If you power the backpack at 5V, its
I2C pull-ups are often 5V too. Pico GPIOs are not 5V tolerant.

Safe options:

- Keep the backpack at 3.3V (dimmer backlight, safe I2C levels).
- Use a proper I2C level shifter if you want to run the backpack at 5V.
