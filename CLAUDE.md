# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

A desktop automation tool for recording and replaying keyboard and mouse operations as JavaScript macros. Uses a **gRPC-based architecture** with separate engine (gRPC server) and UI (Wails) components.

## Technology Stack

- **Go 1.24+** with Go workspaces
- **GUI**: Wails (desktop wrapper) + Vue 3 + Monaco Editor
- **gRPC**: google.golang.org/grpc with Protocol Buffers
- **Keyboard/mouse hooking**: gohook + robotgo
- **JavaScript engine**: goja (dop251/goja)

## Architecture

```
macro_operation/              # Go workspace root
├── go.work                   # Go workspace file
├── Makefile                  # Build automation
├── build.py                  # Python build script
│
├── macro-common/             # Shared module
│   └── proto/                # gRPC service definitions
│       └── macro.proto
│
├── macro-engine/             # Core engine (gRPC server, port 50051)
│   ├── cmd/server.go         # Entry point
│   └── internal/
│       ├── engine/
│       │   ├── events/       # Event capture (gohook)
│       │   ├── generator/    # JS script generation
│       │   └── executor/     # JS execution (goja)
│       ├── model/            # Event/macro data models
│       └── service/          # Recording/playback/macro services
│
├── macro-ui/                 # Wails-based GUI
│   ├── main.go               # Wails entry point
│   ├── app.go                # App struct with gRPC client
│   └── frontend/             # Vue 3 + Vite
│
└── macros/                   # Macro storage (JSON files)
```

## Common Commands

```bash
# First-time setup
make init

# Build everything
make build

# Run engine (gRPC server)
make run-engine

# Run UI (Wails)
make run-ui

# Generate gRPC code from proto
make proto

# Clean build artifacts
make clean
```

## Key Features

1. **Recording**: Captures keyboard/mouse events in real-time via gohook
2. **Script Generation**: Converts events to JavaScript (uses KeyDown/KeyUp pairs for IME compatibility)
3. **Playback**: Executes macros using goja with robotgo bindings
4. **Macro Management**: Save/load/delete macros (JSON files in `./macros`)
5. **gRPC Communication**: UI connects to engine on `localhost:50051`

## JavaScript API (macro-common/proto/macro.proto)

```javascript
// Mouse operations
mouseMove(x, y)           // Move mouse to position
mouseClick("left|right|middle")  // Click button
mouseDrag(x, y)           // Drag mouse
mouseScroll(delta)        // Scroll wheel

// Keyboard operations
keyDown(key)              // Press key
keyUp(key)                // Release key
keyType(text)             // Type text (Unicode via clipboard)
keyTap(key)               // Key press (down + up)

// System operations
sleep(ms)                 // Delay
screenshot()              // Take screenshot (returns base64)
log(message)              // Log to console
```

## Data Flow

1. **Recording**: UI calls `StartRecording` via gRPC → engine starts event capture
2. **Generation**: Events optimized (duplicate mousemove removed) → JavaScript generated
3. **Playback**: UI sends script → engine creates goja VM → registers robotgo bindings → executes
