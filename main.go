package main

import (
	"fmt"

	_ "github.com/jackc/pgx/v5/stdlib" // Standard library bindings for pgx
	"github.com/jmoiron/sqlx"
)

/*
	Produces

WITH sff (name, age, job) AS (VALUES ($1, $2::text, $3),($4, $5::text, $6)) SELECT * FROM sff;
Inserted person: {Name:Bob Age:12 Job:Cleaner}
Inserted person: {Name:Erica Age:552 Job:Ice Sculptor}
*/
func main() {

	host := "postgres://postgres:postgrespw@localhost:49153"

	// var db *sqlx.DB

	db := sqlx.MustConnect("pgx", host)

	pp := []*Person{
		{
			Name: "Bob",
			Age:  12,
			Job:  "Cleaner",
		},
		{
			Name: "Erica",
			Age:  552,
			Job:  "Ice Sculptor",
		},
	}
	err := FormatPeopleUsingSQLWith(db, pp)
	if err != nil {
		panic(err)
	}

}

type Person struct {
	Name string `db:"name"`
	Age  int    `db:"age"`
	Job  string `db:"job"`
}

func FormatPeopleUsingSQLWith(db *sqlx.DB, personSlice []*Person) error {
	queryString := `WITH sff (name, age, job) AS (VALUES `

	params := make([]interface{}, 0)
	for i, p := range personSlice {
		params = append(params, p.Name, fmt.Sprint(p.Age), p.Job)
		queryString += fmt.Sprintf("($%d, $%d::text, $%d),", i*3+1, i*3+2, i*3+3)
	}

	queryString = queryString[:len(queryString)-1] // drop the last comma

	queryString += `) SELECT * FROM sff;`
	fmt.Println(queryString)
	rows, err := db.Queryx(queryString, params...)
	if err != nil {
		return err
	}
	defer rows.Close()

	for rows.Next() {
		var person Person
		err := rows.StructScan(&person)
		if err != nil {
			return err
		}
		fmt.Printf("Inserted person: %+v\n", person)
	}

	return rows.Err()
}
