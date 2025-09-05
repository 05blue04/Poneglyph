# Poneglyph API 🌎
> ⚠️ **Note**: This project is currently in active development. Features and API endpoints may change frequently.

A RESTful API for the One Piece universe built with Go and PostgreSQL. Explore characters, devil fruits, and pirate crews with advanced search, filtering, and pagination capabilities.

## Features ✨

- REST API with OpenAPI 3.0 specification
- IP-based rate limiting for fair usage and abuse prevention  
- Endpoints for characters, devil fruits, and crews with full CRUD support
- Built-in pagination, filtering, and sorting
- Validated input with clear error responses
- Custom API key authentication

## Architecture 🌇

- **Language**: Go 1.25
- **Database**: PostgreSQL with full-text search capabilities
- **Migration**: Goose migration tool
- **Authentication**: Custom API key middleware
- **Rate Limiting**: Token bucket algorithm with IP tracking
- **Validation**: Custom validation package with comprehensive rules

## Development 🚀

### Prerequisites
- Go 1.25+
- PostgreSQL 12+
- Goose migration tool

### Local Setup
```bash
# Clone the repository
git clone https://github.com/05blue04/Poneglyph.git
cd Poneglyph

# Set up environment variables
cp .env.example .env
# Edit .env with your database credentials

# Run migrations
goose postgres $DSN up

# Start the server
go run ./cmd/api
```

## Endpoints

- `/characters` - Search and manage One Piece characters
- `/devilfruits` - Browse and manage devil fruits  
- `/crews` - Explore pirate crews and members
- `/healthcheck` - API health status

**Full API documentation**: [OpenAPI Specification](./docs/api.yml)

## Roadmap 🗺️

### v2.0 (Future Enhancements)
- Database normalization improvements
  - Remove redundant captain_name field from crews table
  - Remove redundant current_owner field from devilfruits table
- New entity relationships
  - Create Locations table (islands, seas, regions)
  - Create Episodes/Chapters tracking

### Technical Improvements
- Comprehensive test suite
  - Unit tests for all models and handlers
  - Integration tests for database operations
  - End-to-end API testing
  - Performance benchmarking
- Enhanced bounty handling
  - Consistent null bounty value handling
  - Bounty history tracking
  - Improved format validation

## Contributing 🤝

Suggestions and feedback are welcome through GitHub issues!

## License 📄

This project is licensed under the MIT License 

## Acknowledgments

- Eiichiro Oda for creating the amazing One Piece universe
- Nico Robin the goat
- Monkey D. Luffy