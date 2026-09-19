[![test](https://github.com/f-lombardo/generate/actions/workflows/tests.yml/badge.svg?branch=master)](https://github.com/f-lombardo/generate/actions/workflows/tests.yml)
![coverage](https://raw.githubusercontent.com/f-lombardo/generate/badges/.badges/master/coverage.svg)
[![govulnchecks](https://github.com/f-lombardo/generate/actions/workflows/govulnchecks.yml/badge.svg?branch=master)](https://github.com/f-lombardo/generate/actions/workflows/govulnchecks.yml)

# generate

`generate` is a command-line interface (CLI) tool written in Go for quickly generating sample and test data such as IBAN
numbers, VAT numbers, UUIDs, and random passwords. It is designed to fit seamlessly into developer workflows by automatically copying
generated output directly to the system clipboard.

---

## 🚀 Features

- **IBAN Generation**: Generates valid IBAN numbers for various supported countries (e.g., Italy `IT`, Spain `ES`,
  Netherlands `NL`, etc.).
- **VAT Generation**: Generates valid VAT numbers for various supported countries (e.g., Italy `IT`, Spain `ES`,
  Netherlands `NL`, etc.). **Only IT is supported at this time.**
- **UUID Generation**: Supports generating both Version 4 (random) and Version 7 (time-ordered) UUIDs.
- **Password Generation**: Generates random passwords with customizable length.
- **Clipboard Integration**: Automatically copies generated values to the system clipboard (can be disabled via flag).
- **Flexible Output Formats**: Supports plain text (`tab`) and JSON (`json`) output.

---

## 📦 Installation & Build

### Installing with Homebrew

```bash
# 1. Add repository (Tap)
brew tap f-lombardo/tools

# 2. Trust the generate formula (needed by Homebrew 6.0+)
brew trust --formula f-lombardo/tools/generate

# 3. Install the program
brew install generate
```


### Building from Source

([Go](https://golang.org/) 1.27 or higher is required.)

Clone the repository and build the binary using:

```bash
go build -o generate .
```

You can move the compiled binary into a directory in your `$PATH` (e.g., `/usr/local/bin` or `~/bin`) to run it from
anywhere:

```bash
mv generate /usr/local/bin/
```

### Installing with go install

```bash
go install github.com/f-lombardo/generate
```

## 🛠 Usage

General syntax:

```bash
generate [global options] <command> [command options]
```

### Global Options

| Option       | Type     | Default | Description                                                                              |
|--------------|----------|---------|------------------------------------------------------------------------------------------|
| `-clipboard` | `bool`   | `true`  | Automatically copies results to the system clipboard. Use `-clipboard=false` to disable. |
| `-output`    | `string` | `tab`   | Output format. Supported values: `tab` (plain text) or `json`.                           |
| `-version`   |          |         | Prints current version and exits.                                                        |

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

### 3. `password`

Generates a random password.

**Command Options:**

- `-length` *(int, default: `14`)*: Length of the password (minimum `4`).

**Examples:**

```bash
# Generate a random password of default length (14)
generate password

# Generate a password with custom length of 20
generate password -length 20

# Generate a password in JSON format
generate -output json password
```

---

### 4. `vat`

Generates a syntactically valid VAT number.

**Command Options:**

- `-country` *(string, default: `IT`)*: ISO 3166-1 alpha-2 country code (e.g., `IT`, `ES`, `DE`, `FR`, `NL`).
**Only IT is supported at this tim.e**

**Examples:**

```bash
# Generate an Italian VAT number and copy it to clipboard
generate VAT

# Generate an Italian VAT number and copy it to clipboard
generate iban -country IT

# Format output as JSON
generate -output json vat -country IT
```

## 💡 Advanced Examples

- **Generate an IBAN without copying to clipboard:**
  ```bash
  generate -clipboard=false iban -country DE
  ```

- **JSON output for use in scripts or pipelines:**
  ```bash
  generate -output json uuid -version 7
  ```

- **Print application version:**
  ```bash
  generate -version
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
