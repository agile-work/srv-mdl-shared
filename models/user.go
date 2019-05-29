package models

import (
	"time"

	shared "github.com/agile-work/srv-shared"
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
	ReceiveEmails string    `json:"receive_emails" sql:"receive_emails"`
	Active        bool      `json:"active" sql:"active"`
	CreatedBy     string    `json:"created_by" sql:"created_by"`
	CreatedAt     time.Time `json:"created_at" sql:"created_at"`
	UpdatedBy     string    `json:"updated_by" sql:"updated_by"`
	UpdatedAt     time.Time `json:"updated_at" sql:"updated_at"`
}

// GetUserSelectableFields returns an array with the user fields
func GetUserSelectableFields() []string {
	return []string{
		shared.TableCoreUsers + ".id as user_id",
		shared.TableCoreUsers + ".username",
		shared.TableCoreUsers + ".first_name",
		shared.TableCoreUsers + ".last_name",
		shared.TableCoreUsers + ".email",
	}

}

// ViewUserAllPermissions defines the struct of this object
type ViewUserAllPermissions struct {
	UserID         string `json:"user_id" sql:"user_id"`
	SchemaID       string `json:"schema_id" sql:"schema_id"`
	SchemaCode     string `json:"schema_code" sql:"schema_code"`
	SchemaName     string `json:"schema_name" sql:"schema_name"`
	StructureID    string `json:"structure_id" sql:"structure_id"`
	StructureCode  string `json:"structure_code" sql:"structure_code"`
	StructureType  string `json:"structure_type" sql:"structure_type"`
	StructureClass string `json:"structure_class" sql:"structure_class"`
	StructureName  string `json:"structure_name" sql:"structure_name"`
	LanguageCode   string `json:"language_code" sql:"language_code"`
	PermissionType int    `json:"permission_type" sql:"permission_type"`
	Scope          string `json:"scope" sql:"scope"`
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
	ReceiveEmails string    `json:"receive_emails" sql:"receive_emails"`
	Active        bool      `json:"active" sql:"active"`
	CreatedBy     string    `json:"created_by" sql:"created_by"`
	CreatedAt     time.Time `json:"created_at" sql:"created_at"`
	UpdatedBy     string    `json:"updated_by" sql:"updated_by"`
	UpdatedAt     time.Time `json:"updated_at" sql:"updated_at"`
}

// UserCustomClaims used to parse token payload
type UserCustomClaims struct {
	User User
	jwt.StandardClaims
}
