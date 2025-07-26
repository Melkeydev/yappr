package db

import (
	"database/sql"
	"fmt"
	"log"
	
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/Melkeydev/yappr/util"
)

func NewDatabase() (*sql.DB, error) {
	env := util.GetEnv("ENVIRONMENT", "dev")
	
	var db *sql.DB
	var err error
	
	if env != "prod" {
		// Local/development environment
		dbHost := util.GetEnv("DB_HOST", "localhost")
		dbPort := util.GetEnv("DB_PORT", "5433")
		dbUser := util.GetEnv("DB_USER", "postgres")
		dbPassword := util.GetEnv("DB_PASSWORD", "postgres")
		dbName := util.GetEnv("DB_NAME", "go_chat_db")

		localDSN := fmt.Sprintf(
			"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
			dbHost, dbPort, dbUser, dbPassword, dbName,
		)
		
		log.Printf("=== DATABASE CONNECTION (DEVELOPMENT) ===")
		log.Printf("Environment: %s", env)
		log.Printf("Connecting to: %s:%s/%s", dbHost, dbPort, dbName)
		log.Printf("User: %s", dbUser)
		log.Printf("Full DSN: %s", localDSN)
		log.Printf("========================================")
		
		db, err = sql.Open("pgx", localDSN)
		if err != nil {
			log.Fatalf("Failed to open local database: %v", err)
		}
	} else {
		// Production environment - pgx handles PostgreSQL URLs natively
		connStr := util.GetEnv("CONNECTION_STRING", "")
		if connStr == "" {
			log.Fatal("CONNECTION_STRING must be set in production environment")
		}
		
		log.Printf("=== DATABASE CONNECTION (PRODUCTION) ===")
		log.Printf("Environment: %s", env)
		log.Printf("Full Connection String: %s", connStr)
		log.Printf("========================================")
		
		db, err = sql.Open("pgx", connStr)
		if err != nil {
			log.Fatalf("Failed to open production database: %v", err)
		}
	}

	return db, nil
}