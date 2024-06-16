package main

import (
	"context"
	"fmt"
	"os"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

var db *pgx.Conn

func connect_db(c string) {
	fmt.Println("Connecting to database with connection string: ", c)
	conn, err := pgx.Connect(context.Background(), c)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		os.Exit(1)
	}

	db = conn
}

func checkError(err error) {
	if err != nil {
		panic(err)
	}
}

type migration func() string

func runMigration(name string, fn migration) {
	tx, err := db.Begin(context.Background())
	if err != nil {
		panic(err)
	}
	// Rollback is safe to call even if the tx is already closed, so if
	// the tx commits successfully, this is a no-op
	defer tx.Rollback(context.Background())

	_, err = tx.Exec(context.Background(), fn())

	if err != nil {
		panic(err)
	}

	_, err = tx.Exec(context.Background(), "INSERT INTO PGMIGRATIONS (id, name, created_at) VALUES ($1, $2, $3)", uuid.New().String(), name, "now()")

	if err != nil {
		panic(err)
	}

	err = tx.Commit(context.Background())
	if err != nil {
		panic(err)
	}

	fmt.Println("Migration successful")
}

func setupDBSchema() {
	fmt.Println("Creating migration schema")
	_, err := db.Exec(context.Background(), "CREATE TABLE IF NOT EXISTS PGMIGRATIONS (id TEXT PRIMARY KEY, name TEXT, created_at TIMESTAMP)")
	if err != nil {
		panic(err)
	}
	fmt.Println("Migration schema created")

	runMigration("create_user_table", createUserTable)
}

func createUserTable() string {
	fmt.Println("Running CreateUserTable migration")
	return "CREATE TABLE IF NOT EXISTS users (id UUID PRIMARY KEY, username TEXT, password TEXT)"
}
