# Tip Jar

A self-hosted web application that enables groups of friends to maintain shared ledgers for tracking offenses and managing contributions/penalties.

## Project Status

🚧 **In Development** - This is the initial project skeleton. Core features are being implemented according to the user stories.

## Quick Start

### Prerequisites

- Go 1.25.1+
- PostgreSQL 15+

### Development Setup

1. **Clone and setup**:
   ```bash
   git clone <repository>
   cd tipjar
   make setup
   ```

2. **Configure environment**:
   ```bash
   cp .env.example .env
   # Edit .env with your Google OAuth credentials
   ```

3. **Start database**:
   ```bash
   make docker-dev
   ```

4. **Run migrations**:
   ```bash
   make migration-up
   ```

5. **Generate templates**:
   ```bash
   make templ
   ```

6. **Start development server**:
   ```bash
   make dev
   ```

### Google OAuth Setup

1. Go to [Google Cloud Console](https://console.cloud.google.com/)
2. Create a new project or select existing one
3. Enable Google+ API
4. Create OAuth 2.0 credentials
5. Add `http://localhost:8080/auth/callback` to authorized redirect URIs
6. Update `.env` with your client ID and secret

## Project Structure

```
tipjar/
├── cmd/server/          # Application entry point
├── internal/
│   ├── auth/           # Google OAuth and session handling
│   ├── config/         # Configuration management  
│   ├── database/       # Database connection and migrations
│   ├── handlers/       # HTTP handlers and middleware
│   ├── models/         # Data models and business logic
│   ├── services/       # Business logic services
│   └── templates/      # Templ templates (.templ files)
├── migrations/         # Database migrations
├── static/            # CSS, JS, images (embedded)
├── uploads/           # File upload directory
└── docker/            # Docker configuration
```

## Technology Stack

- **Backend**: Go 1.25.1, Echo framework, PostgreSQL
- **Frontend**: Templ templates, Tailwind CSS (CDN), Alpine.js (CDN), HTMX (CDN)
- **Auth**: Google OAuth 2.0
- **Database**: PostgreSQL with migrations
- **Deployment**: Self-contained binary with embedded assets

## Available Commands

```bash
make build          # Build the application
make run            # Run the built application  
make dev            # Development mode with live reload
make test           # Run tests
make clean          # Clean build artifacts

make migration-up   # Run database migrations
make migration-down # Rollback migrations
make migration-create name=description # Create new migration

make docker-build   # Build Docker image
make docker-run     # Run with docker-compose
make docker-dev     # Start only PostgreSQL

make templ          # Generate templ files
make sqlc           # Generate sqlc files
make generate       # Generate all code (templ + sqlc)
make deps           # Install development dependencies
make setup          # Complete development setup

make build-prod     # Production build with optimizations
make lint           # Lint code
make fmt            # Format code
```

## Frontend Architecture

The application uses modern web technologies served from CDN for optimal performance:

- **Tailwind CSS**: Loaded from CDN with custom configuration
- **Alpine.js**: Reactive JavaScript framework for interactivity
- **HTMX**: HTML-over-HTTP for dynamic server interactions
- **Custom CSS**: Application-specific styles in `/static/css/app.css`
- **Custom JS**: Utility functions and app logic in `/static/js/app.js`

No build step required for frontend assets - everything loads directly from CDN or static files.

## User Stories Implementation

The application is being built according to the user stories in `stories.md`, organized into iterations:

- **Iteration 1**: Core MVP - Basic functional tip jar ✏️ *In Progress*
- **Iteration 2**: Custom Offenses - Personalized offense types 📋 *Planned*
- **Iteration 3**: Payment Verification - Proof and verification 📋 *Planned*
- **Iteration 4**: Personal Dashboard - Cross-jar views 📋 *Planned*
- **Iteration 5**: Timeline & Filtering - Better navigation 📋 *Planned*
- **Iteration 6**: Disputes & Moderation - Fair resolution 📋 *Planned*
- **Iteration 7**: Notifications - Stay informed 📋 *Planned*
- **Iteration 8**: Admin & Analytics - Powerful tools 📋 *Planned*

## Development Workflow

### Template Development
1. Create `.templ` files in `internal/templates/`
2. Run `make templ` to generate Go code
3. Import and use in handlers

### Database Changes
1. Create migration: `make migration-create name=add_new_feature`
2. Edit the generated `.sql` files
3. Apply migration: `make migration-up`

### Static Assets
- CSS: Edit `static/css/app.css` directly
- JS: Edit `static/js/app.js` directly  
- Images: Add to `static/images/`
- All static assets are embedded in the binary

## Contributing

1. Follow the user stories in `stories.md`
2. Maintain the self-contained binary requirement
3. Use the existing technology stack (CDN-based frontend)
4. Write tests for new features
5. Update documentation

## License

[Add your license here]
