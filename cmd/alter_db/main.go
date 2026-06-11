package main

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq"
)

func main() {
	db, err := sql.Open("postgres", "postgres://postgres:root@localhost:5432/maqolameta?sslmode=disable")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	query := `
		ALTER TABLE articles 
		ALTER COLUMN title TYPE TEXT, 
		ALTER COLUMN access_type TYPE TEXT, 
		ALTER COLUMN journal TYPE TEXT, 
		ALTER COLUMN publisher TYPE TEXT, 
		ALTER COLUMN publisher_date TYPE TEXT, 
		ALTER COLUMN doi TYPE TEXT, 
		ALTER COLUMN url TYPE TEXT, 
		ALTER COLUMN pdf_url TYPE TEXT, 
		ALTER COLUMN source_url TYPE TEXT;
	`
	_, err = db.Exec(query)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Table altered successfully")
}
