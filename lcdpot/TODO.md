# TODO

## Smooth out ADC readings

https://chatgpt.com/share/e/694ad780-db38-8003-852b-97e94eaa2cb1

- Replace the 100k, 200k, resistors with 10k / 20k (same 2/3 ratio), or 4.7k / 10k.
- add 0.01–0.1 µF from ADC pin to GND (at the Pico side) to smooth sampling noise.

## Connect to Raspberry Pi

- Connect the Raspberry Pi to the Pico via [UART, USB, or I2C.](https://raspberrytips.com/pi-to-pico-communication/)

## Calc PPM from ADC value

- Read datasheet for ADC value to PPM conversion.
- Add a button to swap between PPM and ADC value display.
