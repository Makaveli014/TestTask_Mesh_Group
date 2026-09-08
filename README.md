# SAP Segmentation Import Module

This Go module imports segmentation data from a third-party ERP system into a PostgreSQL database.

## Features
- **Configurable**: Uses `envconfig` for environment-based configuration.
- **Idempotent**: Uses PostgreSQL `ON CONFLICT` (Upsert) to update existing records based on `address_sap_id`.
- **Resilient**: Implements HTTP timeouts, pagination, and request intervals.
- **Maintenance**: Automatically cleans up old log files based on a configurable age.
- **Tested**: Includes unit tests for API and Log logic, and integration tests for Database operations.

## Project Structure
- `cmd/sap_segmentationd/main.go`: Entry point of the application.
- `model/`: Data models and business logic (API client, DB operations, Log cleanup).
- `setup/install.sql`: SQL migration to initialize the database schema.
- `log/`: Directory where application logs are stored.

## Configuration
The application is configured via environment variables. Defaults are provided in the code.

| Variable | Description | Default |
|----------|-------------|---------|
| `DB_HOST` | DB Server IP | `127.0.0.1` |
| `DB_PORT` | DB Server Port | `5432` |
| `DB_NAME` | DB Name | `mesh_group` |
| `DB_USER` | DB User | `postgres` |
| `DB_PASSWORD` | DB Password | `postgres` |
| `CONN_URI` | API Endpoint | `http://bsm.api.iql.ru/...` |
| `CONN_AUTH_LOGIN_PWD` | API Auth (login:pwd) | `4Dfddf5:jKlljHGH` |
| `CONN_USER_AGENT` | User Agent | `spacecount-test` |
| `CONN_TIMEOUT` | API Timeout (sec) | `5` |
| `CONN_INTERVAL` | Delay between batches (ms) | `1500` |
| `IMPORT_BATCH_SIZE` | Records per batch | `50` |
| `LOG_CLEANUP_MAX_AGE` | Log retention (days) | `7` |

## Getting Started

### 1. Database Setup
Run the migration script in your PostgreSQL instance:
```bash
psql -U postgres -d mesh_group -f setup/install.sql
```

### 2. Run the Application
```bash
go run cmd/sap_segmentationd/main.go
```

### 3. Running Tests
To run all tests (unit and integration):
```bash
go test ./...
```
*Note: Integration tests require a running PostgreSQL instance with a `mesh_group_test` database.*

## Best Practices Applied
- **Dependency Injection**: The API client accepts an `http.Client` interface, making it easily mockable for unit tests.
- **Transaction Management**: Batch inserts are wrapped in transactions for performance and consistency.
- **Safe File I/O**: Log cleanup uses `os.ReadDir` and modification time checks to safely manage disk space.
- **Modular Design**: Clear separation between configuration, data access, and application flow.
