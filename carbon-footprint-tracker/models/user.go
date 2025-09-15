package models

import (
	"carbon-footprint-tracker/config"
	"carbon-footprint-tracker/utils"
	"time"
)

type User struct {
	ID           int       `json:"id"`
	Name         string    `json:"name"`
	Email        string    `json:"email"`
	PasswordHash string    `json:"password_hash,omitempty"` // Omit for security on output
	Role         string    `json:"role"`                    // 'admin', 'staff', 'viewer'
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

func (u *User) Create() error {
	var userID int
	hashedPassword, err := utils.HashPassword(u.PasswordHash)
	if err != nil {
		return err
	}
	u.PasswordHash = hashedPassword

	query := `INSERT INTO users (name, email, password_hash, role) VALUES ($1, $2, $3, $4) RETURNING id`
	err = config.DB.QueryRow(query, u.Name, u.Email, u.PasswordHash, u.Role).Scan(&userID)
	if err != nil {
		return err
	}
	u.ID = userID
	return nil
}

func GetUserByEmail(email string) (*User, error) {
	user := &User{}
	query := `SELECT id, name, email, password_hash, role FROM users WHERE email = $1`
	err := config.DB.QueryRow(query, email).Scan(&user.ID, &user.Name, &user.Email, &user.PasswordHash, &user.Role)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func GetAllUsers() ([]User, error) {
	rows, err := config.DB.Query(`SELECT id, name, email, role FROM users`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	users := []User{}
	for rows.Next() {
		user := User{}
		err := rows.Scan(&user.ID, &user.Name, &user.Email, &user.Role)
		if err != nil {
			return nil, err
		}
		users = append(users, user)
	}
	return users, nil
}

func GetUserByID(id int) (*User, error) {
	user := &User{}
	query := `SELECT id, name, email, role FROM users WHERE id = $1`
	err := config.DB.QueryRow(query, id).Scan(&user.ID, &user.Name, &user.Email, &user.Role)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (u *User) Update() error {
	query := `UPDATE users SET name=$1, email=$2, role=$3 WHERE id=$4`
	_, err := config.DB.Exec(query, u.Name, u.Email, u.Role, u.ID)
	return err
}

func DeleteUser(id int) error {
	query := `DELETE FROM users WHERE id=$1`
	_, err := config.DB.Exec(query, id)
	return err
}
