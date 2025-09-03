# Poneglyph API
> ⚠️ **Note**: This project is currently in active development. Features and API endpoints may change frequently.

A RESTful API for the One Piece universe built with Go and PostgreSQL. Explore characters, devil fruits, and pirate crews with advanced search, filtering, and pagination capabilities.

## Current Features

### Characters
- Full CRUD operations for One Piece characters
- Advanced filtering by race, age, bounty, and origin
- Full-text search across names and descriptions
- Comprehensive bounty and race validation

### Devil Fruits
- Complete devil fruit database with type classification (Paramecia, Zoan, Logia)
- Ownership tracking with current and previous owners
- Automatic character relationship management

### Pirate Crews
- Crew management with captain assignment
- Member addition/removal with automatic bounty calculations
- Crew member listing with filtering capabilities

## API Endpoints

### Characters
- `GET /v1/characters` - List all characters
- `GET /v1/characters/{id}` - Get a specific character
- `POST /v1/characters` - Create a new character
- `PATCH /v1/characters/{id}` - Update an existing character
- `DELETE /v1/characters/{id}` - Delete a character

### Devil Fruits
- `GET /v1/devilfruits` - List all devil fruits
- `GET /v1/devilfruits/{id}` - Get a specific devil fruit
- `POST /v1/devilfruits` - Create a new devil fruit
- `PATCH /v1/devilfruits/{id}` - Update an existing devil fruit
- `DELETE /v1/devilfruits/{id}` - Delete a devil fruit

### Crews
- `GET /v1/crews/{id}` - Get a specific crew
- `POST /v1/crews` - Create a new crew
- `PATCH /v1/crews/{id}` - Update an existing crew
- `DELETE /v1/crews/{id}` - Delete a crew
- `GET /v1/crews/{id}/members` - List crew members
- `POST /v1/crews/{id}/members` - Add crew member
- `DELETE /v1/crews/{crew_id}/members/{character_id}` - Remove crew member

### Query Parameters
Most list endpoints support:
- `search` - Full-text search across names and descriptions
- `page` - Page number for pagination (default: 1)
- `page_size` - Results per page (default: 20, max: 100)
- `sort` - Sort by field. Use `-` prefix for descending order

#### Characters specific:
- `race` - Filter by character race (human, fishman, mink, giant, etc.)
- `age` - Minimum age filter
- `bounty` - Minimum bounty filter
- `origin` - Filter by character's origin location

#### Devil Fruits specific:
- `type` - Filter by devil fruit type (paramecia, zoan, logia)

#### Crew Members specific:
- `bounty` - Minimum bounty filter for crew members

## Project Backlog

### v1 (Current Release Goals)
- [ ] **Implement listCrewHandler**
  - Allow filter by bounty
  - Allow filter by captain
  - Search parameter to search name and descriptions
- [ ] **Implement API key authentication**
  - Create API keys table and management
  - Add middleware for authentication
  - Implement read/write permission levels
- [ ] **Docker containerization**
  - Create Dockerfile and docker-compose.yml
  - Set up development and production configurations
  - Include database initialization scripts
- [ ] **Deploy v1 to production**
  - Choose hosting platform (Railway/Render/DigitalOcean)
  - Set up CI/CD pipeline
  - Configure environment-based settings

### v2 (Future Enhancements)
- [ ] **Database normalization improvements**
  - Remove redundant captain_name field from crews table
  - Remove redundant current_owner field from devilfruits table
- [ ] **New entity relationships**
  - Create Locations table (islands, seas, regions)
  - Create Episodes / Chapters 

### Technical Debt
- [ ] **Add comprehensive test suite**
  - Unit tests for all models and handlers
  - Integration tests for database operations
  - End-to-end API testing
  - Performance benchmarking tests
- [ ] **Documentation enhancements**
  - Complete OpenAPI/Swagger documentation

### Bugs & Fixes
- [ ] **Improve bounty handling**
  - Handle null bounty values more consistently
  - Add bounty history tracking
  - Validate bounty format and reasonable limits

## Architecture

- **Language**: Go 1.25+
- **Database**: PostgreSQL with full-text search
- **Migrations**: Goose migration tool
- **Validation**: Custom validation package
- **Testing**: Go's built-in testing framework

## Contributing

This project is primarily for learning and portfolio purposes, but suggestions and feedback are welcome through GitHub issues.

## License

This project is licensed under the MIT License - see the LICENSE file for details.

## Acknowledgments

- Eiichiro Oda for creating the amazing One Piece universe
- Nico Robin the goat
- Monkey D. Luffy