```
# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview
This is a desktop automation tool for recording and replaying keyboard and mouse operations as JavaScript macros. It provides a command-line interface for macro management.

## Technology Stack
- Go 1.21+
- Keyboard/mouse hooking: gohook + robotgo
- JavaScript engine: goja
- Configuration: viper

## Project Structure
```
macro_operation/
├── cmd/macro/              # Application entry point (main.go)
├── internal/app/           # Core business services
│   ├── config_service.go   # Configuration management
│   ├── macro_service.go    # Macro management
│   ├── play_service.go     # Macro playback
│   └── record_service.go   # Operation recording
├── pkg/bindings/           # JavaScript API bindings
│   └── api_bindings.go
├── configs/                # Configuration files
├── macros/                 # Macro storage directory
├── docs/                   # Documentation
├── Makefile                # Build/test scripts
└── README.md               # Main documentation
```

## Common Commands
```bash
# Install dependencies
make deps

# Build the application
make build

# Run the application
make run

# Run all tests
make test

# Code formatting
make fmt

# Linting
make lint

# Run all checks (fmt + lint + test)
make check

# Build for all platforms
make build-all

# Clean build artifacts
make clean
```

## Key Features
1. **Recording**: Captures keyboard and mouse events in real-time
2. **Script Generation**: Automatically converts events to JavaScript scripts
3. **Playback**: Executes macros with configurable speed
4. **Macro Management**: Save, load, list, and delete macros

## JavaScript API (pkg/bindings/api_bindings.go)
The core API includes:
- Mouse operations: mouseMove, mouseClick, mouseDrag, mouseScroll
- Keyboard operations: keyDown, keyUp, keyType, keyTap
- System operations: sleep, screenshot, log
