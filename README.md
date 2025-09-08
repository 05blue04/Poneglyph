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
- Runtime and API metrics including memory stats, request/response counts, and processing times

## Architecture 🌇

- **Language**: Go 1.25
- **Database**: PostgreSQL with full-text search capabilities
- **Migration**: Goose migration tool
- **Authentication**: Custom API key middleware
- **Rate Limiting**: Token bucket algorithm with IP tracking

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
- `/metrics` - API metrics

**Full API documentation**: [OpenAPI Specification](./docs/api.yml)

# API Examples

## Authentication

All write operations (POST, PATCH, DELETE) require authentication using a Bearer token:

```bash
curl -H "Authorization: Bearer your-token-here" \
  https://api.poneglyph.dev/v1/characters
```

## Characters

### Get All Characters
```bash
# Basic request
curl "https://api.poneglyph.dev/v1/characters"

# With pagination
curl "https://api.poneglyph.dev/v1/characters?page=1&page_size=10"

# Search and filter
curl "https://api.poneglyph.dev/v1/characters?search=luffy&race=human&bounty=1000000"
```

**Response:**
```json
{
  "characters": [
    {
      "id": 1,
      "name": "Monkey D. Luffy",
      "age": 19,
      "description": "Captain of the Straw Hat Pirates with rubber powers",
      "origin": "Foosha Village, East Blue",
      "bounty": 3000000000,
      "race": "human",
      "episode": 1
    }
  ],
  "metadata": {
    "current_page": 1,
    "page_size": 20,
    "first_page": 1,
    "last_page": 5,
    "total_records": 100
  }
}
```

### Create a Character
```bash
curl -X POST "https://api.poneglyph.dev/v1/characters" \
  -H "Authorization: Bearer your-token" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Roronoa Zoro",
    "age": 21,
    "description": "First mate and swordsman of the Straw Hat Pirates",
    "origin": "Shimotsuki Village, East Blue",
    "bounty": 1111000000,
    "race": "human",
    "episode": 2
  }'
```

### Update a Character
```bash
curl -X PATCH "https://api.poneglyph.dev/v1/characters/1" \
  -H "Authorization: Bearer your-token" \
  -H "Content-Type: application/json" \
  -d '{
    "bounty": 3500000000,
    "age": 20
  }'
```

## Devil Fruits

### Get Devil Fruits by Type
```bash
# Get all Logia fruits
curl "https://api.poneglyph.dev/v1/devilfruits?type=logia"

# Search for specific fruit
curl "https://api.poneglyph.dev/v1/devilfruits?search=gomu"
```

**Response:**
```json
{
  "devil_fruits": [
    {
      "id": 1,
      "name": "Gomu Gomu no Mi",
      "description": "Allows the user's body to stretch like rubber",
      "type": "paramecia",
      "current_owner": "Monkey D. Luffy",
      "character_id": 1,
      "previous_owners": [],
      "episode": 1
    }
  ],
  "metadata": {
    "current_page": 1,
    "page_size": 20,
    "first_page": 1,
    "last_page": 3,
    "total_records": 50
  }
}
```

### Create a Devil Fruit
```bash
curl -X POST "https://api.poneglyph.dev/v1/devilfruits" \
  -H "Authorization: Bearer your-token" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Mera Mera no Mi",
    "description": "Allows the user to create, control, and transform into fire",
    "type": "logia",
    "character_id": 2,
    "previous_owners": ["Portgas D. Ace"],
    "episode": 94
  }'
```

### Transfer Devil Fruit Ownership
```bash
curl -X PATCH "https://api.poneglyph.dev/v1/devilfruits/1" \
  -H "Authorization: Bearer your-token" \
  -H "Content-Type: application/json" \
  -d '{
    "character_id": 5,
    "previous_owners": ["Monkey D. Luffy"]
  }'
```

## Crews

### Get Crew Details
```bash
curl "https://api.poneglyph.dev/v1/crews/1"
```

**Response:**
```json
{
  "crew": {
    "id": 1,
    "name": "Straw Hat Pirates",
    "description": "A pirate crew led by Monkey D. Luffy",
    "ship_name": "Thousand Sunny",
    "captain_id": 1,
    "captain_name": "Monkey D. Luffy",
    "total_bounty": 8816000000,
    "member_count": 10
  }
}
```

### Create a Crew
```bash
curl -X POST "https://api.poneglyph.dev/v1/crews" \
  -H "Authorization: Bearer your-token" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Heart Pirates",
    "description": "A pirate crew led by Trafalgar Law",
    "ship_name": "Polar Tang",
    "captain_id": 5
  }'
```

### Manage Crew Members
```bash
# Add member to crew
curl -X POST "https://api.poneglyph.dev/v1/crews/1/members" \
  -H "Authorization: Bearer your-token" \
  -H "Content-Type: application/json" \
  -d '{"character_id": 3}'

# Get crew members
curl "https://api.poneglyph.dev/v1/crews/1/members?bounty=100000000"

# Remove member from crew
curl -X DELETE "https://api.poneglyph.dev/v1/crews/1/members/3" \
  -H "Authorization: Bearer your-token"
```

## Health Check

```bash
curl "https://api.poneglyph.dev/v1/healthcheck"
```

**Response:**
```json
{
  "status": "available",
  "system_info": {
    "environment": "production",
    "version": "1.0.0"
  }
}
```

**Response:**
```json
{
  "goroutines": 7,
  "memstats": {
    "Alloc": 717192,
    "TotalAlloc": 717192,
    "Sys": 7168016,
    "HeapAlloc": 717192
  },
  "timestamp": 1757298141,
  "total_requests_received": 1250,
  "total_responses_sent": 1250,
  "total_responses_sent_by_status": {
    "200": 1100,
    "404": 50,
    "422": 100
  },
  "version": "1.0.0"
}
```

## Error Responses

### Validation Error (422)
```json
{
  "error": {
    "name": "must be provided",
    "age": "must be a positive integer",
    "bounty": "must be between 100 and 10000000000"
  }
}
```

### Not Found (404)
```json
{
  "error": "character not found"
}
```

### Rate Limited (429)
```json
{
  "error": "rate limit exceeded"
}
```

## Sorting and Filtering

list endpoints support sorting and filtering:

```bash
# Sort characters by bounty (descending)
curl "https://api.poneglyph.dev/v1/characters?sort=-bounty"

# Filter crews by minimum total bounty
curl "https://api.poneglyph.dev/v1/crews?total_bounty=1B berries"

# Complex filtering
curl "https://api.poneglyph.dev/v1/characters?race=human&age=18&origin=East Blue&sort=name"
```

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