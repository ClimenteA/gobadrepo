package gobadrepo

import (
	"testing"

	"github.com/jmoiron/sqlx"
)

type Book struct {
	Id            int    `db:"id"`
	Author        string `db:"author"`
	DatePublished string `db:"date_published"`
	Title         string `db:"title"`
	Genre         string `db:"genre"`
	Preface       string `db:"preface"`
}

// MySQL

func setupMySQL() (*sqlx.DB, error) {
	db := sqlx.MustConnect("mysql", "admin:admin@tcp(localhost:3306)/testdb")

	err := CreateTable(db, Book{})
	if err != nil {
		return nil, err
	}

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

	interfaceBooks := make([]interface{}, len(books))
	for i, d := range books {
		interfaceBooks[i] = d
	}

	err = InsertMany(db, interfaceBooks)
	if err != nil {
		return nil, err
	}

	return db, nil
}

func TestFindManyMySQL(t *testing.T) {
	db, err := setupMySQL()
	if err != nil {
		t.Fatalf("Failed to setup: %v", err)
	}

	var programmingBooks []Book
	err = FindMany(db, Book{Genre: "Programming"}, &programmingBooks)
	if err != nil || len(programmingBooks) == 0 {
		t.Error("Failed to find many books", err)
	}

	var programmingBooksLS []Book
	err = FindManyLimitSkip(db, Book{Genre: "Programming"}, &programmingBooksLS, 2, 0)
	if err != nil || len(programmingBooksLS) != 2 {
		t.Errorf("Failed to find many books with limit and skip: %v, %d", err, len(programmingBooksLS))
	}

}

func TestInsertOneAndFindOneMySQL(t *testing.T) {
	var err error
	db, err := setupMySQL()
	if err != nil {
		t.Fatalf("Failed to setup: %v", err)
	}

	book := Book{
		Author:        "Alin Dev",
		DatePublished: "UTC ISO format date",
		Title:         "SQL Tutorial",
		Genre:         "IT",
		Preface:       "Learning SQLITE3",
	}

	err = InsertOne(db, book)
	if err != nil {
		t.Errorf("Failed to insert many books: %v", err)
	}

	var alinDevBook Book
	err = FindOne(db, Book{Author: "Alin Dev"}, &alinDevBook)
	if err != nil || alinDevBook.Author != "Alin Dev" {
		t.Errorf("Failed to find one book: %v", err)
	}

}

func TestUpdateManyMySQL(t *testing.T) {
	var err error
	db, err := setupMySQL()
	if err != nil {
		t.Fatalf("Failed to setup: %v", err)
	}

	err = UpdateMany(db, Book{Author: "Alin Dev"}, Book{Author: "Alin Developer"})
	if err != nil {
		t.Errorf("Failed to update many books: %v", err)
	}
}

func TestDeleteManyMySQL(t *testing.T) {
	var err error
	db, err := setupMySQL()
	if err != nil {
		t.Fatalf("Failed to setup: %v", err)
	}

	err = DeleteMany(db, Book{Author: "Alin Developer"})
	if err != nil {
		t.Errorf("Failed to delete many books: %v", err)
	}
}

func TestDeleteAllRowsMySQL(t *testing.T) {
	var err error
	db, err := setupMySQL()
	if err != nil {
		t.Fatalf("Failed to setup: %v", err)
	}

	err = DeleteAllRows(db, Book{})
	if err != nil {
		t.Errorf("Failed to delete all rows: %v", err)
	}
}

// POSTGRES

func setupPostgres() (*sqlx.DB, error) {
	db := sqlx.MustConnect("postgres", "host=localhost port=5432 user=admin password=admin dbname=testdb sslmode=disable")

	err := CreateTable(db, Book{})
	if err != nil {
		return nil, err
	}

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

	interfaceBooks := make([]interface{}, len(books))
	for i, d := range books {
		interfaceBooks[i] = d
	}

	err = InsertMany(db, interfaceBooks)
	if err != nil {
		return nil, err
	}

	return db, nil
}

func TestFindManyPostgres(t *testing.T) {
	db, err := setupPostgres()
	if err != nil {
		t.Fatalf("Failed to setup: %v", err)
	}

	var programmingBooks []Book
	err = FindMany(db, Book{Genre: "Programming"}, &programmingBooks)
	if err != nil || len(programmingBooks) == 0 {
		t.Error("Failed to find many books", err)
	}

	var programmingBooksLS []Book
	err = FindManyLimitSkip(db, Book{Genre: "Programming"}, &programmingBooksLS, 2, 0)
	if err != nil || len(programmingBooksLS) != 2 {
		t.Errorf("Failed to find many books with limit and skip: %v, %d", err, len(programmingBooksLS))
	}

}

