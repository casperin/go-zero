package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/golang-migrate/migrate/v4/source/gitlab"
	_ "github.com/lib/pq"
)

func main() {
	connString := flag.String(
		"connection_string",
		"postgres://postgres:postgres@localhost/zero?sslmode=disable",
		"Connection string for connecting to Postgres",
	)
	flag.Parse()

	if len(os.Args) < 2 {
		fmt.Println("Call with either 'up' or 'down'")
		os.Exit(1)
	}

	migrations, e := migrate.New("file://./migrations/", *connString)
	if e != nil {
		log.Fatalln(e)
	}

	switch os.Args[1] {
	case "up":
		e = migrations.Up()
		if e != nil && e != migrate.ErrNoChange {
			log.Fatalln(e)
		}
	case "down":
		e = migrations.Down()
		if e != nil && e != migrate.ErrNoChange {
			log.Fatalln(e)
		}
		e = migrations.Drop()
		if e != nil && e != migrate.ErrNoChange {
			log.Fatalln(e)
		}
	}

	e1, e2 := migrations.Close()
	if e1 != nil {
		log.Fatalln(e1)
	}
	if e2 != nil {
		log.Fatalln(e2)
	}
}
