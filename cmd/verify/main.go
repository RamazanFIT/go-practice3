package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

func main() {
	// Load .env file
	if err := godotenv.Load(); err != nil {
		log.Println("Warning: .env file not found, using environment variables")
	}

	// Get database URL from environment
	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		log.Fatal("DATABASE_URL environment variable is required")
	}

	// Open database connection
	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatalf("Failed to open database: %v", err)
	}
	defer db.Close()

	// Verify connection
	if err := db.Ping(); err != nil {
		log.Fatalf("Failed to ping database: %v", err)
	}

	fmt.Println("✓ Database connection successful")

	// Check for required tables
	tables := []string{"users", "categories", "expenses"}
	for _, table := range tables {
		var name string
		query := `SELECT table_name FROM information_schema.tables WHERE table_schema='public' AND table_name=$1`
		err := db.QueryRow(query, table).Scan(&name)
		if err == sql.ErrNoRows {
			log.Fatalf("✗ Table '%s' does not exist", table)
		} else if err != nil {
			log.Fatalf("✗ Error checking table '%s': %v", table, err)
		}
		fmt.Printf("✓ Table '%s' exists\n", table)
	}

	// Verify table structures
	fmt.Println("\nVerifying table structures...")
	verifyUsersTable(db)
	verifyCategoriesTable(db)
	verifyExpensesTable(db)

	// Check migration version
	var version int
	var dirty bool
	err = db.QueryRow("SELECT version, dirty FROM schema_migrations").Scan(&version, &dirty)
	if err == nil {
		fmt.Printf("\n✓ Current migration version: %d (dirty: %v)\n", version, dirty)
	}

	fmt.Println("\n✓ All verifications passed!")
}

func verifyUsersTable(db *sql.DB) {
	query := `
		SELECT column_name, data_type, is_nullable
		FROM information_schema.columns
		WHERE table_schema='public' AND table_name='users'
		ORDER BY ordinal_position
	`
	rows, err := db.Query(query)
	if err != nil {
		log.Fatalf("✗ Failed to get users table schema: %v", err)
	}
	defer rows.Close()

	expectedColumns := map[string]bool{
		"id": false, "email": false, "name": false, "created_at": false,
	}

	for rows.Next() {
		var colName, dataType, isNullable string
		if err := rows.Scan(&colName, &dataType, &isNullable); err != nil {
			log.Fatalf("✗ Failed to read column: %v", err)
		}
		if _, exists := expectedColumns[colName]; exists {
			expectedColumns[colName] = true
		}
	}

	for col, found := range expectedColumns {
		if !found {
			log.Fatalf("✗ Users table missing column: %s", col)
		}
	}

	fmt.Println("✓ Users table structure verified")
}

func verifyCategoriesTable(db *sql.DB) {
	query := `
		SELECT column_name, data_type, is_nullable
		FROM information_schema.columns
		WHERE table_schema='public' AND table_name='categories'
		ORDER BY ordinal_position
	`
	rows, err := db.Query(query)
	if err != nil {
		log.Fatalf("✗ Failed to get categories table schema: %v", err)
	}
	defer rows.Close()

	expectedColumns := map[string]bool{
		"id": false, "name": false, "user_id": false,
	}

	for rows.Next() {
		var colName, dataType, isNullable string
		if err := rows.Scan(&colName, &dataType, &isNullable); err != nil {
			log.Fatalf("✗ Failed to read column: %v", err)
		}
		if _, exists := expectedColumns[colName]; exists {
			expectedColumns[colName] = true
		}
	}

	for col, found := range expectedColumns {
		if !found {
			log.Fatalf("✗ Categories table missing column: %s", col)
		}
	}

	fmt.Println("✓ Categories table structure verified")
}

func verifyExpensesTable(db *sql.DB) {
	query := `
		SELECT column_name, data_type, is_nullable
		FROM information_schema.columns
		WHERE table_schema='public' AND table_name='expenses'
		ORDER BY ordinal_position
	`
	rows, err := db.Query(query)
	if err != nil {
		log.Fatalf("✗ Failed to get expenses table schema: %v", err)
	}
	defer rows.Close()

	expectedColumns := map[string]bool{
		"id": false, "user_id": false, "category_id": false,
		"amount": false, "currency": false, "spent_at": false,
		"created_at": false, "note": false,
	}

	for rows.Next() {
		var colName, dataType, isNullable string
		if err := rows.Scan(&colName, &dataType, &isNullable); err != nil {
			log.Fatalf("✗ Failed to read column: %v", err)
		}
		if _, exists := expectedColumns[colName]; exists {
			expectedColumns[colName] = true
		}
	}

	for col, found := range expectedColumns {
		if !found {
			log.Fatalf("✗ Expenses table missing column: %s", col)
		}
	}

	fmt.Println("✓ Expenses table structure verified")
}