func TestInsertOneAndFindOnePostgres(t *testing.T) {
	var err error
	db, err := setupPostgres()
	if err != nil {
		t.Fatalf("Failed to setup: %v", err)
	}

	book := Book{
		Author:        "Alin Dev",
		DatePublished: "UTC ISO format date",
		Title:         "SQL Tutorial",
		Genre:         "IT",
		Preface:       "Learning SQLITE3",
	}

	err = InsertOne(db, book)
	if err != nil {
		t.Errorf("Failed to insert many books: %v", err)
	}

	var alinDevBook Book
	err = FindOne(db, Book{Author: "Alin Dev"}, &alinDevBook)
	if err != nil || alinDevBook.Author != "Alin Dev" {
		t.Errorf("Failed to find one book: %v", err)
	}
}

func TestUpdateManyPostgres(t *testing.T) {
	var err error
	db, err := setupPostgres()
	if err != nil {
		t.Fatalf("Failed to setup: %v", err)
	}

	err = UpdateMany(db, Book{Author: "Alin Dev"}, Book{Author: "Alin Developer"})
	if err != nil {
		t.Errorf("Failed to update many books: %v", err)
	}
}

func TestDeleteManyPostgres(t *testing.T) {
	var err error
	db, err := setupPostgres()
	if err != nil {
		t.Fatalf("Failed to setup: %v", err)
	}

	err = DeleteMany(db, Book{Author: "Alin Developer"})
	if err != nil {
		t.Errorf("Failed to delete many books: %v", err)
	}
}

func TestDeleteAllRowsPostgres(t *testing.T) {
	var err error
	db, err := setupPostgres()
	if err != nil {
		t.Fatalf("Failed to setup: %v", err)
	}

	err = DeleteAllRows(db, Book{})
	if err != nil {
		t.Errorf("Failed to delete all rows: %v", err)
	}
}

// SQLITE3

func setupSqlite3() (*sqlx.DB, error) {
	db := sqlx.MustConnect("sqlite3", ":memory:")

	err := CreateTable(db, Book{})
	if err != nil {
		return nil, err
	}

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

	interfaceBooks := make([]interface{}, len(books))
	for i, d := range books {
		interfaceBooks[i] = d
	}

	err = InsertMany(db, interfaceBooks)
	if err != nil {
		return nil, err
	}

	return db, nil
}

func TestFindManySqlite3(t *testing.T) {
	db, err := setupSqlite3()
	if err != nil {
		t.Fatalf("Failed to setup: %v", err)
	}

	var programmingBooks []Book
	err = FindMany(db, Book{Genre: "Programming"}, &programmingBooks)
	if err != nil || len(programmingBooks) == 0 {
		t.Error("Failed to find many books", err)
	}

	var programmingBooksLS []Book
	err = FindManyLimitSkip(db, Book{Genre: "Programming"}, &programmingBooksLS, 2, 0)
	if err != nil || len(programmingBooksLS) != 2 {
		t.Errorf("Failed to find many books with limit and skip: %v, %d", err, len(programmingBooksLS))
	}

}

func TestInsertOneAndFindOneSqlite3(t *testing.T) {
	var err error
	db, err := setupSqlite3()
	if err != nil {
		t.Fatalf("Failed to setup: %v", err)
	}

	book := Book{
		Author:        "Alin Dev",
		DatePublished: "UTC ISO format date",
		Title:         "SQL Tutorial",
		Genre:         "IT",
		Preface:       "Learning SQLITE3",
	}

	err = InsertOne(db, book)
	if err != nil {
		t.Errorf("Failed to insert many books: %v", err)
	}

	var alinDevBook Book
	err = FindOne(db, Book{Author: "Alin Dev"}, &alinDevBook)
	if err != nil || alinDevBook.Author != "Alin Dev" {
		t.Errorf("Failed to find one book: %v", err)
	}
}

func TestUpdateManySqlite3(t *testing.T) {
	var err error
	db, err := setupSqlite3()
	if err != nil {
		t.Fatalf("Failed to setup: %v", err)
	}

	err = UpdateMany(db, Book{Author: "Alin Dev"}, Book{Author: "Alin Developer"})
	if err != nil {
		t.Errorf("Failed to update many books: %v", err)
	}
}

func TestDeleteManySqlite3(t *testing.T) {
	var err error
	db, err := setupSqlite3()
	if err != nil {
		t.Fatalf("Failed to setup: %v", err)
	}

	err = DeleteMany(db, Book{Author: "Alin Developer"})
	if err != nil {
		t.Errorf("Failed to delete many books: %v", err)
	}
}

func TestDeleteAllRowsSqlite3(t *testing.T) {
	var err error
	db, err := setupSqlite3()
	if err != nil {
		t.Fatalf("Failed to setup: %v", err)
	}

	err = DeleteAllRows(db, Book{})
	if err != nil {
		t.Errorf("Failed to delete all rows: %v", err)
	}
}
