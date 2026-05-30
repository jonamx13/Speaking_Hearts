# ADR 005: Strategic Deferral of Containerization (Docker)

## Status
Accepted

## Context
During the active development and prototyping phases of the **Speaking_Hearts** system, the project requires rapid iteration cycles, frequent code modifications, and immediate feedback loops. We evaluated the use of Docker containers as a primary development environment versus utilizing the native host compiler and runtime.

## Decision
We have decided to **defer the implementation of Dockerization** until the final phase of the project (Phase 8). For the duration of active development, the system will rely on Go’s native toolchain (facilitated by `make dev`).

## Rationale

### 1. Developer Velocity and Feedback Loops
Native development allows for near-instantaneous compilation and execution via the Go compiler. Conversely, developing within a containerized environment necessitates a rebuild of the Docker image or the use of complex volume mounts for every significant change. These delays, however marginal individually, accumulate into significant overhead that impedes developer velocity during the early stages of architectural refinement.

### 2. Reduced Architectural Footprint
By removing the containerization layer from the daily development workflow, we minimize the cognitive overhead and environmental complexity for contributors. This ensures that focus remains squarely on the core Go concurrency engine and AI pipeline optimization, rather than troubleshooting container orchestration or networking issues internal to the Docker daemon.

### 3. Native Hardware Access
The system requires low-level access to hardware devices (e.g., USB microphones via `portaudio`). While hardware passthrough is possible in Docker, it introduces an additional layer of configuration that varies significantly across host operating systems (Windows vs. Linux). Native development bypasses these compatibility hurdles during the critical hardware-software integration phases.

## The Role of Docker in Final Deployment
It is important to clarify that this decision does not eliminate Docker from the project lifecycle. Docker is designated strictly as a **deployment and packaging tool**. In Phase 8, we will implement optimized multi-stage builds to ensure a consistent, reproducible production environment. This approach allows us to enjoy the benefits of rapid native development while retaining the portability and isolation guarantees of containerization for the final product delivery.

## Consequences
*   Developers must ensure their local Go environments meet the version requirements specified in `docs/SETUP.md`.
*   Final validation of environmental dependencies will be consolidated into the concluding Phase 8.
