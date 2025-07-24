package db

import (
	"database/sql"
	"fmt"
	"log"
	
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/melkeydev/chat-go/util"
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
		
		log.Printf("Connecting to local database: %s:%s/%s", dbHost, dbPort, dbName)
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
		
		log.Printf("Connecting to production database with pgx driver")
		log.Printf("Connection string length: %d", len(connStr))
		
		db, err = sql.Open("pgx", connStr)
		if err != nil {
			log.Fatalf("Failed to open production database: %v", err)
		}
	}

	return db, nil
}