# ADR 004: Language Choice Justification - Go vs. Rust

## Status
Accepted

## Context
In the design of a real-time, multi-threaded audio routing system such as **Speaking_Hearts**, the choice between high-performance systems languages is critical. Rust was extensively evaluated as an alternative to Go, given its reputation for zero-cost abstractions and strict memory safety. However, for the initial implementation and architectural blueprint of this streaming pipeline, Go was selected.

## Decision
We chose **Go (Golang)** over **Rust** for the primary execution engine.

## Rationale

### 1. Architectural Alignment (CSP vs. Async/Await)
The **Speaking_Hearts** engine is fundamentally a streaming pipeline utilizing a **Fan-Out/Fan-In** pattern. Go’s native implementation of **Communicating Sequential Processes (CSP)** through goroutines and channels provides a more intuitive and direct mapping to this data-flow architecture. While Rust’s `async/await` (e.g., via the Tokio runtime) is highly performant, it introduces significant complexity in managing lifetimes and ownership across asynchronous boundaries—challenges that are natively abstracted by Go’s runtime without sacrificing the required real-time latency targets.

### 2. Developer Velocity and Academic Constraints
Given the constraints of an academic research timeline, the speed of iteration is a primary concern. Binding to complex C++ AI libraries (such as `faster-whisper`) is more streamlined via **CGo** than through Rust's FFI and bindgen layers. Go's lower cognitive overhead allows for rapid prototyping of the concurrency orchestration layer, ensuring that the system's structural integrity can be validated before the final academic submission.

### 3. Native Multicasting
Go’s channels allow for simple, synchronous-looking orchestration of the Fan-Out distributor. Implementing a similar multicasting audio stream in Rust often requires complex synchronization primitives (e.g., `Arc<Mutex<T>>` or specialized broadcast channels) that can obscure the architectural intent of the streaming layers.

## The Olive Branch: Future Rust Iterations
It is important to note that the current Go implementation serves as a functional and concurrent **blueprint**. We explicitly acknowledge that as the system matures beyond the proof-of-concept phase, a partial or complete rewrite in **Rust** would be a logical progression. Rust would offer:
*   Strict compile-time guarantees against data races.
*   Zero-cost abstractions for the high-frequency audio manipulation layer.
*   Improved CPU efficiency on extremely resource-constrained hardware.

The current architecture, defined by its channel-based boundaries, is designed to be language-agnostic, thus welcoming a future Rust implementation that adheres to the established Go-based concurrent specifications.
