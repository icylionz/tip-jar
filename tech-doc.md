# Tip Jar - Technical Architecture Document

## Project Overview

The Tip Jar is a self-hosted web application that enables groups of friends to maintain shared ledgers for tracking offenses and managing contributions/penalties. The application supports custom offense types, payment verification, and group management with role-based permissions.

## Technical Constraints

### Core Requirements
- **Self-Contained Binary**: The Go binary must include all necessary assets (templates, static files, migrations)
- **Single Binary Deployment**: No external file dependencies beyond the database
- **Go Version**: Go 1.25.1 (as specified in go.mod)
- **Database**: PostgreSQL as the primary data store
- **Self-Hosted**: Designed for personal/group hosting, not SaaS

### Performance Constraints
- Support for multiple concurrent tip jars (estimated 10-100 active jars)
- Responsive web interface for mobile and desktop
- File upload handling for payment proofs (images, receipts)
- Real-time notifications for offense reporting

## Technology Stack

### Backend Framework
- **Go 1.25.1** with standard library HTTP server
- **Echo Web Framework** - High performance, extensible, minimalist Go web framework
- **Air** - Live reload for development

### Database Layer
- **PostgreSQL 15+** - Primary database
- **golang-migrate/migrate** - Database migrations
- **pgx/v5** - PostgreSQL driver and toolkit
- **sqlc** - Type-safe SQL code generation

### Authentication & Authorization
- **Google OAuth 2.0** - Single OAuth provider for authentication
- **OAuth token-based sessions** - Use Google's OAuth tokens directly
- **Session-based authentication** with secure cookies

### Frontend & Templates
- **templ** - Type-safe HTML templating language that generates Go code
- **Embedded static assets** using Go embed
- **Tailwind CSS** - Utility-first CSS framework
- **Alpine.js** - Minimal JavaScript framework for interactivity
- **HTMX** - HTML over HTTP for dynamic interactions

### File Handling
- **Local file storage** with configurable directory
- **Image processing** for payment proof uploads
- **Disintegration/imaging** - Image resizing and optimization

### Configuration & Deployment
- **Viper** - Configuration management
- **Docker** support for containerized deployment
- **Environment variables** for sensitive configuration
- **Embedded migrations** and static assets

## System Architecture

### Application Structure
```
tipjar/
├── cmd/
│   └── server/           # Application entrypoint
├── internal/
│   ├── auth/            # Google OAuth and session handling
│   ├── config/          # Configuration management
│   ├── database/        # Database connection and migrations
│   ├── handlers/        # HTTP handlers and middleware
│   ├── models/          # Data models and business logic
│   ├── services/        # Business logic services
│   └── templates/       # Templ templates (.templ files)
├── migrations/          # Database migration files
├── static/             # CSS, JS, images (embedded)
├── uploads/            # File upload directory
└── docker/             # Docker configuration
```

### Database Schema Design

The database will support the core entities needed for tip jar functionality:
- **Users** - OAuth-authenticated users with profile information
- **Tip Jars** - Groups/communities with invite codes and admin roles
- **Jar Memberships** - User participation in tip jars with role management
- **Offense Types** - Customizable offense definitions with various cost types
- **Offenses** - Reported incidents with status tracking
- **Payments** - Settlement records with optional proof verification

## Security Considerations

### Authentication Security
- Secure Google OAuth 2.0 implementation with state parameter
- Google OAuth tokens used directly for session management
- Secure session cookies with HttpOnly and Secure flags
- CSRF protection for state-changing operations

### Data Protection
- Input validation and sanitization
- SQL injection prevention using parameterized queries
- File upload restrictions (type, size, scan for malware)
- Rate limiting on API endpoints

### Access Control
- Role-based permissions (jar admin vs member)
- User can only access jars they're members of
- Payment verification requires jar membership
- Offense reporting restricted to jar members

## Development Workflow

### Local Development Setup
1. Install Go 1.25.1
2. Install PostgreSQL
3. Install Air for live reload: `go install github.com/cosmtrek/air@latest`
4. Install migrate CLI: `go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest`
5. Install templ: `go install github.com/a-h/templ/cmd/templ@latest`
6. Run migrations: `migrate -path migrations -database $DATABASE_URL up`
7. Generate templ files: `templ generate`
8. Start development server: `air`

### Code Generation
- Use `sqlc generate` to generate type-safe database code
- Use `templ generate` to generate Go code from templ templates
- Use `go generate` for embedding static assets

### Testing Strategy
- Unit tests for business logic services
- Integration tests for database operations
- End-to-end tests for critical user flows
- Use `testcontainers-go` for database testing

## Performance Considerations

### Database Optimization
- Proper indexing on foreign keys and frequently queried columns
- Connection pooling with pgxpool
- Query optimization for jar member lookups
- Pagination for large offense lists

### Caching Strategy
- In-memory caching for user sessions
- Static asset caching with appropriate headers
- Database query result caching for jar metadata

### File Handling
- Image compression for uploaded payment proofs
- Configurable file size limits
- Background cleanup of orphaned files

## Monitoring & Observability

### Logging
- Structured logging with slog (Go 1.21+)
- Request/response logging middleware
- Error tracking and alerting

### Metrics
- HTTP request metrics (response time, status codes)
- Database connection pool metrics
- File upload metrics

### Health Checks
- Database connectivity check
- File system write permissions check
- OAuth provider connectivity

## Future Considerations

### Scalability
- Database read replicas for reporting
- CDN for static assets
- Background job processing for notifications

### Features
- Mobile app using the same backend API
- Webhook integrations for external notifications
- Data export functionality (CSV, JSON)
- Advanced analytics and reporting

This architecture provides a solid foundation for the Tip Jar application while maintaining simplicity and ensuring all requirements are met through modern Go best practices.