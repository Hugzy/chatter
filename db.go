package main

import (
	"fmt"
	"log"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

type DbMigrationRow struct {
	ID   string
	Name string
}

var db *sqlx.DB

func connect_db(c string) {
	_db, err := sqlx.Connect("postgres", c)
	if err != nil {
		log.Fatalln(err)
	}

	db = _db
}

func GetMigrations() {
	query := "SELECT * FROM PGMIGRATIONS"

	dbMigrations := []DbMigrationRow{}
	err := db.Select(&dbMigrations, query)
	if err != nil {
		fmt.Println(err)
	}

	for _, dmr := range dbMigrations {
		println(dmr.Name)
	}
}

func checkError(err error) {
	if err != nil {
		panic(err)
	}
}

type migration func() string

func HasMigrationRun(name string) bool {
	fmt.Printf("checking if migratoin %s has run: ", name)
	query := fmt.Sprintf("SELECT * FROM PGMIGRATIONS WHERE NAME = '%s'", name)
	migration := DbMigrationRow{}
	err := db.Get(migration, query)
	return err == nil
}

func runMigration(name string, fn migration) {
	tx := db.MustBegin()

	db.MustBegin()

	tx.MustExec(fn())
	tx.MustExec("INSERT INTO PGMIGRATIONS (id, name, created_at) VALUES ($1, $2, $3)", uuid.New().String(), name, "now()")

	err := tx.Commit()
	if err != nil {
		err = tx.Rollback()
		if err != nil {
			panic(err)
		}
	}

	fmt.Println("Migration successful")
}

func setupDBSchema() {
	fmt.Println("Creating migration schema")
	_, err := db.Exec("CREATE TABLE IF NOT EXISTS PGMIGRATIONS (id TEXT PRIMARY KEY, name TEXT, created_at TIMESTAMP)")
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
