package models

import (
	"time"

	jwt "github.com/dgrijalva/jwt-go"
)

// User defines the struct of this object
type User struct {
	ID            string    `json:"id" sql:"id" pk:"true"`
	Username      string    `json:"username" sql:"username"`
	FirstName     string    `json:"first_name" sql:"first_name"`
	LastName      string    `json:"last_name" sql:"last_name"`
	Email         string    `json:"email" sql:"email"`
	Password      string    `json:"password" sql:"password"`
	LanguageCode  string    `json:"language_code" sql:"language_code"`
	Active        bool      `json:"active" sql:"active"`
	CreatedBy     string    `json:"created_by" sql:"created_by"`
	CreatedByUser *User     `json:"created_by_user" table:"core_users" alias:"created_by_user" on:"created_by_user.id = core_users.created_by"`
	CreatedAt     time.Time `json:"created_at" sql:"created_at"`
	UpdatedBy     string    `json:"updated_by" sql:"updated_by"`
	UpdatedByUser *User     `json:"updated_by_user" table:"core_users" alias:"updated_by_user" on:"updated_by_user.id = core_users.updated_by"`
	UpdatedAt     time.Time `json:"updated_at" sql:"updated_at"`
}

// ViewGroupUser defines the struct of this object
type ViewGroupUser struct {
	ID            string    `json:"id" sql:"id" pk:"true"`
	GroupID       string    `json:"group_id" sql:"group_id" fk:"true"`
	Username      string    `json:"username" sql:"username"`
	FirstName     string    `json:"first_name" sql:"first_name"`
	LastName      string    `json:"last_name" sql:"last_name"`
	Email         string    `json:"email" sql:"email"`
	Password      string    `json:"password" sql:"password"`
	LanguageCode  string    `json:"language_code" sql:"language_code"`
	Active        bool      `json:"active" sql:"active"`
	CreatedBy     string    `json:"created_by" sql:"created_by"`
	CreatedByUser *User     `json:"created_by_user" table:"core_users" alias:"created_by_user" on:"created_by_user.id = core_v_group_users.created_by"`
	CreatedAt     time.Time `json:"created_at" sql:"created_at"`
	UpdatedBy     string    `json:"updated_by" sql:"updated_by"`
	UpdatedByUser *User     `json:"updated_by_user" table:"core_users" alias:"updated_by_user" on:"updated_by_user.id = core_v_group_users.updated_by"`
	UpdatedAt     time.Time `json:"updated_at" sql:"updated_at"`
}

// UserCustomClaims used to parse token payload
type UserCustomClaims struct {
	User User
	jwt.StandardClaims
}
