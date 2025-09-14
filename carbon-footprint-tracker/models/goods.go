package models

import (
	"carbon-footprint-tracker/config"
	"time"
)

type GoodsPurchased struct {
	ID       int       `json:"id"`
	ItemName string    `json:"item_name"`
	Quantity int       `json:"quantity"`
	Cost     float64   `json:"cost"`
	Date     time.Time `json:"date"`
	Location string    `json:"location"`
}

func (g *GoodsPurchased) Create() error {
	query := `INSERT INTO goods_purchased (item_name, quantity, cost, date, location) VALUES ($1, $2, $3, $4, $5) RETURNING id`
	return config.DB.QueryRow(query, g.ItemName, g.Quantity, g.Cost, g.Date, g.Location).Scan(&g.ID)
}

func GetAllGoodsPurchased() ([]GoodsPurchased, error) {
	rows, err := config.DB.Query(`SELECT id, item_name, quantity, cost, date, location FROM goods_purchased ORDER BY date DESC`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var goods []GoodsPurchased
	for rows.Next() {
		g := GoodsPurchased{}
		err := rows.Scan(&g.ID, &g.ItemName, &g.Quantity, &g.Cost, &g.Date, &g.Location)
		if err != nil {
			return nil, err
		}
		goods = append(goods, g)
	}
	return goods, nil
}

func GetGoodsPurchasedByID(id int) (*GoodsPurchased, error) {
	g := &GoodsPurchased{}
	query := `SELECT id, item_name, quantity, cost, date, location FROM goods_purchased WHERE id = $1`
	err := config.DB.QueryRow(query, id).Scan(&g.ID, &g.ItemName, &g.Quantity, &g.Cost, &g.Date, &g.Location)
	if err != nil {
		return nil, err
	}
	return g, nil
}

func (g *GoodsPurchased) Update() error {
	query := `UPDATE goods_purchased SET item_name=$1, quantity=$2, cost=$3, date=$4, location=$5 WHERE id=$6`
	_, err := config.DB.Exec(query, g.ItemName, g.Quantity, g.Cost, g.Date, g.Location, g.ID)
	return err
}

func DeleteGoodsPurchased(id int) error {
	query := `DELETE FROM goods_purchased WHERE id=$1`
	_, err := config.DB.Exec(query, id)
	return err
}
