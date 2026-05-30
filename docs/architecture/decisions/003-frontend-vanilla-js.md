# ADR 003: Frontend Architecture - Vanilla JavaScript

## Status
Accepted

## Context
The project needs a proof-of-concept (PoC) interface for the Operator Dashboard and the Subtitle Screen. We need to ensure the system remains lightweight, portable, and easy to run on resource-constrained hardware like a Raspberry Pi.

## Decision
We opted for **Vanilla JavaScript** and **Native CSS** instead of modern frameworks like React, Vue, or Angular during the initial development phase.

## Rationale
*   **Zero Build Step:** Using Vanilla JS allows the server to serve static files directly via `http.FileServer`. This avoids the complexity of Node.js build tools (Webpack, Vite, npm) on edge hardware, simplifying the deployment process.
*   **Performance:** Native DOM manipulation and standard CSS have the lowest possible footprint on the CPU and RAM of low-powered display devices.
*   **Focus on Core:** By minimizing frontend boilerplate, we keep the engineering focus on the core backend concurrency engine and the AI processing pipeline.
*   **Platform Independence:** Any modern browser can render the interface without needing a specific runtime environment or heavy library downloads.

## Consequences
*   UI state management becomes more manual as the dashboard complexity grows.
*   A future transition to a framework (likely Svelte or React via Wails) will be required when the product moves beyond the PoC phase and requires a more complex UI.
