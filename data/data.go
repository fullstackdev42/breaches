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
	ID1,
	ID2,
	FirstName,
	LastName,
	Gender,
	BirthPlace,
	CurrentPlace,
	Job,
	Date string
}

type DataHandler struct {
	filename string
	db       *sql.DB // Add db as a field of the DataHandler struct
}

func NewDataHandler(filename string, db *sql.DB) *DataHandler { // Pass db as a parameter to the NewDataHandler function
	return &DataHandler{filename: filename, db: db}
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

	stmt, err := tx.Prepare("INSERT OR IGNORE INTO people(ID1, ID2, FirstName, LastName, Gender, BirthPlace, CurrentPlace, Job, Date) values(?, ?, ?, ?, ?, ?, ?, ?, ?)")
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

func (d *DataHandler) CreatePeopleTable(db *sql.DB) error {
	createTableSQL := `CREATE TABLE IF NOT EXISTS people (
		"ID1" TEXT PRIMARY KEY,
		"ID2" TEXT,
		"FirstName" TEXT,
		"LastName" TEXT,
		"Gender" TEXT,
		"BirthPlace" TEXT,
		"CurrentPlace" TEXT,
		"Job" TEXT,
		"Date" TEXT
	);`

	_, err := db.Exec(createTableSQL)
	if err != nil {
		return fmt.Errorf("error creating table: %w", err)
	}

	return nil
}

func (d *DataHandler) FetchDataFromDB(offset, limit int) ([]Person, error) { // Remove db as a parameter from the FetchDataFromDB function
	rows, err := d.db.Query("SELECT * FROM people LIMIT ? OFFSET ?", limit, offset) // Use d.db instead of db
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var people []Person
	for rows.Next() {
		var person Person
		err = rows.Scan(&person.ID1, &person.ID2, &person.FirstName, &person.LastName, &person.Gender, &person.BirthPlace, &person.CurrentPlace, &person.Job, &person.Date)
		if err != nil {
			return nil, err
		}
		people = append(people, person)
	}

	return people, nil
}

func (d *DataHandler) OpenDB(path string) (*sql.DB, error) {
	db, err := sql.Open("sqlite3", path)
	if err != nil {
		return nil, fmt.Errorf("error opening database: %w", err)
	}
	return db, nil
}

func (d *DataHandler) GetTotalItems() (int, error) {
	var total int
	err := d.db.QueryRow("SELECT COUNT(*) FROM people").Scan(&total)
	if err != nil {
		return 0, err
	}
	return total, nil
}
