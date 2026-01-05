#!/usr/bin/env python3
"""Build macro-engine and macro-ui executables."""

import argparse
import os
import subprocess
import sys

BASE_DIR = os.path.dirname(__file__)


def build_engine():
    """Build macro-engine executable."""
    os.chdir(os.path.join(BASE_DIR, "macro-engine"))

    env = os.environ.copy()
    env["CGO_ENABLED"] = "1"

    print("Building engine...")
    result = subprocess.run(
        ["go", "build", "-ldflags=-s -w", "-o", "engine.exe", "cmd/server.go"],
        env=env,
        capture_output=True,
        text=True,
        encoding="utf-8",
        errors="replace"
    )
    if result.returncode != 0:
        print(f"Build failed: {result.stderr}", file=sys.stderr)
        return 1

    # Compress with UPX if available
    if subprocess.run(["where", "upx"], capture_output=True, text=True, encoding="utf-8", errors="replace").returncode == 0:
        subprocess.run(["upx", "engine.exe"], capture_output=True)

    # Copy to bin directory
    os.makedirs("../bin", exist_ok=True)
    os.replace("engine.exe", "../bin/engine.exe")

    print("Engine built: bin/engine.exe")
    return 0


def build_ui():
    """Build macro-ui executable with Wails."""
    # Build frontend
    frontend_dir = os.path.join(BASE_DIR, "macro-ui", "frontend")
    os.chdir(frontend_dir)

    print("Installing frontend dependencies...")
    result = subprocess.run(["npm", "i"], capture_output=True, text=True, encoding="utf-8", errors="replace")
    if result.returncode != 0:
        print(f"npm install failed: {result.stderr}", file=sys.stderr)
        return 1

    print("Building frontend...")
    result = subprocess.run(["npm", "run", "build"], capture_output=True, text=True, encoding="utf-8", errors="replace")
    if result.returncode != 0:
        print(f"npm build failed: {result.stderr}", file=sys.stderr)
        return 1

    # Build UI with Wails
    ui_dir = os.path.join(BASE_DIR, "macro-ui")
    os.chdir(ui_dir)

    print("Building UI with Wails...")
    result = subprocess.run(["wails", "build"], capture_output=True, text=True, encoding="utf-8", errors="replace")
    if result.returncode != 0:
        print(f"Wails build failed: {result.stderr}", file=sys.stderr)
        return 1

    # Compress with UPX if available
    exe_path = os.path.join(ui_dir, "build", "bin", "macro-ui.exe")
    if subprocess.run(["where", "upx"], capture_output=True, text=True, encoding="utf-8", errors="replace").returncode == 0:
        subprocess.run(["upx", exe_path], capture_output=True)

    # Copy to bin directory
    bin_dir = os.path.join(BASE_DIR, "bin")
    os.makedirs(bin_dir, exist_ok=True)
    os.replace(exe_path, os.path.join(bin_dir, "macro-ui.exe"))

    print("UI built: bin/macro-ui.exe")
    return 0


def main():
    parser = argparse.ArgumentParser(description="Build macro project")
    parser.add_argument("--engine", action="store_true", help="Build engine only")
    parser.add_argument("--ui", action="store_true", help="Build UI only")
    parser.add_argument("--all", action="store_true", help="Build all (default)")

    args = parser.parse_args()

    # Default to all if no specific target
    if not args.engine and not args.ui:
        args.all = True

    exit_code = 0

    if args.all or args.engine:
        exit_code = build_engine()
        if exit_code != 0:
            return exit_code

    if args.all or args.ui:
        exit_code = build_ui()

    return exit_code


if __name__ == "__main__":
    sys.exit(main())
