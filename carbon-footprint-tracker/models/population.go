package models

import (
	"carbon-footprint-tracker/config"
	"time"
)

type Population struct {
	ID              int       `json:"id"`
	RegisteredCount int       `json:"registered_count"`
	FloatingCount   int       `json:"floating_count"`
	Date            time.Time `json:"date"`
	Location        string    `json:"location"`
}

func (p *Population) Create() error {
	query := `INSERT INTO population (registered_count, floating_count, date, location) VALUES ($1, $2, $3, $4) RETURNING id`
	return config.DB.QueryRow(query, p.RegisteredCount, p.FloatingCount, p.Date, p.Location).Scan(&p.ID)
}

func GetAllPopulations() ([]Population, error) {
	rows, err := config.DB.Query(`SELECT id, registered_count, floating_count, date, location FROM population ORDER BY date DESC`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var populations []Population
	for rows.Next() {
		p := Population{}
		err := rows.Scan(&p.ID, &p.RegisteredCount, &p.FloatingCount, &p.Date, &p.Location)
		if err != nil {
			return nil, err
		}
		populations = append(populations, p)
	}
	return populations, nil
}

func GetPopulationByID(id int) (*Population, error) {
	p := &Population{}
	query := `SELECT id, registered_count, floating_count, date, location FROM population WHERE id = $1`
	err := config.DB.QueryRow(query, id).Scan(&p.ID, &p.RegisteredCount, &p.FloatingCount, &p.Date, &p.Location)
	if err != nil {
		return nil, err
	}
	return p, nil
}

func (p *Population) Update() error {
	query := `UPDATE population SET registered_count=$1, floating_count=$2, date=$3, location=$4 WHERE id=$5`
	_, err := config.DB.Exec(query, p.RegisteredCount, p.FloatingCount, p.Date, p.Location, p.ID)
	return err
}

func DeletePopulation(id int) error {
	query := `DELETE FROM population WHERE id=$1`
	_, err := config.DB.Exec(query, id)
	return err
}
