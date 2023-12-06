# GoBadRepo

Useful CRUD funcs over sqlx go package. GoBadRepo uses structs as input to generate safe SQL queries which then are passed to sqlx package. GoBadRepo works with sqlite3, mysql, postgres. 


## Install

Simple install with:

```shell
go get -u github.com/ClimenteA/gobadrepo
```


## Usage

Here is a basic usage of GoBadRepo which you can copy/paste in your project and start using it.

```go
package main

import ( 
	repo "github.com/ClimenteA/gobadrepo"
	"github.com/jmoiron/sqlx"
)

type Book struct {
	Id            int    `db:"id"` // mandatory for each struct
	Author        string `db:"author"`
	DatePublished string `db:"date_published"`
	Title         string `db:"title"`
	Genre         string `db:"genre"`
	Preface       string `db:"preface"`
}

func main() {

	var err error

	db := sqlx.MustConnect("sqlite3", ":memory:")
	// db := sqlx.MustConnect("mysql", "admin:admin@tcp(localhost:3306)/testdb")
	// db := sqlx.MustConnect("postgres", "host=localhost port=5432 user=admin password=admin dbname=testdb sslmode=disable")

	// CREATE TABLE
	err = repo.CreateTable(db, Book{})
	check(err)

	book := Book{
		Author:        "Alin Dev",
		DatePublished: "UTC ISO format date",
		Title:         "SQL Tutorial",
		Genre:         "IT",
		Preface:       "Learning SQLITE3",
	}

	// INSERT ONE / FIND ONE
	err = repo.InsertOne(db, book)
	check(err)

	var alinDevBook Book
	err = repo.FindOne(db, Book{Author: "Alin Dev"}, &alinDevBook)
	check(err)

	// INSERT MANY / FIND MANY
	books := []Book{
		{
			Author:        "Alin Devon",
			DatePublished: "2056-34-34",
			Title:         "Python Tutorial",
			Genre:         "Programming",
			Preface:       "Learning Python",
		},
		{
			Author:        "Cornel Marcon",
			DatePublished: "2020-05-23",
			Title:         "SQL Tutorial",
			Genre:         "SQL",
			Preface:       "Learning SQLITE3",
		},
		{
			Author:        "Razvan Rapden",
			DatePublished: "2020-23-04",
			Title:         "Java Tutorial",
			Genre:         "Programming",
			Preface:       "Learning Java",
		},
	}

	// Open a PR if you know how a way to avoid this intermediary step
	interfaceBooks := make([]interface{}, len(books))
	for i, d := range books {
		interfaceBooks[i] = d
	}

	err = repo.InsertMany(db, interfaceBooks)
	check(err)

	// RETURN ALL ROWS WHERE Genre COLUMN HAVE VALUE "Programming" 
	var programmingBooks []Book
	err = repo.FindMany(db, Book{Genre: "Programming"}, &programmingBooks)
	check(err)

	// PAGINATION
	var programmingBooksLS []Book
	err = repo.FindManyLimitSkip(db, Book{Genre: "Programming"}, &programmingBooksLS, 2, 0)
	check(err)

	// UPDATE MANY
	// Use Id match to update one
	err = repo.UpdateMany(db, Book{Author: "Alin Dev"}, Book{Author: "Alin Developer"})
	check(err)

	// DELETE MANY
	// Use Id match to delete one
	err = repo.DeleteMany(db, Book{Author: "Alin Developer"})
	check(err)

	// CLEAR TABLE
	err = repo.DeleteAllRows(db, Book{})
	check(err)

}

func check(err error) {
	if err != nil {
		panic(err)
	}
}

```

GoBadRepo is not advanced by any means and has it's hard limitations like:
- Accepts only field struct type of `string` (`TEXT`), `int*` (`INT`), `float*` (`REAL`);
- You have to use sqlx for more complex queries;


## Tests

Every CRUD function is tested for sqlite3, mysql, postgres databases. 

To run tests spin up the databases:
```shell
docker-compose up -d
```
then 
```shell
$ go test .

ok      github.com/ClimenteA/gobadrepo  0.050s
```
or 
```shell
$ go test -cover

PASS
coverage: 94.6% of statements
ok      github.com/ClimenteA/gobadrepo  0.051s
```
 
