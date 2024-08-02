package db

import (
	"database/sql"
	"log"
  "os"
  _ "github.com/mattn/go-sqlite3"

	_ "github.com/joho/godotenv/autoload"
)

var conn *sql.DB

func Connect() {
  var dbpath string
  var err error
  if envstring, ok := os.LookupEnv("DB_PATH"); ok == true {
    dbpath = envstring
  } else {
    dbpath = "./db.sqlite"
  }
  conn, err = sql.Open("sqlite3", dbpath)
  if err != nil {
    log.Fatal("could not open DB:", err)
  }
      // Ping the database to verify connection
    err = conn.Ping()
    if err != nil {
        log.Fatalf("failed to connect to the database: %v", err)
  }

  //Apply schema when connecting initially
  schemaContent, err := fs.ReadFile("schema.sql")
  if err != nil {
    log.Fatalf("Error reading Schema file: %v", err)
  }
  schema := string(schemaContent)

  _, err = conn.Exec(schema)
  if err != nil {
      log.Fatalf("failed to apply schema: %v", err)
  }

}

func Close() {
  err := conn.Close()
  log.Fatalf("Could not disconnect DB: %v", err)
}
