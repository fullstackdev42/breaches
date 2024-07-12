package data

import (
	"bufio"
	"database/sql"
	"fmt"
	"os"
	"strings"

	_ "github.com/mattn/go-sqlite3"
)

type Person struct {
	ID1, ID2, FirstName, LastName, Gender, BirthPlace, CurrentPlace, Job, Date string
}

type DataHandler struct {
	filename string
}

func NewDataHandler(filename string) *DataHandler {
	return &DataHandler{filename: filename}
}

func (d *DataHandler) LoadDataFromFile() ([]Person, error) {
	file, err := os.Open(d.filename)
	if err != nil {
		return nil, fmt.Errorf("error opening file: %w", err)
	}
	defer file.Close()

	var people []Person
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		parts := strings.Split(line, ":")
		if len(parts) >= 10 {
			people = append(people, Person{
				ID1:          parts[0],
				ID2:          parts[1],
				FirstName:    parts[2],
				LastName:     parts[3],
				Gender:       parts[4],
				BirthPlace:   parts[5],
				CurrentPlace: parts[6],
				Job:          parts[8],
				Date:         parts[9],
			})
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("error reading file: %w", err)
	}

	return people, nil
}

func (d *DataHandler) LoadDataIntoDB(db *sql.DB, people []Person) error {
	tx, err := db.Begin()
	if err != nil {
		return fmt.Errorf("error starting transaction: %w", err)
	}

	stmt, err := tx.Prepare("INSERT INTO people(ID1, ID2, FirstName, LastName, Gender, BirthPlace, CurrentPlace, Job, Date) values(?, ?, ?, ?, ?, ?, ?, ?, ?)")
	if err != nil {
		return fmt.Errorf("error preparing statement: %w", err)
	}
	defer stmt.Close()

	for _, person := range people {
		_, err = stmt.Exec(person.ID1, person.ID2, person.FirstName, person.LastName, person.Gender, person.BirthPlace, person.CurrentPlace, person.Job, person.Date)
		if err != nil {
			return fmt.Errorf("error executing statement: %w", err)
		}
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("error committing transaction: %w", err)
	}

	return nil
}
