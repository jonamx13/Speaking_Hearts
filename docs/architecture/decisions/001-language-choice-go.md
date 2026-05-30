# ADR 001: Language Choice - Go (Golang)

## Status
Accepted

## Context
The system requires real-time processing of high-frequency audio data, simultaneous multi-language translation, and low-latency WebSocket distribution. We needed a language that provides high performance, safety, and a robust concurrency model.

## Decision
We chose **Go (Golang)** as the primary backend language for Speaking Hearts.

## Rationale
*   **Concurrency Model:** Go's goroutines and channels are significantly more efficient than Python's threading model, which is limited by the Global Interpreter Lock (GIL). This is critical for routing audio streams across parallel worker pools without performance bottlenecks.
*   **Safety vs. Performance:** While C++ offers maximum performance, Go provides comparable execution speed with the added benefit of memory safety and a significantly simpler development experience for HTTP/WebSocket integration.
*   **Deployment:** Go compiles into a single, static binary with no external runtime dependencies. This is ideal for deployment on edge hardware (like Raspberry Pi) in offline ceremonial environments where installing dependencies could be problematic.
*   **CGo Integration:** Go allows for clean integration with high-performance C libraries (like `portaudio` and `faster-whisper`), which are required for the acquisition and AI layers.

## Consequences
*   We must manage CGo overhead for AI model bindings.
*   The team must adhere to idiomatic Go patterns (e.g., CSP-style concurrency) to maintain a clean and maintainable codebase.
