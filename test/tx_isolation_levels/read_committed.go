package main

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq"
)

type Object1 struct {
	objectName string
	md5 string
	count int
}

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

	_, err = tx1.Exec("SET TRANSACTION ISOLATION LEVEL READ UNCOMMITTED;")
	if err != nil {
		log.Println(err)
	}

	_, err = tx1.Exec("INSERT INTO myschema.object(objectname, md5, count) VALUES ('cat.png', 'err', 2);")
	if err != nil {
		_ = tx1.Rollback()
		log.Fatal(err)
	}

	tx2, err := db.Begin()
	if err != nil {
		log.Fatal(err)
	}

	//_, err = tx2.Exec("SET TRANSACTION ISOLATION LEVEL READ UNCOMMITTED;")
	//if err != nil {
	//	log.Println(err)
	//}

	var object Object1
	rows, err := tx2.Query("SELECT * FROM myschema.object WHERE md5 = 'err';")
	if err != nil {
		_ = tx2.Rollback()
		log.Fatal(err)
	}
	defer rows.Close()

	for rows.Next() {
		err = rows.Scan(&object.objectName, &object.md5, &object.count)
		if err != nil {
			log.Println(err)
		}
	}
	log.Println(&object)

	err = tx2.Commit()
	if err != nil {
		log.Fatal(err)
	}

	err = tx1.Rollback()
	if err != nil {
		log.Fatal(err)
	}
}
