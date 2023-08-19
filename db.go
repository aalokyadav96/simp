package main
 
import (
    "database/sql"
    "fmt"
    _ "github.com/lib/pq"
)
 
 var db = dbConnect()

func dbConnect() *sql.DB{
	const (
		host     = "localhost"
		port     = 5432
		user     = "postgres"
		password = "apeman"
		dbname   = "two"
	)
    psqlconn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)
    db, err := sql.Open("postgres", psqlconn)
    if err != nil {
        fmt.Println(err)
    }
    defer db.Close()
  return db
}

func main() {
    err := db.Ping()
    if err != nil {
        fmt.Println(err)
    }
    fmt.Println("Connected!")

}

func dbGet() {
	
}