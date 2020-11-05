package main

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq"
)

func main() {
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		"localhost", 5432, "postgres", "123456", "minio")
	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		log.Println(err, "open failed")
		return
	}
	defer db.Close()
	err = db.Ping()
	if err != nil {
		log.Println(err, "ping failed")
	}
	fmt.Println("Successfully connected!")

	tx1, err := db.Begin()
	if err != nil {
		log.Fatal(err)
	}

	_, err = tx1.Exec("SET TRANSACTION ISOLATION LEVEL REPEATABLE READ;")
	if err != nil {
		log.Println(err)
	}

	var sum int
	rows := tx1.QueryRow("SELECT SUM(value) FROM myschema.users WHERE class = 1;")
	_ = rows.Scan(&sum)
	log.Println(sum)

	_, err = tx1.Exec("INSERT INTO myschema.users(class, value) VALUES (2, $1);", sum)
	if err != nil {
		_ = tx1.Rollback()
		log.Fatal(err)
	}

	tx2, err := db.Begin()
	if err != nil {
		log.Fatal(err)
	}

	_, err = tx2.Exec("SET TRANSACTION ISOLATION LEVEL REPEATABLE READ;")
	if err != nil {
		log.Println(err)
	}

	rows = tx2.QueryRow("SELECT SUM(value) FROM myschema.users WHERE class = 2;")
	_ = rows.Scan(&sum)
	log.Println(sum)

	_, err = tx2.Exec("INSERT INTO myschema.users(class, value) VALUES (1, $1);", sum)
	if err != nil {
		_ = tx1.Rollback()
		log.Fatal(err)
	}

	err = tx2.Commit()
	if err != nil {
		log.Fatal(err)
	}

	err = tx1.Commit()
	if err != nil {
		log.Fatal(err)
	}
}
