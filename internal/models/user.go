package models

import "time"

type User struct {
	ID        string    `json:"id" db:"id"`
	Sub       string    `json:"sub"`   
	Email     string    `json:"email" db:"email"`
	Password  string    `json:"-" db:"password"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}