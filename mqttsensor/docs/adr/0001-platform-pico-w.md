# ADR-0001: Platform Selection - Raspberry Pi Pico W

## Status

Accepted

## Context

The project requires a microcontroller platform to read sensor data and transmit it to the cloud via WiFi. Key requirements:

- WiFi connectivity for cloud communication
- Low power consumption (battery or USB powered)
- Sufficient GPIO for sensor interfacing
- Cost-effective solution

A full Raspberry Pi was initially considered but consumes significantly more power (500mA+) than needed for this simple data collection task.

## Decision

Use the **Raspberry Pi Pico W** as the microcontroller platform.

The Pico W provides:

- RP2040 dual-core ARM Cortex-M0+ processor
- 264KB SRAM
- Integrated CYW43439 WiFi chip (2.4GHz 802.11n)
- 26 GPIO pins
- Low power consumption (~50-100mA active with WiFi)
- $6 USD price point

## Alternatives Considered

### Raspberry Pi (Full Linux SBC)

- **Rejected** - High power consumption (500mA-2A), overkill for sensor data collection, requires SD card and full OS management

### Arduino + WiFi Shield

- **Rejected** - Less integrated solution, additional cost and complexity, more wiring required

### ESP32

- **Considered viable** - Similar capabilities and price, but Pico W was chosen due to familiarity with Raspberry Pi ecosystem and better TinyGo support at the time

### Pico2 W

- **Rejected** - Newer platform with less community support and documentation compared to Pico W. Found un-resolved freeze bug on RP2350 when using TinyGo.

## Consequences

### Positive

- Low power consumption enables battery operation or simple USB power
- Integrated WiFi eliminates need for external modules
- Wide community support and documentation
- Cost-effective platform

### Negative

- Limited to 2.4GHz WiFi only (no 5GHz support)
- Constrained memory (264KB RAM) requires careful resource management
- No built-in filesystem (must use external flash for persistent storage)
- Single-threaded networking in current TinyGo implementation
