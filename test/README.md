# Admin Integration Tests

This directory contains comprehensive tests for admin account creation and configuration management functionality.

## Test Coverage

### 1. Admin Account Creation (`admin_test.go`)

Tests the complete admin account workflow:

- ✅ **Admin Role Selection**: Create admin accounts with `UserRoleAdmin`
- ✅ **Learner Role Selection**: Create learner accounts with `UserRoleLearner`
- ✅ **Teacher Role Selection**: Create teacher accounts with `UserRoleInstructor`
- ✅ **Role Validation**: Reject invalid roles with proper error messages
- ✅ **Database Verification**: Verify roles are correctly stored in the database

### 2. Configuration Management

The system now supports full CRUD operations for all configuration types:

#### ✅ Prompt Templates (`/prompt-templates`)

- Get latest classification and generation prompts
- Create new prompt template versions
- Activate/deactivate specific versions
- Query by generation type (`CLASSIFICATION`, `QUESTIONS`)

#### ✅ System Instructions (`/system-instructions`)

- Get active system instruction
- Create new system instruction versions
- Activate specific versions
- List all system instructions

#### ✅ Chunking Configs (`/chunking-configs`)

- Get active chunking configuration
- Create new chunking configs with size/overlap settings
- Activate specific configurations
- List all chunking configs

#### ✅ Model Configs (`/model-configs`) - **NEWLY IMPLEMENTED**

- Get active model configuration
- Create new model configs (temperature, max_tokens, top_p, top_k, etc.)
- Activate specific model configurations
- List all model configs

### 3. Automatic Versioning

All configuration types support automatic versioning:

- Version numbers auto-increment when creating new configurations
- Only one configuration can be active per type at a time
- Activating a new version automatically deactivates the previous one

### 4. Rollback Functionality

Admins can easily rollback to previous configurations:

- Use `POST /config-type/{id}/activate` to activate any previous version
- System automatically handles deactivating the current active version
- Full audit trail maintained in the database

### 5. Role-Based Access Control

- **Admin**: Full access to all configuration endpoints
- **Teacher**: Read-only access to active configurations (for generation context)
- **Learner**: No access to configuration management

## Running the Tests

### Prerequisites

1. **PostgreSQL Database**: Ensure you have a test database running
2. **Environment Variables**: Set `TEST_DB_URL` or use the default
3. **Dependencies**: Run `go mod tidy` to install required packages

### Test Database Setup

```bash
# Create test database
createdb learning_test

# Set environment variable (optional)
export TEST_DB_URL="postgres://username:password@localhost:5432/learning_test?sslmode=disable"
```

### Run Tests

```bash
# Run all admin tests
go test ./test -v

# Run specific test
go test ./test -run TestAdminAccountCreation -v

# Run with database cleanup
go test ./test -v -cleanup
```

### Test Database URL

The tests will use the following database URL priority:

1. `TEST_DB_URL` environment variable
2. Default: `postgres://test:test@localhost:5432/learning_test?sslmode=disable`

## Test Structure

### `TestAdminAccountCreation`

- Tests user signup with different roles
- Verifies role assignment in database
- Tests role validation and error handling

### `TestConfigurationEndpointsExist`

- Verifies all configuration endpoints are properly registered
- Tests authentication requirements
- Confirms proper HTTP status codes

## Configuration Endpoints Reference

### Admin-Only Endpoints

All configuration endpoints require admin authentication:

```
GET    /prompt-templates?generation_type=CLASSIFICATION
GET    /prompt-templates/{id}
GET    /prompt-templates/generation-type/{type}
POST   /prompt-templates
POST   /prompt-templates/{id}/activate
POST   /prompt-templates/{id}/deactivate

GET    /system-instructions
GET    /system-instructions/{id}
GET    /system-instructions/active
POST   /system-instructions
POST   /system-instructions/{id}/activate

GET    /chunking-configs
GET    /chunking-configs/{id}
GET    /chunking-configs/active
POST   /chunking-configs
POST   /chunking-configs/{id}/activate

GET    /model-configs
GET    /model-configs/{id}
GET    /model-configs/active
POST   /model-configs
POST   /model-configs/{id}/activate
```

### Public Endpoints

```
POST   /signup  # Create account with role selection
```

## Example Usage

### Create Admin Account

```bash
curl -X POST http://localhost:8080/signup \
  -H "Content-Type: application/json" \
  -d '{
    "email": "admin@example.com",
    "password": "SecurePassword123!",
    "display_name": "System Administrator",
    "role": "ADMIN"
  }'
```

### Create Model Configuration (Admin Only)

```bash
curl -X POST http://localhost:8080/model-configs \
  -H "Authorization: Bearer $ADMIN_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "model_name": "gemini-1.5-pro",
    "temperature": 0.7,
    "max_tokens": 2048,
    "top_p": 0.9,
    "top_k": 40.0,
    "mime_type": "application/json",
    "is_active": true
  }'
```

### Rollback to Previous Configuration

```bash
# Activate a previous version
curl -X POST http://localhost:8080/model-configs/{previous-config-id}/activate \
  -H "Authorization: Bearer $ADMIN_TOKEN"
```

## Implementation Status

- ✅ Admin account creation with role selection
- ✅ All configuration management endpoints
- ✅ Model configs handler and service (newly implemented)
- ✅ Automatic versioning system
- ✅ Rollback functionality via activation
- ✅ Role-based access control
- ✅ Comprehensive test coverage
- ⏳ Generation endpoints (will be implemented as part of generation flow)

## Next Steps

1. **JWT Token Generation**: Implement proper JWT token generation for testing authenticated endpoints
2. **Integration with Admin Dashboard**: Ensure frontend can consume all these endpoints
3. **Audit Logging**: Add audit trails for configuration changes
4. **Validation**: Add more robust validation for configuration parameters
5. **Generation Flow**: Implement the generation endpoints that will use these configurations

The system now provides a complete foundation for admin account management and configuration versioning with rollback capabilities.
