package models

import (
	"carbon-footprint-tracker/config"
	"database/sql"
	"time"
)

type Accommodation struct {
	ID                        int             `json:"id"`
	Date                      time.Time       `json:"date"`
	ParticipantGuestName      sql.NullString  `json:"participant_guest_name,omitempty"`
	Category                  sql.NullString  `json:"category,omitempty"` // 'Staff', 'Student', 'VIP', 'Volunteer'
	PeopleCount               int             `json:"people_count"`
	AccommodationFacilityName string          `json:"accommodation_facility_name"`  // Renamed from Location
	AccommodationType         sql.NullString  `json:"accommodation_type,omitempty"` // 'Hotel', 'Hostel', 'Guest House', 'Campus Stay'
	RoomType                  sql.NullString  `json:"room_type,omitempty"`          // 'Single', 'Double', 'Dormitory'
	NoOfRooms                 sql.NullInt32   `json:"no_of_rooms,omitempty"`
	Nights                    int             `json:"nights"`
	ElectricityConsumptionKWH sql.NullFloat64 `json:"electricity_consumption_kwh,omitempty"`
	WaterConsumptionLPD       sql.NullFloat64 `json:"water_consumption_lpd,omitempty"` // L/person/day
	MealsProvided             sql.NullBool    `json:"meals_provided,omitempty"`
	TransportModeToVenue      sql.NullString  `json:"transport_mode_to_venue,omitempty"`
	Remarks                   sql.NullString  `json:"remarks,omitempty"`
}

func (a *Accommodation) Create() error {
	query := `INSERT INTO accommodation (
		date, participant_guest_name, category, people_count, accommodation_facility_name,
		accommodation_type, room_type, no_of_rooms, nights, electricity_consumption_kwh,
		water_consumption_lpd, meals_provided, transport_mode_to_venue, remarks
	) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14) RETURNING id`
	return config.DB.QueryRow(query,
		a.Date, a.ParticipantGuestName, a.Category, a.PeopleCount, a.AccommodationFacilityName,
		a.AccommodationType, a.RoomType, a.NoOfRooms, a.Nights, a.ElectricityConsumptionKWH,
		a.WaterConsumptionLPD, a.MealsProvided, a.TransportModeToVenue, a.Remarks,
	).Scan(&a.ID)
}

func GetAllAccommodations() ([]Accommodation, error) {
	rows, err := config.DB.Query(`SELECT
		id, date, participant_guest_name, category, people_count, accommodation_facility_name,
		accommodation_type, room_type, no_of_rooms, nights, electricity_consumption_kwh,
		water_consumption_lpd, meals_provided, transport_mode_to_venue, remarks
		FROM accommodation ORDER BY date DESC`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var accommodations []Accommodation
	for rows.Next() {
		a := Accommodation{}
		err := rows.Scan(
			&a.ID, &a.Date, &a.ParticipantGuestName, &a.Category, &a.PeopleCount, &a.AccommodationFacilityName,
			&a.AccommodationType, &a.RoomType, &a.NoOfRooms, &a.Nights, &a.ElectricityConsumptionKWH,
			&a.WaterConsumptionLPD, &a.MealsProvided, &a.TransportModeToVenue, &a.Remarks,
		)
		if err != nil {
			return nil, err
		}
		accommodations = append(accommodations, a)
	}
	return accommodations, nil
}

func GetAccommodationByID(id int) (*Accommodation, error) {
	a := &Accommodation{}
	query := `SELECT
		id, date, participant_guest_name, category, people_count, accommodation_facility_name,
		accommodation_type, room_type, no_of_rooms, nights, electricity_consumption_kwh,
		water_consumption_lpd, meals_provided, transport_mode_to_venue, remarks
		FROM accommodation WHERE id = $1`
	err := config.DB.QueryRow(query, id).Scan(
		&a.ID, &a.Date, &a.ParticipantGuestName, &a.Category, &a.PeopleCount, &a.AccommodationFacilityName,
		&a.AccommodationType, &a.RoomType, &a.NoOfRooms, &a.Nights, &a.ElectricityConsumptionKWH,
		&a.WaterConsumptionLPD, &a.MealsProvided, &a.TransportModeToVenue, &a.Remarks,
	)
	if err != nil {
		return nil, err
	}
	return a, nil
}

func (a *Accommodation) Update() error {
	query := `UPDATE accommodation SET
		date=$1, participant_guest_name=$2, category=$3, people_count=$4, accommodation_facility_name=$5,
		accommodation_type=$6, room_type=$7, no_of_rooms=$8, nights=$9, electricity_consumption_kwh=$10,
		water_consumption_lpd=$11, meals_provided=$12, transport_mode_to_venue=$13, remarks=$14
		WHERE id=$15`
	_, err := config.DB.Exec(query,
		a.Date, a.ParticipantGuestName, a.Category, a.PeopleCount, a.AccommodationFacilityName,
		a.AccommodationType, a.RoomType, a.NoOfRooms, a.Nights, a.ElectricityConsumptionKWH,
		a.WaterConsumptionLPD, a.MealsProvided, a.TransportModeToVenue, a.Remarks, a.ID,
	)
	return err
}

func DeleteAccommodation(id int) error {
	query := `DELETE FROM accommodation WHERE id=$1`
	_, err := config.DB.Exec(query, id)
	return err
}
