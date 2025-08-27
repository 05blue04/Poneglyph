# Poneglyph API

> ⚠️ **Note**: This project is currently in active development. Features and API endpoints may change frequently.

A RESTful API for the One Piece universe built with Go and PostgreSQL. Explore characters, devil fruits, organizations, and locations with advanced search, filtering, and pagination capabilities.

## API Endpoints

### Characters
- `GET /characters` - List all characters
- `GET /characters/{id}` - Get a specific character
- `POST /characters` - Create a new character
- `PUT /characters/{id}` - Update an existing character
- `DELETE /characters/{id}` - Delete a character

#### Query Parameters for GET /v1/characters
- `search` - Full-text search across character names and descriptions
- `race` - Filter by character race (human, fishman, mink, giant, etc.)
- `time_skip` - Filter by time period (pre, post)
- `age` - Minimum age filter
- `bounty` - Minimum bounty filter (e.g., "1B berries", "500M berries")
- `sort` - Sort by field (id, name, age, bounty, race). Use `-` prefix for descending (e.g., `-bounty`)
- `page` - Page number for pagination (default: 1)
- `page_size` - Results per page (default: 20)

## TODO List

### High Priority
- [ ] **Write tests for Character endpoints**
  - Unit tests for handlers
  - Integration tests for database operations
  - Test validation logic

- [ ] **Create Devil Fruits table and CRUD operations**
  - Design and implement devil_fruits table schema
  - Create DevilFruit model with validation
  - Implement Devil Fruit CRUD handlers:
  - Add relationship endpoints:
    - GET /v1/characters/{id}/devil-fruit
    - PATCH /v1/characters/{id}/devil-fruit (assign/unassign)

- [ ] **Migrate database to support Organizations and Locations tables**
  - Create organizations table with seed data (Marines, Straw Hats, etc.)
  - Create character_organizations junction table for many-to-many relationships
  - Create locations table with hierarchical structure (islands → seas → regions)
  - Migrate existing character.organizations string data to proper relationships
  - Migrate character.origin strings to location references
  - Update Character model to use foreign keys instead of string fields
  - Add CRUD endpoints for organizations and locations

### Medium Priority
- [ ] **Add devil fruit user history tracking**
  - Implement devil_fruit_users_history table
  - Add endpoints for tracking fruit ownership changes
  - Handle fruit inheritance (Ace → Sabo scenarios)

- [ ] **Enhance Character-Devil Fruit integration**
  - Update character endpoints to include devil fruit data
  - Add validation to prevent multiple users per fruit
  - Implement fruit assignment/transfer operations

### Future Enhancements
- [ ] Add authentication and authorization
- [ ] Implement crew/organization relationships as separate entities
- [ ] Add rate limiting
- [ ] Add Docker containerization for easy deployment and development setup

## License

This project is licensed under the MIT License - see the LICENSE file for details.

## Acknowledgments

- Eiichiro Oda for creating the amazing One Piece universe
- Nico Robin the goat
- Monkey D. Luffy