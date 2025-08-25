# Poneglyph API

> ⚠️ **Note**: This project is currently in active development. Features and API endpoints may change frequently.

## API Endpoints

### Characters
- `GET /characters` - List all characters
- `GET /characters/{id}` - Get a specific character
- `POST /characters` - Create a new character
- `PUT /characters/{id}` - Update an existing character
- `DELETE /characters/{id}` - Delete a character

## TODO List

### High Priority
- [ ] **Write tests for Character endpoints**
  - Unit tests for handlers
  - Integration tests for database operations
  - Test validation logic

- [ ] **Implement GET /v1/characters endpoint with filtering and pagination**
  - Add query parameters for filtering (race, organization, time_skip, bounty range, etc.)
  - Implement pagination (limit, offset/cursor-based)
  - Add sorting options (by name, age, bounty, episode)
  - Support search by name/description

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

## License

This project is licensed under the MIT License - see the LICENSE file for details.

## Acknowledgments

- Eiichiro Oda for creating the amazing One Piece universe
- Nico Robin the goat
- Monkey D. Luffy