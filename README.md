![coverage](https://raw.githubusercontent.com/f-lombardo/generate/badges/.badges/master/coverage.svg)

# generate

`generate` is a command-line interface (CLI) tool written in Go for quickly generating sample and test data such as IBAN numbers and UUIDs. It is designed to fit seamlessly into developer workflows by automatically copying generated output directly to the system clipboard.

---

## 🚀 Features

- **IBAN Generation**: Generates valid IBAN numbers for various supported countries (e.g., Italy `IT`, Spain `ES`, Netherlands `NL`, etc.).
- **UUID Generation**: Supports generating both Version 4 (random) and Version 7 (time-ordered) UUIDs.
- **Clipboard Integration**: Automatically copies generated values to the system clipboard (can be disabled via flag).
- **Flexible Output Formats**: Supports plain text (`tab`) and JSON (`json`) output.

---

## 📦 Installation & Build

### Prerequisites
- [Go](https://golang.org/) 1.27 or higher.

### Building from Source
Clone the repository and build the binary using:

```bash
go build -o generate .
```

You can move the compiled binary into a directory in your `$PATH` (e.g., `/usr/local/bin` or `~/bin`) to run it from anywhere:

```bash
mv generate /usr/local/bin/
```

---

## 🛠 Usage

General syntax:

```bash
generate [global options] <command> [command options]
```

### Global Options

| Option | Type | Default | Description |
|---|---|---|---|
| `-clipboard` | `bool` | `true` | Automatically copies results to the system clipboard. Use `-clipboard=false` to disable. |
| `-output` | `string` | `tab` | Output format. Supported values: `tab` (plain text) or `json`. |

---

## 📋 Available Commands

### 1. `iban`
Generates a syntactically valid IBAN.

**Command Options:**
- `-country` *(string, default: `IT`)*: ISO 3166-1 alpha-2 country code (e.g., `IT`, `ES`, `DE`, `FR`, `NL`).

**Examples:**
```bash
# Generate an Italian IBAN and copy it to clipboard
generate iban

# Generate a Spanish IBAN
generate iban -country ES

# Format output as JSON
generate -output json iban -country IT
```

### 2. `uuid`
Generates a Universally Unique Identifier (UUID).

**Command Options:**
- `-version` *(string, default: `4`)*: UUID version to generate (`4` for random UUID, `7` for time-ordered UUID).

**Examples:**
```bash
# Generate a UUID v4 (default)
generate uuid

# Generate a UUID v7 (time-ordered)
generate uuid -version 7

# Generate a UUID v4 in JSON format
generate -output json uuid
```

---

## 💡 Advanced Examples

- **Generate an IBAN without copying to clipboard:**
  ```bash
  generate -clipboard=false iban -country DE
  ```

- **JSON output for use in scripts or pipelines:**
  ```bash
  generate -output json uuid -version 7
  ```

---

## 🧪 Running Tests

To run the unit test suite:

```bash
go test -v ./...
```

---

## 📄 License

This project is licensed under the [MIT License](LICENSE).
