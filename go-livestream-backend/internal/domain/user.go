package domain

import (
	"time"
)

type Role string

const (
	RoleAdmin  Role = "Admin"
	RoleAgency Role = "Agency"
	RoleHost   Role = "Host"
	RoleUser   Role = "User"
)

type User struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	Email     string    `gorm:"unique;not null" json:"email"`
	Password  string    `gorm:"not null" json:"-"`
	Role      Role      `gorm:"type:varchar(20);not null" json:"role"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
