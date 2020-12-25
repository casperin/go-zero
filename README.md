# Zero

Base for writing *simple* web stuff that requires users and auth.

The idea is to clone the repo, remove the `.git` folder, and rename everything
(search for "zero") to fit your needs.


## Requirements

* Go (probably anything above 1.0 is okay)
* Docker (to run Postgres)
    * See [Makefile](./Makefile)


## Dependencies

* [github.com/dgrijalva/jwt-go](dgrijalva/jwt-go)
* [github.com/golang-migrate/migrate](golang-migrate/migrate) v4
* [github.com/jmoiron/sqlx](jmoiron/sqlx) (and [github.com/lib/pq](lib/pq))
* [golang.org/x/crypto](x/crypto)


