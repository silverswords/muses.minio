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

	//var object Object2
	var n int
	err = tx1.QueryRow("SELECT COUNT(*) FROM myschema.object WHERE count = 3;").Scan(&n)
	if err != nil {
		_ = tx1.Rollback()
		log.Fatal(err)
	}
	log.Println(n)

	tx2, err := db.Begin()
	if err != nil {
		log.Fatal(err)
	}

	_, err = tx2.Exec("INSERT INTO myschema.object(objectname, md5, count) VALUES ('oiar.png', '00000', 3);")
	if err != nil {
		_ = tx2.Rollback()
		log.Fatal(err)
	}

	err = tx2.Commit()
	if err != nil {
		log.Fatal(err)
	}

	err = tx1.QueryRow("SELECT COUNT(*) FROM myschema.object WHERE count = 3;").Scan(&n)
	if err != nil {
		_ = tx1.Rollback()
		log.Fatal(err)
	}
	log.Println(n)

	err = tx1.Commit()
	if err != nil {
		log.Fatal(err)
	}
}
