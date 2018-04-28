package databases

import (
	"fmt"
	"database/sql"

	_ "github.com/lib/pq"
)

const (
	host     = "localhost"
	port     = 5432
	user     = "forum_admin"
	password = "forum_password"
	dbname   = "db_forum"
)

var db *sql.DB

func ConnectDB() {
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
    "password=%s dbname=%s sslmode=disable",
    host, port, user, password, dbname)

	var err error
	db, err = sql.Open("postgres", psqlInfo)
	if err != nil {
		panic(err)
	}

	err = db.Ping()
	if err != nil {
		panic(err)
	}

	fmt.Println("Successfully connected!")
}

func GetPostgresSession() *sql.DB{
	return db
}

func CloseDB() {
	db.Close()
}
