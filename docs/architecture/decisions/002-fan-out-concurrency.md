# ADR 002: Audio Stream Concurrency - Fan-Out Pattern

## Status
Accepted

## Context
The system must simultaneously process raw audio chunks for Speech-to-Text (STT) and write them to persistent storage (Fragmented Recording). Attaching multiple consumers to a single Go channel would result in "message stealing," where only one consumer receives any given chunk of data.

## Decision
We implemented a **Fan-Out (Multicast) Pattern** for the primary audio stream.

## Rationale
*   **Data Integrity:** By using a dedicated "Distributor" goroutine that iterates over the source microphone channel and sends copies to multiple destination channels, we ensure that both the STT processing layer and the Recording layer receive 100% of the audio data.
*   **Decoupling:** This pattern allows the recording layer (which flushes every 30 seconds) to operate at its own pace without blocking or delaying the high-frequency STT worker pool, which must provide real-time subtitles.
*   **Scalability:** Adding new consumers (e.g., a real-time waveform visualizer or an additional backup writer) becomes a simple configuration change in the distributor logic.

## Consequences
*   Increased memory usage due to multiple buffered channels.
*   A slight increase in complexity to manage the lifecycle of additional channels and the distributor goroutine.
