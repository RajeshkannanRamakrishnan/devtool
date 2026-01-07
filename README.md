# Devtool

![Go Version](https://img.shields.io/badge/go-1.23-blue)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

**Devtool** is a comprehensive developer productivity CLI designed to speed up daily tasks. Built with Go and Cobra.

## Features

-   **Process Management**:
    -   `kill`: Terminate processes by PID or Port number (e.g., kill the process on port 8080).
    -   `ports`: List all processes listening on network ports, with filtering capabilities.
-   **Web & Network**:
    -   `server`: Instantly start a static HTTP file server in the current directory.
-   **Utilities**:
    -   `base64`: Encode and decode Base64 strings.
    -   `md5`: Compute MD5 hashes of input strings (with uppercase support).
    -   `upper`: Convert text to uppercase.

## Installation

### Prerequisites

-   **Go 1.23+** (required for building from source)
-   [mise](https://mise.jdx.dev/) (optional, recommended for development environment management)

### Quick Install (Manual)

We provide convenience scripts to build and install `devtool` to your system path.

#### macOS / Linux

Run the installation script to build and install to `/usr/local/bin`:

```bash
./install.sh
```

> [!NOTE]
> This script builds with `CGO_ENABLED=0` to ensure compatibility and may request `sudo` access to move the binary.

#### Windows

Run the batch script to build and install to `%USERPROFILE%\bin`:

```batch
install.bat
```

> [!NOTE]
> The script will automatically add the install directory to your `PATH` if it's not already present.

### Building from Source

If you prefer to build manually:

```bash
# Clone the repository
git clone https://github.com/rajeshkannanramakrishnan/devtool.git
cd devtool

# Build the binary (disable CGO for better portability)
CGO_ENABLED=0 go build -o devtool main.go

# Move to a directory in your PATH, e.g.:
mv devtool /usr/local/bin/
```

## Usage

Run `devtool --help` to see the full list of commands.

### Base64 Encoding/Decoding

Encode a string:
```bash
devtool base64 "Hello World"
# Output: SGVsbG8gV29ybGQ=
```

Decode a string:
```bash
devtool base64 --decode "SGVsbG8gV29ybGQ="
# Output: Hello World
```

### Process Management (Kill)

Kill a process by PID:
```bash
devtool kill --pid 1234
```

Kill the process listening on a specific port (e.g., 8080):
```bash
devtool kill --port 8080
```

### MD5 Hashing

Generate MD5 hash:
```bash
devtool md5 "hello world"
# Output: 5eb63bbbe01eeed093cb22bb8f5acdc3
```

Generate uppercase hash:
```bash
devtool md5 --upper "hello world"
# Output: 5EB63BBBE01EEED093CB22BB8F5ACDC3
```

### Port Enumeration

List all listening ports:
```bash
devtool ports
```

Filter by process name:
```bash
devtool ports --filter chrome
```

Show executable paths:
```bash
devtool ports --show-path
```

### Static File Server

Start a server in the current directory on port 8080 (default):
```bash
devtool server
```

Start on a specific port:
```bash
devtool server --port 9090
```

### String Manipulation

Convert string to uppercase:
```bash
devtool upper "hello world"
# Output: HELLO WORLD
```

## Development

This project uses `mise` to manage tools like `go` and `goreleaser`.

1.  Install [mise](https://mise.jdx.dev/).
2.  Run `mise install` to install dependencies.
3.  Make your changes.
4.  Run `go build` to verify.

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.
