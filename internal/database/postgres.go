package database

import (
	"fmt"
	"log"

	"github.com/amirtalbi/examen_go/internal/config"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

func NewPostgresConnection(cfg *config.Config) (*sqlx.DB, error) {
	// First try to connect to the database
	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		cfg.Database.Host, cfg.Database.Port, cfg.Database.User, cfg.Database.Password, cfg.Database.Name)

	// Try to connect to the database
	db, err := sqlx.Open("postgres", dsn)
	if err != nil {
		log.Printf("Error opening database connection: %v", err)
		return nil, err
	}

	// Test the connection
	err = db.Ping()
	if err != nil {
		log.Printf("Error pinging database: %v", err)
		
		// If we can't connect to the database, try to connect to postgres and recreate the database
		postgresDSN := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=postgres sslmode=disable",
			cfg.Database.Host, cfg.Database.Port, cfg.Database.User, cfg.Database.Password)
		
		postgresDB, pgErr := sqlx.Connect("postgres", postgresDSN)
		if pgErr != nil {
			log.Printf("Could not connect to postgres database: %v", pgErr)
			return nil, err // Return original error
		}
		
		// Drop and recreate the database
		_, dropErr := postgresDB.Exec(fmt.Sprintf("DROP DATABASE IF EXISTS %s", cfg.Database.Name))
		if dropErr != nil {
			log.Printf("Error dropping database: %v", dropErr)
			return nil, err // Return original error
		}
		
		_, createErr := postgresDB.Exec(fmt.Sprintf("CREATE DATABASE %s", cfg.Database.Name))
		if createErr != nil {
			log.Printf("Error creating database: %v", createErr)
			return nil, err // Return original error
		}
		
		postgresDB.Close()
		
		// Try to connect to the newly created database
		db, err = sqlx.Connect("postgres", dsn)
		if err != nil {
			log.Printf("Error connecting to recreated database: %v", err)
			return nil, err
		}
	}

	// Initialize the schema
	err = initSchema(db)
	if err != nil {
		log.Printf("Error initializing schema: %v", err)
		return nil, err
	}

	log.Println("Successfully connected to the database")
	return db, nil
}

func initSchema(db *sqlx.DB) error {
	schema := `
    CREATE TABLE IF NOT EXISTS users (
        id UUID PRIMARY KEY,
        name TEXT NOT NULL,
        email TEXT UNIQUE NOT NULL,
        password TEXT NOT NULL,
        reset_token TEXT,
        reset_token_expires TIMESTAMP,
        created_at TIMESTAMP NOT NULL,
        updated_at TIMESTAMP NOT NULL
    );
    `

	_, err := db.Exec(schema)
	return err
}
