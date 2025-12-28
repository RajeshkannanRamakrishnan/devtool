# Devtool

Devtool is a developer productivity CLI featuring a collection of handy utilities to speed up daily tasks. Built with Go and Cobra.

## Features

- **Process Management**:
  - `kill`: Terminate processes by PID or Port number.
  - `ports`: List all processes listening on network ports (with filtering options).
- **String Utilities**:
  - `md5`: Compute MD5 hashes of input strings.
  - `upper`: Convert text to uppercase.

## Installation

### Prerequisites

- Go 1.13+ (if building from source)
- [mise](https://mise.jdx.dev/) (optional, for managing dependencies)

### Manual Installation

We provide convenience scripts to build and install `devtool` to your system path.

#### macOS / Linux
Run the installation script to build and install valid for `/usr/local/bin`:
```bash
./install.sh
```
*Note: May require sudo password to move the binary to `/usr/local/bin`.*

#### Windows
Run the batch script to build and install to `%USERPROFILE%\bin`:
```batch
install.bat
```
*Note: This script will automatically update your user `PATH` environment variable if needed.*

### Building from Source

```bash
go build -o devtool main.go
# Add the binary/current directory to your PATH
```

## Usage

### Kill Process
Kill a process by PID:
```bash
devtool kill --pid 1234
```

Kill a process listening on a specific port (e.g., 8080):
```bash
devtool kill --port 8080
```

### List Ports
List all processes listening on ports:
```bash
devtool ports
```

List ports with filter by process name:
```bash
devtool ports --filter python
```

Show full executable paths:
```bash
devtool ports --show-path
```

### MD5 Hash
Generate MD5 hash:
```bash
devtool md5 "hello world"
```

Generate uppercase hash:
```bash
devtool md5 -u "hello world"
```

### Uppercase
Convert string to uppercase:
```bash
devtool upper "hello world"
```

## Development

This project uses `mise` to manage tools like `go` and `goreleaser`.

1. Install [mise](https://mise.jdx.dev/).
2. Run `mise install` to install dependencies.
3. Make your changes.
4. Run `go build` to verify.


## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.
