package bucketStorage

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
)

type checksum struct {
	db *sql.DB
}

func newdb(configName, configPath string) (*checksum, error) {
	ac, err := GetConfig(configName, configPath)
	if err != nil {
		return nil, err
	}

	host := ac.Database["host"]
	port := ac.Database["port"]
	user := ac.Database["user"]
	password := ac.Database["password"]
	dbname := ac.Database["dbname"]

	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		host.(string), port.(int), user.(string), password.(string), dbname.(string))
	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		return nil, err
	}

	defer db.Close()
	err = db.Ping()
	if err != nil {
		return nil, err
	}

	return &checksum{
		db: db,
	}, nil
}
