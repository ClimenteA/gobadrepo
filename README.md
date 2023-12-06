# Go Bad Repo

Useful CRUD funcs over sqlx go package. 


## Usage

```go
package main

import "github.com/jmoiron/sqlx"

type ProgrammingBook struct {
	Id            int    `db:"id"`
	Author        string `db:"author"`
	DatePublished string `db:"date_published"`
	Title         string `db:"title"`
	Genre         string `db:"genre"`
	Preface       string `db:"preface"`
}

func main() {

	db := sqlx.MustConnect("sqlite3", ":memory:")
	// db := sqlx.MustConnect("mysql", "admin:admin@tcp(localhost:3306)/testdb")
	// db := sqlx.MustConnect("postgres", "host=localhost port=5432 user=admin password=admin dbname=testdb sslmode=disable")

	var err error

	err = CreateTable(db, ProgrammingBook{})
	check(err)

    // TODO

}

func check(err error) {
	if err != nil {
		panic(err)
	}
}

```
