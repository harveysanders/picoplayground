# ADC Sampling Plan – MQ Sensor on Pico (Updated)

## Goal

Obtain a stable, low-noise analog reading from a 5V MQ sensor module using a Pico ADC, sampled once per interval (e.g. 60s), suitable for logging/transmission.

---

## Hardware Assumptions

- MCU: Raspberry Pi Pico / Pico W / Pico 2 (RP2040/RP2350)
- Sensor module:
  - Powered at 5V
  - AO output
- Voltage divider (ORIENTATION IMPORTANT):
  - **AO → 20kΩ → ADC node → 27kΩ → GND**
  - Divider ratio ≈ 0.574× (5V → ~3.01V nominal)
- Optional but recommended:
  - Capacitor from ADC node → GND
    - 0.1 µF (100 nF) default
    - Can increase to 1 µF if readings are still jumpy
- Grounds:
  - Sensor GND, Pico GND, and AGND tied to the same ground net

---

## ADC Configuration

- ADC input pin: GP26 / ADC0
- ADC reference: default (3.3V via onboard filtering)
- ADC resolution: platform default
- No external ADC_VREF unless explicitly added later

---

## Sampling Strategy

- Sampling interval: once per minute (configurable)
- Each interval:
  - short burst sampling
  - discard first read
  - average burst
  - emit one value

---

## Sampling Algorithm (Per Interval)

1. Sleep until next sampling interval (e.g. 60 seconds).
2. (Optional) If any power gating just occurred:
   - Wait 10–100 ms for analog settling.
3. Perform one ADC read and discard it.
4. Take a burst of `N` ADC samples:
   - Recommended `N = 32` (64 if desired)
   - Delay between samples: 100 µs – 1 ms
5. Aggregate samples:
   - Compute arithmetic mean
   - Optional: trimmed mean (drop top/bottom 1–2 samples)
6. Store / publish / log the averaged value.
7. Return to sleep.

---

## Validation / Debug Checks

- Verify divider wiring matches the orientation above:
  - 20k is the "top" resistor (AO → ADC node)
  - 27k is the "bottom" resistor (ADC node → GND)
- Confirm max AO does not cause ADC node to exceed 3.3V.
- If readings still noisy:
  - Increase ADC node capacitor
  - Increase burst size
  - Shorten wires / improve grounding

---

## Non-goals / Explicit Exclusions

- No continuous streaming ADC
- No external precision voltage reference (yet)
- No AGND isolation experiments
- No DSP filtering beyond burst averaging
