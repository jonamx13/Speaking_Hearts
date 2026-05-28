# Setup Instructions

This guide will help you set up the **Speaking Hearts** project on your local machine for development.

## Pre-requisites

Before you begin, ensure you have the following tools installed:

1.  **Go 1.24+**: The core backend is built with Go.
    *   [Download Go](https://go.dev/dl/)
2.  **Git**: For version control.
    *   [Download Git](https://git-scm.com/downloads)
3.  **Make**: Used for automating build and development tasks.

---

## Installing Make on Windows

Since `make` is not natively included in Windows, you can install it using one of the following methods:

### Option 1: Using winget (Recommended)
Open PowerShell and run:
```powershell
winget install ezwinports.make
```

### Option 2: Using Chocolatey
If you have [Chocolatey](https://chocolatey.org/) installed:
```powershell
choco install make
```

### Option 3: Using MinGW (via MSYS2)
If you prefer a full Unix-like environment:
1.  Install [MSYS2](https://www.msys2.org/).
2.  Run the following command in the MSYS2 terminal:
    ```bash
    pacman -S make
    ```

---

## Getting Started

1.  **Clone the repository**:
    ```bash
    git clone `https://github.com/jonamx13/Speaking_Hearts.git`
    cd Speaking_Hearts
    ```

2.  **Verify the installation**:
    Run the following command to see if the environment is ready:
    ```bash
    make build
    ```

3.  **Run the development server**:
    ```bash
    make dev
    ```
    The dashboard will be available at `http://localhost:8080/dashboard.html`.

## Additional Dependencies (Future Phases)
---
**This section will be documented in-progress as the projects progresses**
---