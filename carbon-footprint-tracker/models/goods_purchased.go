package models

import (
	"carbon-footprint-tracker/config"
	"database/sql"
	"time"
)

type GoodsPurchased struct {
	ID                  int             `json:"id"`
	Date                time.Time       `json:"date"`
	Location            string          `json:"location"` // Location where goods are delivered/used
	ItemName            string          `json:"item_name"`
	Category            sql.NullString  `json:"category,omitempty"` // 'Stationery', 'Hardware', 'Furniture', etc.
	Quantity            int             `json:"quantity"`
	Unit                sql.NullString  `json:"unit,omitempty"` // e.g., 'pcs', 'kg', 'liters'
	VendorName          sql.NullString  `json:"vendor_name,omitempty"`
	Origin              sql.NullString  `json:"origin,omitempty"`         // 'Local', 'Imported'
	TransportMode       sql.NullString  `json:"transport_mode,omitempty"` // 'Road', 'Rail', 'Air'
	TransportDistanceKM sql.NullFloat64 `json:"transport_distance_km,omitempty"`
	BillAmountINR       float64         `json:"bill_amount_inr"` // Renamed from Cost
	BillAttachmentURL   sql.NullString  `json:"bill_attachment_url,omitempty"`
	PackagingType       sql.NullString  `json:"packaging_type,omitempty"` // 'Plastic', 'Paper', 'None'
	IsRecyclable        sql.NullBool    `json:"is_recyclable,omitempty"`
	Remarks             sql.NullString  `json:"remarks,omitempty"`
}

func (g *GoodsPurchased) Create() error {
	query := `INSERT INTO goods_purchased (
		date, location, item_name, category, quantity, unit, vendor_name, origin,
		transport_mode, transport_distance_km, bill_amount_inr, bill_attachment_url,
		packaging_type, is_recyclable, remarks
	) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15) RETURNING id`
	return config.DB.QueryRow(query,
		g.Date, g.Location, g.ItemName, g.Category, g.Quantity, g.Unit, g.VendorName, g.Origin,
		g.TransportMode, g.TransportDistanceKM, g.BillAmountINR, g.BillAttachmentURL,
		g.PackagingType, g.IsRecyclable, g.Remarks,
	).Scan(&g.ID)
}

func GetAllGoodsPurchased() ([]GoodsPurchased, error) {
	rows, err := config.DB.Query(`SELECT
		id, date, location, item_name, category, quantity, unit, vendor_name, origin,
		transport_mode, transport_distance_km, bill_amount_inr, bill_attachment_url,
		packaging_type, is_recyclable, remarks
		FROM goods_purchased ORDER BY date DESC`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var goods []GoodsPurchased
	for rows.Next() {
		g := GoodsPurchased{}
		err := rows.Scan(
			&g.ID, &g.Date, &g.Location, &g.ItemName, &g.Category, &g.Quantity, &g.Unit, &g.VendorName,
			&g.Origin, &g.TransportMode, &g.TransportDistanceKM, &g.BillAmountINR, &g.BillAttachmentURL,
			&g.PackagingType, &g.IsRecyclable, &g.Remarks,
		)
		if err != nil {
			return nil, err
		}
		goods = append(goods, g)
	}
	return goods, nil
}

func GetGoodsPurchasedByID(id int) (*GoodsPurchased, error) {
	g := &GoodsPurchased{}
	query := `SELECT
		id, date, location, item_name, category, quantity, unit, vendor_name, origin,
		transport_mode, transport_distance_km, bill_amount_inr, bill_attachment_url,
		packaging_type, is_recyclable, remarks
		FROM goods_purchased WHERE id = $1`
	err := config.DB.QueryRow(query, id).Scan(
		&g.ID, &g.Date, &g.Location, &g.ItemName, &g.Category, &g.Quantity, &g.Unit, &g.VendorName,
		&g.Origin, &g.TransportMode, &g.TransportDistanceKM, &g.BillAmountINR, &g.BillAttachmentURL,
		&g.PackagingType, &g.IsRecyclable, &g.Remarks,
	)
	if err != nil {
		return nil, err
	}
	return g, nil
}

func (g *GoodsPurchased) Update() error {
	query := `UPDATE goods_purchased SET
		date=$1, location=$2, item_name=$3, category=$4, quantity=$5, unit=$6, vendor_name=$7, origin=$8,
		transport_mode=$9, transport_distance_km=$10, bill_amount_inr=$11, bill_attachment_url=$12,
		packaging_type=$13, is_recyclable=$14, remarks=$15
		WHERE id=$16`
	_, err := config.DB.Exec(query,
		g.Date, g.Location, g.ItemName, g.Category, g.Quantity, g.Unit, g.VendorName, g.Origin,
		g.TransportMode, g.TransportDistanceKM, g.BillAmountINR, g.BillAttachmentURL,
		g.PackagingType, g.IsRecyclable, g.Remarks, g.ID,
	)
	return err
}

func DeleteGoodsPurchased(id int) error {
	query := `DELETE FROM goods_purchased WHERE id=$1`
	_, err := config.DB.Exec(query, id)
	return err
}
