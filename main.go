package main

import (
  "database/sql"
  "fmt"
  "log"
  "net/http"
  "time"

  "github.com/gorilla/mux"
  _ "github.com/lib/pq"
)

const (
  DatabaseUser     = "postgres"
  DatabasePassword = "mypassword"
  DatabaseHost     = "pgdb"
  DatabaseName     = "postgres"
)

var (
  db *sql.DB
)

/*********
* MODEL *
*********/
// Game is a simple model to store video game data and encoded/decoded through JSON.
// Has an id, title, console, rating, completed status, time created, and time updated.
type Game struct {
  ID       int       `json:"id"`
  Title    string    `json:"title"`
  Console  string    `json:"console"`
  Rating   float64   `json:"rating"`
  Complete bool      `json:"complete"`
  Created  time.Time `json:"created"`
  Updated  time.Time `json:"updated"`
}

// JsonErr is an error wrapper that can easily be Marshaled
// into JSON and retruned through a handler
type JsonErr struct {
  Error string `json:"error"`
}

/************
* DATABASE *
************/
// ConnectDB builds a connection string and connects to a
// PostgreSQL database using the standard libray's sql.DB type.
// Returns an error if the connection fails.
func ConnectDB() error {
  // build database connection string
  dbinfo := fmt.Sprintf(
    "user=%s password=%s host=%s dbname=%s sslmode=disable",
    DatabaseUser,
    DatabasePassword,
    DatabaseHost,
    DatabaseName,
  )

  // connect to database
  var err error
  db, err = sql.Open("postgres", dbinfo)
  if err != nil {
    return err
  }

  return db.Ping()
}

func main() {
  err := ConnectDB()
  if err != nil {
    log.Fatal(err)
  }

  router := mux.NewRouter()
  router.HandleFunc("/games", CreateGameHandler).Methods(http.MethodPost)
  router.HandleFunc("/games", RetrieveGamesHandler).Methods(http.MethodGet)
  router.HandleFunc("/games/{id:[0-9]+}", RetrieveGameHandler).Methods(http.MethodGet)
  router.HandleFunc("/games/{id:[0-9]+}", UpdateGameHandler).Methods(http.MethodPut, http.MethodPatch)
  router.HandleFunc("/games/{id:[0-9]+}", DeleteGameHandler).Methods(http.MethodDelete)

  log.Println("starting server on port 8000")
  log.Fatal(http.ListenAndServe(":8000", router))
}
