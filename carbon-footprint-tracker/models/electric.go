package models

import (
	"carbon-footprint-tracker/config"
	"database/sql"
	"time"
)

type ElectricConsumption struct {
	ID       int       `json:"id"`
	Date     time.Time `json:"date"`
	Location string    `json:"location"`
	Source   string    `json:"source"` // 'Main Board', 'Diesel Generator', 'Biofuel Generator', 'Solar Generation'

	// Fields specific to Generators
	DgCapacityKVA        sql.NullFloat64 `json:"dg_capacity_kva,omitempty"`
	RunningTimeHours     sql.NullFloat64 `json:"running_time_hours,omitempty"`
	FuelConsumedLiters   sql.NullFloat64 `json:"fuel_consumed_liters,omitempty"`
	FuelType             sql.NullString  `json:"fuel_type,omitempty"` // 'Diesel', 'Biofuel'
	EnergyGeneratedDgKWH sql.NullFloat64 `json:"energy_generated_dg_kwh,omitempty"`

	// Fields specific to Main Electric Board / Grid
	GridElectricityUsedKWH    sql.NullFloat64 `json:"grid_electricity_used_kwh,omitempty"`
	ElectricityBillKWH        sql.NullFloat64 `json:"electricity_bill_kwh,omitempty"`
	ElectricityBillCostINR    sql.NullFloat64 `json:"electricity_bill_cost_inr,omitempty"`
	ElectricalAppliancesCount sql.NullInt32   `json:"electrical_appliances_count,omitempty"`

	// Fields specific to Solar
	SolarGeneratedKWH sql.NullFloat64 `json:"solar_generated_kwh,omitempty"`

	Remarks sql.NullString `json:"remarks,omitempty"`
}

func (e *ElectricConsumption) Create() error {
	query := `INSERT INTO electric_consumption (
		date, location, source, dg_capacity_kva, running_time_hours, fuel_consumed_liters, 
		fuel_type, energy_generated_dg_kwh, grid_electricity_used_kwh, electricity_bill_kwh, 
		electricity_bill_cost_inr, electrical_appliances_count, solar_generated_kwh, remarks
	) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14) RETURNING id`

	return config.DB.QueryRow(query,
		e.Date, e.Location, e.Source, e.DgCapacityKVA, e.RunningTimeHours, e.FuelConsumedLiters,
		e.FuelType, e.EnergyGeneratedDgKWH, e.GridElectricityUsedKWH, e.ElectricityBillKWH,
		e.ElectricityBillCostINR, e.ElectricalAppliancesCount, e.SolarGeneratedKWH, e.Remarks,
	).Scan(&e.ID)
}

func GetAllElectricConsumptions() ([]ElectricConsumption, error) {
	rows, err := config.DB.Query(`SELECT
		id, date, location, source, dg_capacity_kva, running_time_hours, fuel_consumed_liters,
		fuel_type, energy_generated_dg_kwh, grid_electricity_used_kwh, electricity_bill_kwh,
		electricity_bill_cost_inr, electrical_appliances_count, solar_generated_kwh, remarks
		FROM electric_consumption ORDER BY date DESC`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var consumptions []ElectricConsumption
	for rows.Next() {
		e := ElectricConsumption{}
		err := rows.Scan(
			&e.ID, &e.Date, &e.Location, &e.Source, &e.DgCapacityKVA, &e.RunningTimeHours,
			&e.FuelConsumedLiters, &e.FuelType, &e.EnergyGeneratedDgKWH, &e.GridElectricityUsedKWH,
			&e.ElectricityBillKWH, &e.ElectricityBillCostINR, &e.ElectricalAppliancesCount, &e.SolarGeneratedKWH, &e.Remarks,
		)
		if err != nil {
			return nil, err
		}
		consumptions = append(consumptions, e)
	}
	return consumptions, nil
}

func GetElectricConsumptionByID(id int) (*ElectricConsumption, error) {
	e := &ElectricConsumption{}
	query := `SELECT
		id, date, location, source, dg_capacity_kva, running_time_hours, fuel_consumed_liters,
		fuel_type, energy_generated_dg_kwh, grid_electricity_used_kwh, electricity_bill_kwh,
		electricity_bill_cost_inr, electrical_appliances_count, solar_generated_kwh, remarks
		FROM electric_consumption WHERE id = $1`
	err := config.DB.QueryRow(query, id).Scan(
		&e.ID, &e.Date, &e.Location, &e.Source, &e.DgCapacityKVA, &e.RunningTimeHours,
		&e.FuelConsumedLiters, &e.FuelType, &e.EnergyGeneratedDgKWH, &e.GridElectricityUsedKWH,
		&e.ElectricityBillKWH, &e.ElectricityBillCostINR, &e.ElectricalAppliancesCount, &e.SolarGeneratedKWH, &e.Remarks,
	)
	if err != nil {
		return nil, err
	}
	return e, nil
}

func (e *ElectricConsumption) Update() error {
	query := `UPDATE electric_consumption SET
		date=$1, location=$2, source=$3, dg_capacity_kva=$4, running_time_hours=$5, fuel_consumed_liters=$6,
		fuel_type=$7, energy_generated_dg_kwh=$8, grid_electricity_used_kwh=$9, electricity_bill_kwh=$10,
		electricity_bill_cost_inr=$11, electrical_appliances_count=$12, solar_generated_kwh=$13, remarks=$14
		WHERE id=$15`
	_, err := config.DB.Exec(query,
		e.Date, e.Location, e.Source, e.DgCapacityKVA, e.RunningTimeHours, e.FuelConsumedLiters,
		e.FuelType, e.EnergyGeneratedDgKWH, e.GridElectricityUsedKWH, e.ElectricityBillKWH,
		e.ElectricityBillCostINR, e.ElectricalAppliancesCount, e.SolarGeneratedKWH, e.Remarks, e.ID,
	)
	return err
}

func DeleteElectricConsumption(id int) error {
	query := `DELETE FROM electric_consumption WHERE id=$1`
	_, err := config.DB.Exec(query, id)
	return err
}
