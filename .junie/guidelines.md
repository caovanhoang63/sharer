# Sharer - Project Guidelines

## Project Overview

**Sharer** is a Go web application that allows users to share HTML content through unique URLs. Users can create, categorize, and share HTML pages that are accessible via generated slugs. The application provides both a web interface and REST API endpoints for content management.

### Key Features
- Share HTML content via unique URLs
- Categorize shared content for organization
- Modern, responsive web UI with HTMX interactivity
- Beautiful styling with Tailwind CSS and DaisyUI components
- REST API for programmatic access
- SQLite database for data persistence
- GORM for database query
- Type-safe HTML templating with Templ

## Architecture & Structure

The project follows **Clean Architecture** principles with clear separation of concerns:

```
sharer/
├── internal/
│   ├── database/           # Database configuration and migrations
│   └── modules/
│       ├── category/       # Category management module
│       │   ├── controller.go
│       │   ├── interfaces.go
│       │   ├── models.go
│       │   ├── repository.go
│       │   └── service.go
│       ├── page/          # Page/content sharing module
│       │   ├── controller.go
│       │   ├── interfaces.go
│       │   ├── models.go
│       │   ├── repository.go
│       │   └── service.go
│       └── user/          # User models (future expansion)
├── views/                 # Templ templates
│   ├── components/        # Reusable UI components
│   ├── layouts/          # Page layouts
│   └── pages/            # Page templates
├── templates/            # Legacy HTML templates
├── main.go              # Application entry point
├── go.mod               # Go module definition
└── Dockerfile           # Container configuration
```

### Layer Responsibilities
- **Controllers**: Handle HTTP requests/responses, input validation
- **Services**: Business logic implementation
- **Repositories**: Data access layer, database operations
- **Models**: Data structures and DTOs

## Technology Stack

- **Language**: Go 1.24
- **Web Framework**: Gin (HTTP router and middleware)
- **Database**: SQLite with GORM ORM
- **Templating**: Templ (type-safe HTML templates)
- **Frontend**: HTMX + Tailwind CSS + DaisyUI
- **Containerization**: Docker
- **Architecture**: Clean Architecture with Repository pattern

### Key Dependencies
- `github.com/gin-gonic/gin` - HTTP web framework
- `gorm.io/gorm` - ORM for database operations
- `gorm.io/driver/sqlite` - SQLite database driver
- `github.com/a-h/templ` - Type-safe HTML templating

### Frontend Technologies
- **HTMX** - Modern HTML-driven interactivity (loaded from CDN)
- **Tailwind CSS** - Utility-first CSS framework (loaded from CDN)
- **DaisyUI** - Tailwind CSS component library (loaded from CDN)

## Development Workflow

### Prerequisites
- Go 1.24 or later
- SQLite (for local development)
- Templ CLI (for template generation)

### Setup
1. Clone the repository
2. Install dependencies: `go mod download`
3. Install Templ CLI: `go install github.com/a-h/templ/cmd/templ@latest`
4. Generate templates: `templ generate`
5. Run the application: `go run main.go`

### Template Development
- Templates are written in `.templ` files
- Run `templ generate` to compile templates to Go code
- Generated `*_templ.go` files should not be manually edited

## Testing Approach

### Test Strategy
- **Unit Tests**: Test individual functions and methods in isolation
- **Integration Tests**: Test module interactions and database operations
- **End-to-End Tests**: Test complete user workflows via HTTP endpoints

### Running Tests
```bash
# Run all tests
go test ./...

# Run tests with coverage
go test -cover ./...

# Run tests for specific module
go test ./internal/modules/page/...
```

### Test Guidelines
- Write tests for all business logic in services
- Test repository methods with database interactions
- Mock external dependencies in unit tests
- Use table-driven tests for multiple test cases

## Build & Deployment

### Local Development
```bash
# Run with hot reload (if using air)
air

# Build binary
go build -o sharer main.go

# Run binary
./sharer
```

### Docker Deployment
```bash
# Build Docker image
docker build -t sharer .

# Run container
docker run -p 8080:8080 sharer
```

### Production Considerations
- Application runs on port 8080 by default
- SQLite database file: `./sharer.db`
- Gin runs in release mode for production
- Database migrations run automatically on startup

## Code Style Guidelines

### General Principles
- Follow standard Go conventions and idioms
- Use `gofmt` for code formatting
- Follow Clean Architecture principles
- Maintain clear separation between layers

### Naming Conventions
- Use descriptive names for variables and functions
- Follow Go naming conventions (PascalCase for exported, camelCase for unexported)
- Use consistent naming across modules

### Error Handling
- Always handle errors explicitly
- Use meaningful error messages
- Log errors appropriately at service layer
- Return appropriate HTTP status codes in controllers

### Database Operations
- Use GORM best practices
- Implement proper indexing for performance
- Use transactions for multi-step operations
- Handle soft deletes consistently

### API Design
- Follow RESTful conventions
- Use consistent JSON response formats
- Implement proper input validation
- Return appropriate HTTP status codes

### Template Guidelines
- Keep templates focused and reusable
- Use components for common UI elements
- Follow consistent naming for template files
- Generate templates after changes: `templ generate`
- Use DaisyUI component classes for consistent styling (e.g., `btn`, `card`, `form-control`)
- Leverage Tailwind utility classes for custom styling and layout
- Implement HTMX attributes for dynamic interactions (`hx-post`, `hx-get`, `hx-target`, etc.)
- Use semantic HTML with appropriate ARIA attributes for accessibility

## Development Notes for Junie

- **Always run tests** after making changes to ensure correctness
- **Build the project** before submitting to verify compilation
- **Generate templates** if modifying `.templ` files: `templ generate`
- **Follow the existing architecture** - don't break layer boundaries
- **Maintain database consistency** - consider migrations for schema changes
- **Test both web UI and API endpoints** when making changes
- **Use Context7 MCP** for reading library documentation when working with external libraries or frameworks
- **Use Exa for web search** when links are provided in prompts to gather additional context and information
- **Just Use Exa when Context7 not found the document**
