package db

import (
	"database/sql"
	"log"
  "os"
  _ "github.com/mattn/go-sqlite3"

	_ "github.com/joho/godotenv/autoload"
)

var db *sql.DB

func Connect() {
  var dbpath string
  var err error
  if envstring, ok := os.LookupEnv("DB_PATH"); ok == true {
    dbpath = envstring
  } else {
    dbpath = "./db.sqlite"
  }
  db, err = sql.Open("sqlite3", dbpath)
  if err != nil {
    log.Fatal("could not open DB:", err)
  }

}


