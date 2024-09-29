package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

var db *pgxpool.Pool

func connect_db(c string) {
	fmt.Println("Connecting to database with connection string: ", c)
	conn, err := pgxpool.New(context.Background(), c)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		os.Exit(1)
	}

	db = conn
}

func GetMigrations() {
	query := "SELECT * FROM PGMIGRATIONS"
	rows, err := db.Query(context.Background(), query)
	if err != nil {
		fmt.Fprintf(os.Stderr, "QueryRow failed: %v\n", err)
		os.Exit(1)
	}

	rowNum := 1
	for rows.Next() {
		var v []interface{}
		v, err = rows.Values()
		if err != nil {
			log.Fatal(err)
		}

		fmt.Println("Row", rowNum)

		for i := range v {
			fmt.Println(rows.FieldDescriptions()[i].Name, v[i])
		}
		rowNum++
	}
}

func checkError(err error) {
	if err != nil {
		panic(err)
	}
}

type DbMigrationRow struct {
	ID   string
	Name string
}

type migration func() string

func HasMigrationRun(name string) bool {
	fmt.Printf("checking if migratoin %s has run: ", name)
	query := fmt.Sprintf("SELECT * FROM PGMIGRATIONS WHERE NAME = '%s'", name)
	rows, err := db.Query(context.Background(), query)
	if err != nil {
		fmt.Fprintf(os.Stderr, "QueryRow failed: %v\n", err)
		os.Exit(1)
	}

	hasRun := rows.Next()

	rows.Close()

	fmt.Printf("%t \n", hasRun)

	return hasRun
}

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

	hasRun := false
	hasRun = HasMigrationRun("create_user_table")
	if !hasRun {
		runMigration("create_user_table", createUserTable)
	}
}

func createUserTable() string {
	fmt.Println("Running CreateUserTable migration")
	return "CREATE TABLE IF NOT EXISTS users (id UUID PRIMARY KEY, username TEXT, password TEXT)"
}

func createAccesskeyTable() string {
	fmt.Println("Running createAccesskeyTable migration")
	return "CREATE TABLE IF NOT EXISTS access_key (id UUID PRIMARY KEY, accesskey UUID, created_at TIMESTAMP, keep_alive TEXT, user_id UUID references users (id)"
}
