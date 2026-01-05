# ADR-0004: Communication Protocol - MQTT

## Status

Accepted

## Context

The sensor device needs to transmit data to a cloud backend. The protocol must:

- Work well on constrained devices with limited memory and bandwidth
- Support reliable message delivery
- Handle intermittent connectivity gracefully
- Be suitable for streaming sensor data (temperature, humidity readings)

## Decision

Use **MQTT** (Message Queuing Telemetry Transport) as the communication protocol, implemented via the **natiu-mqtt** library.

MQTT is a lightweight publish/subscribe messaging protocol designed for IoT and constrained devices. It operates over TCP and uses a broker-based architecture.

## Alternatives Considered

### HTTP/REST
- **Rejected** - Higher overhead per message (headers, connection setup), request/response model less suited for continuous sensor streaming, no built-in message queuing

### WebSockets
- **Rejected** - More complex to implement, designed for bidirectional real-time communication which is overkill for periodic sensor readings, higher memory footprint

### CoAP (Constrained Application Protocol)
- **Considered** - UDP-based alternative designed for IoT, but MQTT's TCP foundation provides better reliability and wider broker support

## Consequences

### Positive
- **Lightweight**: Minimal packet overhead (as low as 2 bytes for small messages)
- **Pub/Sub model**: Natural fit for sensor data streaming; device publishes, subscribers receive
- **QoS levels**: Configurable delivery guarantees (0: at most once, 1: at least once, 2: exactly once)
- **Retained messages**: Broker can store last message for new subscribers
- **Wide ecosystem**: Many broker options (Mosquitto, HiveMQ, cloud providers)
- **natiu-mqtt**: TinyGo-compatible, minimal memory allocation

### Negative
- **Broker dependency**: Requires MQTT broker infrastructure
- **TCP overhead**: More expensive than UDP-based protocols for very constrained networks
- **Security**: Requires TLS for encrypted communication (adds memory/CPU overhead)

## Implementation

Using the `natiu-mqtt` library which is designed for:
- Minimal heap allocations
- Small binary size
- TinyGo compatibility
- Simple API

## Related Decisions

- Built on networking provided by [ADR-0003](0003-network-stack-lneto.md) (lneto stack)
