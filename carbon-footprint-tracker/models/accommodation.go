package models

import (
	"carbon-footprint-tracker/config"
	"time"
)

type Accommodation struct {
	ID          int       `json:"id"`
	PeopleCount int       `json:"people_count"`
	Nights      int       `json:"nights"`
	Date        time.Time `json:"date"`
	Location    string    `json:"location"`
}

func (a *Accommodation) Create() error {
	query := `INSERT INTO accommodation (people_count, nights, date, location) VALUES ($1, $2, $3, $4) RETURNING id`
	return config.DB.QueryRow(query, a.PeopleCount, a.Nights, a.Date, a.Location).Scan(&a.ID)
}

func GetAllAccommodations() ([]Accommodation, error) {
	rows, err := config.DB.Query(`SELECT id, people_count, nights, date, location FROM accommodation ORDER BY date DESC`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var accommodations []Accommodation
	for rows.Next() {
		a := Accommodation{}
		err := rows.Scan(&a.ID, &a.PeopleCount, &a.Nights, &a.Date, &a.Location)
		if err != nil {
			return nil, err
		}
		accommodations = append(accommodations, a)
	}
	return accommodations, nil
}

func GetAccommodationByID(id int) (*Accommodation, error) {
	a := &Accommodation{}
	query := `SELECT id, people_count, nights, date, location FROM accommodation WHERE id = $1`
	err := config.DB.QueryRow(query, id).Scan(&a.ID, &a.PeopleCount, &a.Nights, &a.Date, &a.Location)
	if err != nil {
		return nil, err
	}
	return a, nil
}

func (a *Accommodation) Update() error {
	query := `UPDATE accommodation SET people_count=$1, nights=$2, date=$3, location=$4 WHERE id=$5`
	_, err := config.DB.Exec(query, a.PeopleCount, a.Nights, a.Date, a.Location, a.ID)
	return err
}

func DeleteAccommodation(id int) error {
	query := `DELETE FROM accommodation WHERE id=$1`
	_, err := config.DB.Exec(query, id)
	return err
}
