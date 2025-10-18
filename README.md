Here’s a complete `README.md` file for **Global Build Manager** ready to use:

````markdown
# Global Build Manager (GBM)

**Global Build Manager (GBM)** is a smart, adaptive build and run manager for Go projects (or any CLI-based projects). GBM watches your project files, automatically rebuilds when you stop typing, and runs the latest version while gracefully stopping previous processes. Inspired by Claude CLI, it provides a clear and interactive console UI.

---

## Features

- **Adaptive Rebuild Timing:** Waits until you stop typing before rebuilding. Tracks per-file idle time and adapts for new files.
- **Parallel File Watching:** Handles multiple files simultaneously for efficient rebuild triggering.
- **Automatic Process Management:** Stops the previous version before running the new build.
- **Console UI:** Colored output for builds, errors, running process PID, and status messages.
- **Configurable via `exec.gbf`:** Define build directory, build/run commands, and ignored directories/files.
- **Ignore Libraries/Dependencies:** Skip watching directories like `vendor` or `node_modules` to avoid unnecessary rebuilds.

---

## Installation

1. Make sure you have [Go](https://golang.org/dl/) installed.
2. Clone the repository:

```bash
git clone https://github.com/samuelgirmametaferia/Global-Build-Manager.git
cd global-build-manager
````

3. Build GBM:

```bash
go build -o gbm main.go
```

4. Add GBM to your system PATH or run directly:

```bash
./gbm <project_path>
```

---

## Configuration (`exec.gbf`)

Place an `exec.gbf` file in the root of your project. Example:

```ini
# Build output directory
build=./build

# Build command
build_cmd=go build -o app .

# Run command
run_cmd=./build/app

# Ignore folders (comma-separated)
ignore=vendor,node_modules
```

* **build:** Directory where your compiled output should be placed.
* **build_cmd:** Command to build your project.
* **run_cmd:** Command to run the built project.
* **ignore:** Comma-separated directories or paths to ignore during watching.

---

## Usage

```bash
# Watch current folder
gbm

# Watch specific folder
gbm /path/to/project
```

### Console Output

* **Green:** Build and run success.
* **Red:** Errors during build or runtime.
* **Yellow:** Stopping previous process.
* **Blue:** PID of the currently running process.

---

## How It Works

1. GBM recursively watches your project directory.
2. Detects changes in `.go` files or `exec.gbf`.
3. Waits for adaptive idle time before rebuilding.
4. Stops the previous process gracefully.
5. Builds and runs the new version.
6. Updates idle time per file and project-wide average for adaptive behavior.

---

## Contributing

1. Fork the repository.
2. Create a feature branch.
3. Commit your changes.
4. Open a pull request.

---

## License

MIT License © 2025 Your Name

