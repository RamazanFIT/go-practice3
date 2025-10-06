# Expense Tracker - Database Migrations

This project implements a minimal Expense Tracker database schema using versioned migrations with golang-migrate and PostgreSQL.

## Project Structure

```
practice_3/
├── internal/
│   └── db/
│       └── migrations/
│           ├── 000001_create_users_table.up.sql
│           ├── 000001_create_users_table.down.sql
│           ├── 000002_create_categories_table.up.sql
│           ├── 000002_create_categories_table.down.sql
│           ├── 000003_create_expenses_table.up.sql
│           └── 000003_create_expenses_table.down.sql
├── cmd/
│   ├── migrate/
│   │   └── main.go
│   └── verify/
│       └── main.go
├── docker-compose.yml
├── .env.example
├── go.mod
└── README.md
```

## Database Schema

### Users Table
- `id` - SERIAL primary key
- `email` - VARCHAR(255), required, unique
- `name` - VARCHAR(255), required
- `created_at` - TIMESTAMP, default NOW()

### Categories Table
- `id` - SERIAL primary key
- `name` - VARCHAR(255), required, unique per user
- `user_id` - INTEGER, nullable (null = global category)
- Foreign key: `user_id` → `users(id)` ON DELETE CASCADE
- Unique constraint: `(user_id, name)`
- Index on `user_id`

### Expenses Table
- `id` - SERIAL primary key
- `user_id` - INTEGER, required → references `users(id)` ON DELETE CASCADE
- `category_id` - INTEGER, required → references `categories(id)` ON DELETE RESTRICT
- `amount` - DECIMAL(15,2), required, positive
- `currency` - CHAR(3), required (e.g., "USD", "KZT")
- `spent_at` - TIMESTAMP, required
- `created_at` - TIMESTAMP, default NOW()
- `note` - TEXT, optional
- Check constraint: `amount > 0`
- Indexes on `user_id` and `(user_id, spent_at)`

## Prerequisites

- Go 1.25.1 or higher
- Docker and Docker Compose
- golang-migrate CLI tool (optional, for manual migrations)

## Installation

### 1. Install golang-migrate (optional)

```bash
brew install golang-migrate
```

### 2. Install Go dependencies

```bash
go mod download
```

### 3. Set up environment variables

Copy the example environment file and configure your database connection:

```bash
cp .env.example .env
```

Edit `.env` and set your database URL:
```
DATABASE_URL=postgres://postgres:postgres@localhost:5433/expense_tracker?sslmode=disable
```

## Running PostgreSQL with Docker Compose

### Start PostgreSQL

```bash
docker compose up -d
```

This will:
- Start a PostgreSQL 15 container
- Create a database named `expense_tracker`
- Expose PostgreSQL on port `5433` (to avoid conflicts with local installations)
- Create a persistent volume for database data

### Check PostgreSQL status

```bash
docker compose ps
```

### View PostgreSQL logs

```bash
docker compose logs -f postgres
```

### Stop PostgreSQL

```bash
docker compose down
```

### Stop and remove all data

```bash
docker compose down -v
```

## Running Migrations

### Using the custom migration tool (recommended)

```bash
# Apply all migrations (up)
go run cmd/migrate/main.go up

# Rollback all migrations (down)
go run cmd/migrate/main.go down

# Check current version
go run cmd/migrate/main.go version
```

### Using golang-migrate CLI

If you prefer using the CLI tool directly:

```bash
# Apply migrations
migrate -path internal/db/migrations -database "postgres://postgres:postgres@localhost:5433/expense_tracker?sslmode=disable" up

# Rollback one step
migrate -path internal/db/migrations -database "postgres://postgres:postgres@localhost:5433/expense_tracker?sslmode=disable" down 1

# Check version
migrate -path internal/db/migrations -database "postgres://postgres:postgres@localhost:5433/expense_tracker?sslmode=disable" version
```

## Verifying Migrations

After applying migrations, verify the database structure:

```bash
go run cmd/verify/main.go
```

Expected output:
```
✓ Database connection successful
✓ Table 'users' exists
✓ Table 'categories' exists
✓ Table 'expenses' exists

Verifying table structures...
✓ Users table structure verified
✓ Categories table structure verified
✓ Expenses table structure verified

✓ Current migration version: 3 (dirty: false)

✓ All verifications passed!
```

## Creating New Migrations

To create a new migration:

```bash
migrate create -ext sql -dir internal/db/migrations -seq migration_name
```

This will create two files:
- `XXXXXX_migration_name.up.sql` - for applying the migration
- `XXXXXX_migration_name.down.sql` - for rolling back the migration

## Accessing PostgreSQL Database

### Using psql

```bash
# Connect to the database
docker compose exec postgres psql -U postgres -d expense_tracker

# Or from your host machine (if you have psql installed)
psql -h localhost -p 5433 -U postgres -d expense_tracker
```

### Common SQL queries

```sql
-- List all tables
\dt

-- Describe a table
\d users

-- View migration history
SELECT * FROM schema_migrations;

-- Count records in each table
SELECT 'users' as table_name, COUNT(*) FROM users
UNION ALL
SELECT 'categories', COUNT(*) FROM categories
UNION ALL
SELECT 'expenses', COUNT(*) FROM expenses;
```

## Migration Files

1. **000001_create_users_table** - Creates the users table with email, name, and timestamps
2. **000002_create_categories_table** - Creates categories with user relationships and unique constraints
3. **000003_create_expenses_table** - Creates expenses with foreign keys to users and categories, includes amount validation

## Development Workflow

1. Start PostgreSQL:
   ```bash
   docker compose up -d
   ```

2. Run migrations:
   ```bash
   go run cmd/migrate/main.go up
   ```

3. Verify database:
   ```bash
   go run cmd/verify/main.go
   ```

4. Develop your application...

5. When done, stop PostgreSQL:
   ```bash
   docker compose down
   ```

## Troubleshooting

### Port already in use

If port 5433 is already in use, you can change it in `docker-compose.yml`:

```yaml
ports:
  - "5434:5432"  # Change 5433 to another port
```

Don't forget to update the `DATABASE_URL` in your `.env` file accordingly.

### Cannot connect to database

1. Check if PostgreSQL is running:
   ```bash
   docker compose ps
   ```

2. Check PostgreSQL logs:
   ```bash
   docker compose logs postgres
   ```

3. Verify your `.env` file has the correct `DATABASE_URL`

### Reset database completely

```bash
# Stop and remove everything
docker compose down -v

# Start fresh
docker compose up -d

# Wait a few seconds for PostgreSQL to be ready
sleep 5

# Run migrations
go run cmd/migrate/main.go up
```

## Notes

- The project uses PostgreSQL 15 Alpine for a lightweight container
- Migration files use PostgreSQL-specific syntax (SERIAL, VARCHAR, etc.)
- All foreign keys use appropriate cascade rules
- The database includes check constraints for data validation
- Persistent data is stored in a Docker volume named `practice_3_postgres_data`
