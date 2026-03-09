<p align="center">
  <strong>🪿 Goose Framework</strong><br>
  <sub>Modular • Scalable • Multi-Platform</sub>
</p>

---

# Goose CLI

A command-line tool for scaffolding Goose framework applications and modules.

---

## Table of Contents

- [Installation](#installation)
  - [Quick Install](#quick-install)
  - [Building from Source](#building-from-source)
  - [Verifying Downloads](#verifying-downloads)
  - [Available Platforms](#available-platforms)
- [Quick Start](#quick-start)
- [Commands](#commands)
  - [version](#version)
  - [app](#app)
  - [generate module / g module](#generate-module--g-module)
- [Application Templates](#application-templates)
  - [API Template](#api-template)
  - [CLI Template](#cli-template)
  - [Web Template](#web-template)
  - [Multi-Platform Template](#multi-platform-template)
- [Module Types](#module-types)
  - [Plain Module](#plain-module)
  - [Resource Module](#resource-module)
- [Examples](#examples)
- [Project Structure](#project-structure)
- [Environment Variables](#environment-variables)
- [Troubleshooting](#troubleshooting)
- [Updating & Uninstalling](#updating--uninstalling)
- [License](#license)

---

## Installation

### Quick Install

#### Linux (amd64)

```bash
curl -LO https://github.com/awesome-goose/goose-cli/releases/latest/download/goose-linux-amd64.tar.gz
tar -xzf goose-linux-amd64.tar.gz
sudo mv goose-* /usr/local/bin/goose
sudo chmod +x /usr/local/bin/goose
goose version
```

#### macOS (Apple Silicon - arm64)

```bash
curl -LO https://github.com/awesome-goose/goose-cli/releases/latest/download/goose-darwin-arm64.tar.gz
tar -xzf goose-darwin-arm64.tar.gz
sudo mv goose-* /usr/local/bin/goose
sudo chmod +x /usr/local/bin/goose
goose version
```

#### macOS (Intel - amd64)

```bash
curl -LO https://github.com/awesome-goose/goose-cli/releases/latest/download/goose-darwin-amd64.tar.gz
tar -xzf goose-darwin-amd64.tar.gz
sudo mv goose-* /usr/local/bin/goose
sudo chmod +x /usr/local/bin/goose
goose version
```

#### Windows (amd64)

Using PowerShell:

```powershell
Invoke-WebRequest -Uri "https://github.com/awesome-goose/goose-cli/releases/latest/download/goose-windows-amd64.zip" -OutFile "goose-windows-amd64.zip"
Expand-Archive -Path "goose-windows-amd64.zip" -DestinationPath "."
New-Item -ItemType Directory -Force -Path "C:\Program Files\goose"
Move-Item -Path "goose-*.exe" -Destination "C:\Program Files\goose\goose.exe"
[Environment]::SetEnvironmentVariable("Path", $env:Path + ";C:\Program Files\goose", "Machine")
# Open a new terminal, then verify:
goose version
```

### Building from Source

```bash
git clone https://github.com/awesome-goose/goose-cli.git
cd goose-cli/cli
go build -o goose .
sudo mv goose /usr/local/bin/
```

**Requirements:** Go 1.22+, Git

### Verifying Downloads

Each release includes a `checksums.txt` file with SHA256 hashes:

```bash
curl -LO https://github.com/awesome-goose/goose-cli/releases/latest/download/checksums.txt
sha256sum -c checksums.txt --ignore-missing
```

### Available Platforms

| OS      | Architecture  | Archive Name                          |
| ------- | ------------- | ------------------------------------- |
| Linux   | amd64 (x64)   | `goose-<version>-linux-amd64.tar.gz`  |
| Linux   | arm64         | `goose-<version>-linux-arm64.tar.gz`  |
| Linux   | 386 (x86)     | `goose-<version>-linux-386.tar.gz`    |
| Linux   | arm           | `goose-<version>-linux-arm.tar.gz`    |
| macOS   | amd64 (Intel) | `goose-<version>-darwin-amd64.tar.gz` |
| macOS   | arm64 (M1/M2) | `goose-<version>-darwin-arm64.tar.gz` |
| Windows | amd64 (x64)   | `goose-<version>-windows-amd64.zip`   |
| Windows | arm64         | `goose-<version>-windows-arm64.zip`   |
| Windows | 386 (x86)     | `goose-<version>-windows-386.zip`     |

---

## Quick Start

```bash
# Check CLI version
goose version

# Create a new API application
goose app --name=myapi --template=api

# Navigate to the new project and generate a module
cd myapi
goose g module --name=users --type=resource

# Run the application
go mod tidy
go run main.go
```

---

## Commands

| Command                                     | Description                                      |
| ------------------------------------------- | ------------------------------------------------ |
| `goose` or `goose version`                  | Display version and help information             |
| `goose app`                                 | Create a new Goose application                   |
| `goose generate module` or `goose g module` | Generate a new module in an existing application |

### version

Display CLI version and help information.

```bash
goose
goose version
```

**Output:**

```
┌──────────────────────────────────────────────────────────────────┐
│ Goose CLI                                                        │
├──────────────────────────────────────────────────────────────────┤
│ Version: 0.0.0                                                   │
│ A tool for scaffolding Goose applications                        │
│                                                                  │
│ Commands:                                                        │
│   goose app --name=<name> --template=<api|cli|web|multi>         │
│   goose g module --name=<name> --type=<plain|resource>           │
│   goose version                                                  │
│                                                                  │
│ Templates:                                                       │
│   api   - REST API application                                   │
│   cli   - Command-line interface application                     │
│   web   - Web application with HTML templates                    │
│   multi - Multi-platform app (API + Web + CLI in one)            │
└──────────────────────────────────────────────────────────────────┘
```

### app

Create a new Goose application from a template.

```bash
goose app --name=<name> --template=<template> [--path=<path>]
```

**Flags:**

| Flag         | Required | Description                                                                         |
| ------------ | -------- | ----------------------------------------------------------------------------------- |
| `--name`     | Yes      | Name of the application (used for directory and module name)                        |
| `--template` | Yes      | Application template type: `api`, `cli`, `web`, or `multi`                          |
| `--path`     | No       | Custom path where the application should be created (defaults to current directory) |

**Examples:**

```bash
# Create an API application
goose app --name=myapi --template=api

# Create a CLI application
goose app --name=mycli --template=cli

# Create a Web application
goose app --name=myweb --template=web

# Create a multi-platform application (API + Web + CLI)
goose app --name=myapp --template=multi

# Create in a specific directory
goose app --name=myapp --template=api --path=/path/to/projects
```

### generate module / g module

Generate a new module within an existing Goose application. Run from the project root directory.

```bash
goose generate module --name=<name> [--type=<type>] [--template=<template>]
goose g module --name=<name> [--type=<type>] [--template=<template>]
```

**Flags:**

| Flag         | Required | Default       | Description                                              |
| ------------ | -------- | ------------- | -------------------------------------------------------- |
| `--name`     | Yes      | -             | Name of the module (e.g., `users`, `products`, `orders`) |
| `--type`     | No       | `plain`       | Module type: `plain` or `resource`                       |
| `--template` | No       | Auto-detected | Target platform template: `api`, `cli`, or `web`         |

**Examples:**

```bash
# Generate a plain module (auto-detects app type)
goose g module --name=auth --type=plain

# Generate a resource module (with CRUD operations)
goose g module --name=product --type=resource

# Specify the template type explicitly
goose g module --name=user --type=resource --template=api
```

---

## Application Templates

### API Template

REST API application optimized for JSON-based services.

```bash
goose app --name=myapi --template=api
```

**Structure:**

```
myapi/
├── .env
├── go.mod
├── main.go
└── app/
    ├── app.controller.go
    ├── app.dtos.go
    ├── app.module.go
    ├── app.routes.go
    ├── app.service.go
    ├── jobs/
    │   └── sample.job.go
    └── queries/
        └── sample.query.go
```

**Features:** HTTP server with JSON responses, structured routing, background job support, query handling, environment configuration.

### CLI Template

Command-line interface application.

```bash
goose app --name=mycli --template=cli
```

**Structure:**

```
mycli/
├── .env
├── go.mod
├── main.go
└── app/
    ├── app.controller.go
    ├── app.dtos.go
    ├── app.module.go
    ├── app.routes.go
    └── app.service.go
```

**Features:** Command parsing with flags, structured command handlers, console output formatting, environment configuration.

### Web Template

Web application with HTML template rendering.

```bash
goose app --name=myweb --template=web
```

**Structure:**

```
myweb/
├── .env
├── go.mod
├── main.go
└── app/
    ├── app.controller.go
    ├── app.dtos.go
    ├── app.module.go
    ├── app.routes.go
    ├── app.service.go
    └── templates/
        ├── base/
        │   └── layout.html
        ├── pages/
        │   └── home.html
        └── partials/
            ├── header.html
            └── footer.html
```

**Features:** HTTP server with HTML responses, template rendering, static file serving, view helpers, environment configuration.

### Multi-Platform Template

Unified application supporting API, Web, and CLI interfaces.

```bash
goose app --name=myapp --template=multi
```

**Structure:**

```
myapp/
├── .env
├── go.mod
├── main.go
└── app/
    ├── api/
    │   ├── api.controller.go
    │   ├── api.module.go
    │   ├── api.routes.go
    │   └── api.service.go
    ├── cli/
    │   ├── cli.controller.go
    │   ├── cli.module.go
    │   ├── cli.routes.go
    │   └── cli.service.go
    ├── web/
    │   ├── web.controller.go
    │   ├── web.module.go
    │   ├── web.routes.go
    │   ├── web.service.go
    │   └── templates/
    └── shared/
        └── (shared services and utilities)
```

**Features:** Single codebase for multiple interfaces, shared business logic, platform-specific controllers, unified configuration.

---

## Module Types

### Plain Module

Basic module structure for general functionality.

```bash
goose g module --name=notifications --type=plain
```

**Generated Files:**

| File            | Purpose                                          |
| --------------- | ------------------------------------------------ |
| `controller.go` | Request handlers and business logic coordination |
| `service.go`    | Business logic implementation                    |
| `routes.go`     | Route definitions                                |
| `dtos.go`       | Data transfer objects for request/response       |
| `module.go`     | Module registration and dependency injection     |

**Structure:**

```
app/
└── notifications/
    ├── notifications.controller.go
    ├── notifications.dtos.go
    ├── notifications.module.go
    ├── notifications.routes.go
    └── notifications.service.go
```

### Resource Module

Module with full CRUD support and database entity.

```bash
goose g module --name=products --type=resource
```

**Generated Files:** All plain module files plus:

| File        | Purpose                   |
| ----------- | ------------------------- |
| `entity.go` | Database model definition |

**Structure:**

```
app/
└── products/
    ├── products.controller.go
    ├── products.dtos.go
    ├── products.entity.go
    ├── products.module.go
    ├── products.routes.go
    └── products.service.go
```

---

## Examples

### Building a REST API

```bash
goose app --name=store-api --template=api
cd store-api

# Generate resource modules
goose g module --name=products --type=resource
goose g module --name=categories --type=resource
goose g module --name=orders --type=resource

# Generate plain modules for utilities
goose g module --name=auth --type=plain
goose g module --name=notifications --type=plain

go run main.go
```

### Building a CLI Tool

```bash
goose app --name=devtools --template=cli
cd devtools

goose g module --name=migrate --type=plain
goose g module --name=seed --type=plain
goose g module --name=backup --type=plain

go build -o devtools .
./devtools migrate
```

### Building a Full-Stack Application

```bash
goose app --name=myapp --template=multi
cd myapp

# Generate shared resource modules
goose g module --name=users --type=resource --template=api
goose g module --name=posts --type=resource --template=api

# Generate web-specific modules
goose g module --name=dashboard --type=plain --template=web
goose g module --name=admin --type=plain --template=web

go run main.go
```

---

## Project Structure

After creating an application and adding modules:

```
myproject/
├── .env                    # Environment configuration
├── go.mod                  # Go module definition
├── go.sum                  # Dependency checksums
├── main.go                 # Application entry point
└── app/
    ├── app.controller.go   # Root controller
    ├── app.dtos.go         # Root DTOs
    ├── app.module.go       # Root module (registers sub-modules)
    ├── app.routes.go       # Root routes
    ├── app.service.go      # Root service
    │
    ├── users/              # Generated module
    │   ├── users.controller.go
    │   ├── users.dtos.go
    │   ├── users.entity.go
    │   ├── users.module.go
    │   ├── users.routes.go
    │   └── users.service.go
    │
    └── products/           # Another generated module
        ├── products.controller.go
        └── ...
```

---

## Environment Variables

Generated applications include a `.env` file with common configuration options. Customize as needed for your environment.

---

## Troubleshooting

### "command not found" after installation

Make sure the installation directory is in your PATH:

```bash
which goose

# If not found, add to ~/.bashrc, ~/.zshrc, or equivalent:
export PATH="$PATH:/usr/local/bin"
```

### Permission denied on Linux/macOS

```bash
chmod +x /usr/local/bin/goose
```

### macOS security warning ("unidentified developer")

Option 1: Go to **System Preferences > Security & Privacy > General** and click **"Allow Anyway"**.

Option 2: Run in terminal:

```bash
xattr -d com.apple.quarantine /usr/local/bin/goose
```

---

## Updating & Uninstalling

### Updating

Download and install the new version following the installation steps. The new binary will replace the old one.

### Uninstalling

**Linux/macOS:**

```bash
sudo rm /usr/local/bin/goose
```

**Windows:**

```powershell
Remove-Item "C:\Program Files\goose" -Recurse -Force
# Remove "C:\Program Files\goose" from your PATH environment variable
```

---

## Running Tests

```bash
# Run all tests
go test ./tests/...

# Run with verbose output
go test ./tests/... -v
```

---

## Code Coverage

```bash
# Coverage for all goose packages
go test ./tests/... -coverpkg=./...

```

---

## Tips

1. **Naming Conventions:** Use lowercase, singular names for modules (e.g., `user` not `Users`)
2. **Resource vs Plain:** Use `resource` type when you need database CRUD operations
3. **Multi-Platform:** For shared code in multi-platform apps, place it in the `shared/` directory
4. **Run from Project Root:** Module generation commands should be run from the project root directory

---

## License

MIT
