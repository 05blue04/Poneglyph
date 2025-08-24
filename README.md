# Poneglyph API

> ⚠️ **Note**: This project is currently in active development. Features and API endpoints may change frequently. Not recommended for production use at this time.

## API Endpoints

### Characters
- `GET /characters` - List all characters
- `GET /characters/{id}` - Get a specific character
- `POST /characters` - Create a new character
- `PUT /characters/{id}` - Update an existing character
- `DELETE /characters/{id}` - Delete a character

## TODO List

### High Priority
- [ ] **Finish up validator for creating Characters**
  - Complete validation rules for all character fields
  - Add custom validation for One Piece specific data

- [ ] **Finish CRUD endpoints for Characters**
  - Implement remaining HTTP handlers
  - Add proper error handling and responses

- [ ] **Write tests for Character endpoints**
  - Unit tests for handlers
  - Integration tests for database operations
  - Test validation logic

### Medium Priority
- [ ] **Implement Pre & Post timeskip flag on characters table**
  - Add database migration for timeskip fields
  - Update character model and validation
  - Modify endpoints to support timeskip data

### Future Enhancements
- [ ] Add authentication and authorization
- [ ] Implement crew/organization relationships
- [ ] Add search and filtering capabilities
- [ ] Create API documentation with Swagger
- [ ] Add rate limiting
- [ ] Implement caching layer

## License

This project is licensed under the MIT License - see the LICENSE file for details.

## Acknowledgments

- Eiichiro Oda for creating the amazing One Piece universe
- Nico Robin the goat
- Monkey D. Luffy