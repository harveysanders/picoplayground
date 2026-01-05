# ADR-0003: Network Stack - lneto

## Status

Accepted

## Context

With TinyGo selected as the programming language ([ADR-0002](0002-language-tinygo.md)), a TCP/IP network stack is needed that:

- Works with TinyGo's compiler constraints
- Integrates with the Pico W's CYW43439 WiFi chip
- Has a small memory footprint suitable for embedded use
- Provides standard networking primitives (TCP, UDP, DNS)

Go's standard `net` package is not available in TinyGo for embedded targets.

## Decision

Use **lneto** (lightweight network stack) by soypat.

lneto is a pure-Go network stack designed for embedded systems and TinyGo. It provides:
- TCP/IP implementation
- DHCP client
- DNS resolution
- Low memory footprint
- Direct integration with hardware drivers

## Alternatives Considered

### lwIP (via CGo)
- **Rejected** - Would require CGo bindings, adding complexity and potentially breaking TinyGo compatibility

### Custom implementation
- **Rejected** - TCP/IP is complex; using a tested implementation is more reliable

### Standard Go net package
- **Not available** - TinyGo does not support the standard library's net package on embedded targets

## Consequences

### Positive
- **TinyGo native**: Written in pure Go, works seamlessly with TinyGo
- **Lightweight**: Designed for memory-constrained environments
- **Actively maintained**: Regular updates and bug fixes
- **Official integration**: Used in TinyGo's official cyw43439 driver for Pico W
- **Good documentation**: Clear examples and API documentation

### Negative
- **Community-supported**: Not part of Go standard library; depends on maintainer availability
- **Less battle-tested**: Smaller user base than lwIP or other established stacks
- **Limited features**: May lack advanced networking features (acceptable for this use case)

## Related Decisions

- Depends on [ADR-0002](0002-language-tinygo.md) (TinyGo language choice)
- Enables [ADR-0004](0004-protocol-mqtt.md) (MQTT protocol implementation)
