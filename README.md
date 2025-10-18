# Global Build Manager

A lightweight, intelligent file watcher that automatically builds and runs your Go projects when files change. Perfect for rapid development cycles and continuous feedback.

## Features

- 🔄 **Automatic Watching**: Recursively watches all directories in your project
- 🚀 **Smart Rebuild**: Automatically rebuilds and restarts your application on file changes
- ⚡ **Adaptive Debouncing**: Learns your project's save patterns for optimal rebuild timing
- 🎨 **Colored Output**: Clear, colored console output for better readability
- 🎯 **Selective Ignoring**: Configure which files and directories to ignore
- 📁 **Directory Detection**: Automatically watches new directories as they're created

## Installation

### From Source

```bash
git clone https://github.com/samuelgirmametaferia/Global-Build-Manager.git
cd Global-Build-Manager
go build -o gbm
```

### Install Globally

```bash
go install github.com/samuelgirmametaferia/compile-when-edited@latest
```

## Usage

### Basic Usage

Run the tool in your project directory:

```bash
./compile-when-edited
```

Or specify a path to watch:

```bash
./compile-when-edited /path/to/your/project
```

### Configuration File

Create an `exec.gbf` file in your project root to configure build and run commands:

```
build=/path/to/build/directory
build_cmd=go build -o myapp
run_cmd=./myapp
ignore=vendor,node_modules,.git
```

#### Configuration Options

- **`build=`**: Directory where build commands will be executed (default: current directory)
- **`build_cmd=`**: Command to build your project (e.g., `go build`, `make`, `npm run build`)
- **`run_cmd=`**: Command to run your application (e.g., `./myapp`, `npm start`)
- **`ignore=`**: Comma-separated list of paths to ignore (e.g., `vendor,tmp,.git`)

## How to Add Files

### Automatic File Watching

**The Global Build Manager automatically watches ALL files in your project!** Here's how it works:

#### 1. **Go Files (`.go`)**
Simply create or edit any `.go` file in your project:

```bash
touch main.go           # Creates a new Go file
vim handlers.go         # Edit an existing file
```

The watcher will **automatically detect** these changes and trigger a rebuild.

#### 2. **New Directories**
When you create a new directory, it's automatically added to the watch list:

```bash
mkdir pkg/utils         # New directory is automatically watched
touch pkg/utils/helper.go  # Files in new directory trigger rebuilds
```

#### 3. **Configuration Changes**
Modify your `exec.gbf` file to change build/run settings:

```bash
vim exec.gbf           # Changes take effect on next rebuild
```

### What Files Are Watched?

By default, the tool watches:
- ✅ All `.go` files in the project
- ✅ The `exec.gbf` configuration file
- ✅ All subdirectories (recursively)

### Ignoring Files and Directories

To **exclude** certain files or directories from triggering rebuilds, add them to the `ignore=` option in your `exec.gbf`:

```
ignore=vendor,tmp,.git,node_modules,test_data
```

This prevents changes in these paths from triggering rebuilds while still allowing you to work on them.

## Examples

### Example 1: Simple Go Application

**Project Structure:**
```
myproject/
├── exec.gbf
├── main.go
└── handlers/
    └── api.go
```

**exec.gbf:**
```
build=.
build_cmd=go build -o myapp
run_cmd=./myapp
ignore=tmp,.git
```

**Workflow:**
1. Start the watcher: `./compile-when-edited`
2. Edit `main.go` or `handlers/api.go`
3. Save your changes
4. Application automatically rebuilds and restarts!

### Example 2: Multi-Module Project

**exec.gbf:**
```
build=/home/user/myproject
build_cmd=go build -o bin/server cmd/server/main.go
run_cmd=./bin/server
ignore=vendor,docs,scripts,.git
```

### Example 3: Project with Tests

**exec.gbf:**
```
build=.
build_cmd=go test ./... && go build -o app
run_cmd=./app --dev
ignore=testdata,vendor
```

## How It Works

1. **Initial Setup**: On startup, the tool recursively adds all directories to the file watcher
2. **File Change Detection**: When a `.go` file or `exec.gbf` is modified, the change is detected
3. **Smart Debouncing**: The tool waits for a brief period (adapts to your save patterns) to batch changes
4. **Build Phase**: Stops any running process, reloads config, and executes your build command
5. **Run Phase**: If the build succeeds, runs your application
6. **Repeat**: Continues watching for more changes

## Advanced Features

### Adaptive Debouncing

The tool learns from your editing patterns:
- Tracks time between file saves
- Calculates optimal wait time before rebuilding
- Reduces unnecessary rebuilds during rapid typing
- Improves performance for different project sizes

### Process Management

- Automatically stops the previous process before starting a new one
- Properly handles process cleanup
- Shows process PID for debugging
- Graceful process termination

## Troubleshooting

### Changes Not Triggering Rebuilds

**Check:**
1. Is the file a `.go` file or `exec.gbf`?
2. Is the path listed in the `ignore=` configuration?
3. Is the file in a hidden directory (starts with `.`)?

### Build Failures

**Check:**
1. Is your `build_cmd` correct in `exec.gbf`?
2. Does the `build=` path exist?
3. Are there compilation errors in your code?

### Process Not Starting

**Check:**
1. Did the build succeed?
2. Is your `run_cmd` correct?
3. Are there runtime errors preventing startup?

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## Credits

Built with:
- [fsnotify](https://github.com/fsnotify/fsnotify) - Cross-platform file system notifications
- [color](https://github.com/fatih/color) - Colored terminal output
