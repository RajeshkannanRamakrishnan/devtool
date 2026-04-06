# Devtool

![Go Version](https://img.shields.io/badge/go-1.23-blue)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

**Devtool** is a comprehensive developer productivity CLI designed to speed up daily tasks. Built with Go and Cobra.

## Features

-   **Process Management**:
    -   `kill`: Terminate processes by PID or Port number (e.g., kill the process on port 8080).
    -   `ports`: List all processes listening on network ports, with filtering capabilities.
-   **Web & Network**:
    - `server`: Start a lightweight HTTP/HTTPS server that responds `200 OK` with JSON to any request and logs request details.
    - `ssl`: Check SSL certificate issuer, expiry date, and days remaining for a domain.
-   **Productivity**:
    -   `standup`: Generate a git daily standup report across multiple repositories.
-   **Utilities**:
    -   `base64`: Encode and decode Base64 strings, files, or stdin.
    -   `csv split`: Split large CSV files into smaller files while preserving headers.
    -   `md2pdf`: Convert Markdown files to PDF.
    -   `md5`: Compute MD5 hashes of strings, files, or stdin.
    -   `sha256`: Compute SHA256 hashes of strings, files, or stdin.
    -   `upper`: Convert text to uppercase.
    -   `json`: Pretty print or minify JSON with colors.

## Installation

### Prerequisites

-   **Go 1.23+** (required for building from source)
-   [mise](https://mise.jdx.dev/) (optional, recommended for development environment management)

### Quick Install (Manual)

We provide convenience scripts to build and install `devtool` to your system path.

#### macOS / Linux

Run the installation script to build and install to `~/.local/bin`:

```bash
./install.sh
```

> [!NOTE]
> This script builds with `CGO_ENABLED=0` to ensure compatibility, creates `~/.local/bin` if needed, and prints a PATH hint if that directory is not already available in your shell.

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
mv devtool ~/.local/bin/
```

## Usage

Run `devtool --help` to see the full list of commands.

### Flexible Input Support
Commands like `md5`, `sha256`, and `base64` support input from:
1. **Arguments**: `devtool md5 "text"`
2. **Files**: `devtool md5 myfile.txt`
3. **Stdin**: `echo "text" | devtool md5`

### Base64 Encoding/Decoding

Encode a string, file, or stdin:
```bash
devtool base64 "Hello World"
# Output: SGVsbG8gV29ybGQ=

devtool base64 myfile.txt
```

Decode:
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

### Hashing (MD5 & SHA256)

Generate MD5 hash:
```bash
devtool md5 "hello world"
# Output: 5eb63bbbe01eeed093cb22bb8f5acdc3
```

Generate SHA256 hash:
```bash
devtool sha256 "hello world"
# Output: b94d27b9934d3e08a52e52d7da7dabfac484efe37a5380ee9088f7ace2efcde9
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

### HTTP Response Server

Start a local server on port 8080 (default). Every request returns `200 OK` with a JSON body:

```bash
devtool server
# Response body: {"status": "ok"}
```

Start on a specific port:
```bash
devtool server --port 9090
```

Start with HTTPS (Self-Signed Certificate):
```bash
devtool server --ssl
# Listening at https://localhost:8080
```

Return a custom status and body:
```bash
devtool server --status 201 --body '{"ok":true}'
```

Add custom headers:
```bash
devtool server --header 'Content-Type: application/json' --header 'X-Debug: true'
```

Simulate latency:
```bash
devtool server --delay 250ms
```

Control how much of the request body is captured in logs:
```bash
devtool server --log-body-limit 8192
```

The server logs each request as structured JSON, including:
```text
- timestamp
- remote address
- method
- URL
- headers
- request body
- status
- response size
- duration
```

If the captured request body exceeds the configured log limit, the log includes `"body_truncated": true`.

### String Manipulation

Convert string to uppercase:
```bash
devtool upper "hello world"
# Output: HELLO WORLD

```

### Git Standup

Generate a report of your commits across git repositories.

Search for git repos in the current directory and list commits by the current user from the last 24 hours:
```bash
devtool standup
```

Look back 5 days:
```bash
devtool standup --days 5
```

Filter by a specific author name:
```bash
devtool standup --author "John Doe"
```

Scan a specific directory for repositories:
```bash
devtool standup --path ~/projects/opensource
```

### JSON Utilities

Pretty print JSON (with colors):
```bash
echo '{"foo":"bar"}' | devtool json
# Output:
# {
#   "foo": "bar"
# }
```

Minify JSON:
```bash
echo '{\n  "foo": "bar"\n}' | devtool json --minify
# Output: {"foo":"bar"}
```

### SSL Certificate Check

Inspect the certificate presented by a domain:
```bash
devtool ssl example.com
```

Use an explicit host and port if needed:
```bash
devtool ssl example.com:8443
```

### CSV Split

Split a CSV file into chunks of 1000 rows each:
```bash
devtool csv split large.csv
```

Customize rows per file, output directory, and filename prefix:
```bash
devtool csv split data.csv --rows 500 --out ./output --prefix part_
```

### Markdown To PDF

Convert a Markdown file to PDF:
```bash
devtool md2pdf notes.md
# Generates notes.pdf
```

Specify the output path explicitly:
```bash
devtool md2pdf notes.md exports/notes.pdf
```


## Development

This project uses `mise` to manage tools like `go` and `goreleaser`.

1.  Install [mise](https://mise.jdx.dev/).
2.  Run `mise install` to install dependencies.
3.  Make your changes.
4.  Run `go build` to verify.

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.
