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

	tx, err := db.Begin()
	if err != nil {
		log.Fatal(err)
	}

	_, err = tx.Exec("INSERT INTO myschema.object(objectname, md5, count) VALUES ('cat.png', 'asdf', 2);")
	if err != nil {
		_ = tx.Rollback()
		log.Fatal(err)
	}

	rows, err := tx.Query("SELECT * FROM myschema.object;")
	if err != nil {
		_ = tx.Rollback()
		log.Fatal(err)
	}
	log.Println(rows)

	err = tx.Commit()
	if err != nil {
		log.Fatal(err)
	}
	//rows, err := db.Query("SELECT * FROM myschema.object;")
	//if err != nil {
	//	log.Println(err, "exec failed")
	//}
	//log.Println(rows)
}
