# Build Paths: Sensor System Development Workflow

This document outlines the development task paths for building sensor systems, showing shared infrastructure vs sensor-specific work.

## High-Level Workflow

```mermaid
flowchart TD
    subgraph shared["Shared Infrastructure"]
        WIFI[WiFi Configuration]
        WIFIHTTP[HTTP Config Server<br/><i>optional future</i>]
        MQTT[MQTT Broker Config]
        NTP[NTP Time Sync]
        LCD[LCD Display Setup]
        CHAN[Channel-Based Data Flow]
    end

    subgraph fork["Sensor Type Fork"]
        DECIDE{Sensor Type?}

        subgraph mq["MQ Analog Path"]
            MQ_HW[Voltage Divider Circuit<br/>5V → 3.3V]
            MQ_SW[ADC Burst Sampling<br/>+ Averaging]
            MQ_CAL[Voltage-to-PPM<br/>Calibration Curves]
            MQ_WARM[Warm-up Period<br/>24-48 hours]
        end

        subgraph echem["Electrochemical AFE Path"]
            EC_HW[AFE Circuit Wiring]
            EC_ADC[External ADC<br/>~6ch per AFE, SPI]
            EC_READ[Read Channels:<br/>Pt1000+/-, WE/AUX x2]
            EC_OFFSET[Subtract Offset]
            EC_SENS[Apply Sensitivity<br/>nA/ppm]
            EC_TEMP[Temperature Comp<br/>via Pt1000]
            EC_CROSS[Cross-Gas Correction<br/><i>optional</i>]
            EC_WARM[Warm-up Period<br/>Minutes]
        end
    end

    subgraph converge["Convergence Point"]
        STRUCT[Define SensorReading Fields]
        SMOOTH[Smoothing Algorithms]
        NORM[Normalization]
        PUB[MQTT Publishing]
    end

    %% Shared infrastructure flow
    WIFI --> WIFIHTTP
    WIFI --> MQTT
    MQTT --> NTP
    NTP --> LCD
    LCD --> CHAN
    CHAN --> DECIDE

    %% Fork paths
    DECIDE -->|MQ Analog| MQ_HW
    DECIDE -->|Electrochemical AFE| EC_HW

    %% MQ path
    MQ_HW --> MQ_SW
    MQ_SW --> MQ_CAL
    MQ_CAL --> MQ_WARM
    MQ_WARM --> STRUCT

    %% Electrochemical path
    EC_HW --> EC_ADC
    EC_ADC --> EC_READ
    EC_READ --> EC_OFFSET
    EC_OFFSET --> EC_SENS
    EC_SENS --> EC_TEMP
    EC_TEMP --> EC_CROSS
    EC_CROSS --> EC_WARM
    EC_WARM --> STRUCT

    %% Convergence flow
    STRUCT --> SMOOTH
    SMOOTH --> NORM
    NORM --> PUB
```

## Task Breakdown

### Shared Infrastructure (All Sensor Types)

| Task | Current Implementation | Notes |
|------|----------------------|-------|
| WiFi Configuration | Flash-time linker flags | Credentials set at compile time |
| HTTP Config Server | Not implemented | Future: web-based settings page |
| MQTT Broker Config | Flash-time linker flags | Broker URL, topic prefix |
| NTP Time Sync | lwIP SNTP | Timestamps for readings |
| LCD Display | HD44780 via I2C | Real-time sensor values |
| Channel Architecture | Go channels | Decoupled producer/consumer |

### Sensor-Specific: MQ Analog Path

| Task | Details |
|------|---------|
| **Hardware** | Voltage divider circuit (5V heater → 3.3V ADC safe) |
| **Software** | ADC burst sampling with averaging to reduce noise |
| **Calibration** | Sensor-specific voltage-to-ppm curves (datasheet + tuning) |
| **Warm-up** | 24-48 hour burn-in for stable readings |

### Sensor-Specific: Electrochemical AFE Path

| Task | Details |
|------|---------|
| **Hardware** | AFE circuit wiring + external ADC |
| **External ADC** | ~6 channels per AFE, SPI interface to MCU |
| **Channel Mapping** | Pt1000+, Pt1000- (temp), WE/AUX per sensor slot |
| **Offset Subtraction** | Remove baseline offset voltage |
| **Sensitivity** | Apply nA/ppm conversion factor (from datasheet) |
| **Temp Compensation** | Use Pt1000 RTD reading for drift correction |
| **Cross-Gas Correction** | Optional: compensate for interfering gases |
| **Warm-up** | Minutes (per datasheet) |

### Convergence Point (Both Paths)

| Task | Notes |
|------|-------|
| `SensorReading` struct | Fields depend on sensor type (ppm, ppb, raw voltage, etc.) |
| Smoothing algorithms | Moving average, exponential smoothing (sensor-dependent tuning) |
| Normalization | Convert to standard units if needed |
| MQTT Publishing | JSON payload → broker |

## Key Differences Summary

| Aspect | MQ Analog | Electrochemical AFE |
|--------|-----------|---------------------|
| Interface | Built-in ADC | External SPI ADC (~6ch/AFE) |
| Voltage handling | Requires divider (5V→3.3V) | AFE handles signal conditioning |
| Signal processing | Voltage → calibration curve → ppm | ADC → offset → sensitivity → temp comp → ppm |
| Calibration | Manual voltage-to-ppm curves | Datasheet sensitivity + offset values |
| Temp compensation | Usually ignored | Required (Pt1000 RTD) |
| Cross-gas | N/A | Optional correction for interferents |
| Warm-up time | 24-48 hours | Minutes |
| Cost | Lower | Higher |
| Accuracy | Moderate | Higher |
