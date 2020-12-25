package main

import (
	"flag"
	"log"
	"net/http"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"

	"github.com/casperin/go-zero/handlers"
)

func main() {
	port := flag.String("port", "3000", "Port for the web server")
	connString := flag.String(
		"connection_string",
		"postgres://postgres:postgres@localhost/zero?sslmode=disable",
		"Connection string for connecting to Postgres",
	)
	flag.Parse()

	db, e := sqlx.Connect("postgres", *connString)
	if e != nil {
		log.Fatalln(e)
	}

	fs := http.FileServer(http.Dir("./assets"))
	http.Handle("/assets/", http.StripPrefix("/assets/", fs))

	http.HandleFunc("/", handlers.Routes(db))

	log.Println(http.ListenAndServe(":"+*port, nil))
}
