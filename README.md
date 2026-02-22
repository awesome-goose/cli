# Goose CLI

A command-line tool for scaffolding Goose framework applications and modules.

## Installation

```bash
# Build from source
go build -o goose .

# Or install globally
go install .
```

## Commands

### Create a New Application

Create a new Goose application from a template:

```bash
goose app --name=myapp --template=<type>
```

Available templates:

- `api` - RESTful API application with database support, sample jobs and queries
- `cli` - Command-line application
- `web` - Web application with HTML templates and database support

#### Examples

```bash
# Create an API application
goose app --name=myapi --template=api

# Create a CLI application
goose app --name=mycli --template=cli

# Create a Web application
goose app --name=myweb --template=web

# Create in a specific directory
goose app --name=myapp --template=api --path=/path/to/projects
```

### Generate a Module

Generate a new module within an existing Goose application:

```bash
goose generate module --name=<name> --type=<type>
# Or use the short alias
goose g module --name=<name> --type=<type>
```

Module types:

- `plain` - Simple module with controller, service, routes, and DTOs
- `resource` - Full CRUD module with entity, controller, service, routes, and DTOs

The `--template` flag is optional; the CLI will auto-detect the app type from `main.go`.

#### Examples

```bash
# Generate a plain module (auto-detects app type)
goose g module --name=auth --type=plain

# Generate a resource module
goose g module --name=product --type=resource

# Specify the template type explicitly
goose g module --name=user --type=resource --template=api
```

## Application Structure

### API Application

```
myapi/
├── .env
├── go.mod
├── main.go
└── app/
    ├── app.module.go
    ├── app.controller.go
    ├── app.service.go
    ├── app.routes.go
    ├── app.dtos.go
    ├── jobs/
    │   └── sample.job.go
    └── queries/
        └── sample.query.go
```

### CLI Application

```
mycli/
├── .env
├── go.mod
├── main.go
└── app/
    ├── app.module.go
    ├── app.controller.go
    ├── app.service.go
    ├── app.routes.go
    └── app.dtos.go
```

### Web Application

```
myweb/
├── .env
├── go.mod
├── main.go
└── app/
    ├── app.module.go
    ├── app.controller.go
    ├── app.service.go
    ├── app.routes.go
    ├── app.dtos.go
    └── templates/
        ├── base/
        │   └── layout.html
        ├── pages/
        │   └── home.html
        └── partials/
            ├── header.html
            └── footer.html
```

## Module Structure

### Plain Module

```
modulename/
├── modulename.module.go
├── modulename.controller.go
├── modulename.service.go
├── modulename.routes.go
└── modulename.dtos.go
```

### Resource Module

```
modulename/
├── modulename.module.go
├── modulename.controller.go
├── modulename.service.go
├── modulename.routes.go
├── modulename.dtos.go
├── modulename.entity.go
└── templates/           # Web only
    └── pages/
        ├── list.html
        └── show.html
```

## After Creating an App

```bash
cd myapp
go mod tidy
go run main.go
```

## License

MIT
