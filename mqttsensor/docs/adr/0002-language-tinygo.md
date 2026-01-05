# ADR-0002: Programming Language - TinyGo

## Status

Accepted

## Context

Firmware development for the Raspberry Pi Pico W requires selecting a programming language. The traditional choice for embedded development is C/C++, but modern alternatives exist.

Key considerations:

- Developer productivity and familiarity
- Memory safety
- Tooling quality (IDE support, testing, formatting)
- Ecosystem and library availability
- Suitability for constrained environments

## Decision

Use **TinyGo** - a Go compiler for small places (microcontrollers and WebAssembly).

TinyGo compiles Go code to machine code suitable for microcontrollers while maintaining most of Go's language features and standard library compatibility.

## Alternatives Considered

### C/C++

- **Rejected** - While the traditional embedded choice with the largest ecosystem, it requires manual memory management, lacks built-in safety guarantees, and the developers are less familiar with it. Development iteration would be slower.

### MicroPython

- **Rejected** - Easy to get started but:
  - No compile-time type checking
  - Higher runtime memory overhead
  - Interpreted execution is slower
  - Harder to catch bugs before deployment

### Rust

- **Considered** - Excellent memory safety, but steeper learning curve and less mature embedded ecosystem compared to established options (C/C++, not TinyGo).

## Consequences

### Positive

- **Familiarity**: Leverages existing Go knowledge
- **Memory safety**: Garbage collection and bounds checking prevent common bugs
- **Fast iteration**: Quick compile times, familiar testing framework
- **Excellent tooling**: gopls LSP support, gofmt, go test all work
- **Type safety**: Compile-time type checking catches errors early
- **Readability**: Go's simplicity makes code maintainable

### Negative

- **Smaller ecosystem**: Fewer embedded libraries than C/C++
- **Feature limitations**: Some Go features unavailable (reflection, full concurrency)
- **Binary size**: Larger than hand-optimized C (though still fits in flash)
- **Garbage collection**: Can cause unpredictable pauses (mitigated by careful allocation)

## Related Decisions

This decision influences [ADR-0003](0003-network-stack-lneto.md) as the network stack must be TinyGo-compatible.
