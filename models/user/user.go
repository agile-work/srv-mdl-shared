package user

import (
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"github.com/agile-work/srv-mdl-shared/models/customerror"
	"github.com/agile-work/srv-mdl-shared/models/instance"

	"github.com/agile-work/srv-shared/constants"
	"github.com/agile-work/srv-shared/rdb"
	"github.com/agile-work/srv-shared/sql-builder/builder"
	"github.com/agile-work/srv-shared/sql-builder/db"
)

// Users defines the array struct of this object
type Users []User

// LoadAll defines all instances from the object
func (u *Users) LoadAll(opt *db.Options) error {
	if err := db.SelectStruct(constants.TableCoreUsers, u, opt); err != nil {
		return customerror.New(http.StatusInternalServerError, "users load", err.Error())
	}
	return nil
}

// User swagg-doc:model
// id,created_by,created_at,updated_by,updated_at,token swagg-doc:attribute:ignore_write
// security,security_instances swagg-doc:attribute:ignore
// Defines a user model
type User struct {
	ID                string             `json:"id" sql:"id" pk:"true"`
	Username          string             `json:"username" sql:"username" updatable:"false" validate:"required"`
	FirstName         string             `json:"first_name" sql:"first_name" validate:"required"`
	LastName          string             `json:"last_name" sql:"last_name" validate:"required"`
	Email             string             `json:"email" sql:"email" updatable:"false" validate:"required"`
	Password          string             `json:"password,omitempty" sql:"password" updatable:"false" validate:"required"`
	LanguageCode      string             `json:"language_code" sql:"language_code"`
	ReceiveEmails     string             `json:"receive_emails" sql:"receive_emails"`
	Security          *security          `json:"security,omitempty" sql:"security" field:"jsonb"`
	SecurityInstances *securityInstances `json:"security_instances,omitempty" sql:"security_instances" field:"jsonb"`
	Active            bool               `json:"active" sql:"active"`
	Token             string             `json:"token"`
	CreatedBy         string             `json:"created_by" sql:"created_by"`
	CreatedAt         time.Time          `json:"created_at" sql:"created_at"`
	UpdatedBy         string             `json:"updated_by" sql:"updated_by"`
	UpdatedAt         time.Time          `json:"updated_at" sql:"updated_at"`
}

// Create persists the struct creating a new object in the database
func (u *User) Create(trs *db.Transaction, columns ...string) error {
	id, err := db.InsertStructTx(trs.Tx, constants.TableCoreUsers, u, columns...)
	if err != nil {
		return customerror.New(http.StatusInternalServerError, "user create", err.Error())
	}
	u.ID = id

	resource := instance.Instance{}
	resource.ID = db.UUID()
	resource.ParentID = u.ID
	resource.CreatedAt = u.CreatedAt
	resource.CreatedBy = u.CreatedBy
	resource.UpdatedAt = u.UpdatedAt
	resource.UpdatedBy = u.UpdatedBy
	db.InsertStructTx(trs.Tx, constants.TableCustomResources, &resource)

	return nil
}

// Load defines only one object from the database
func (u *User) Load() error {
	cache, _ := rdb.Get("instance:user:" + u.Username)

	if cache != "" {
		if err := json.Unmarshal([]byte(cache), u); err != nil {
			return customerror.New(http.StatusInternalServerError, "user parse from cache", err.Error())
		}
	} else {
		if err := db.SelectStruct(constants.TableCoreUsers, u, &db.Options{
			Conditions: builder.Equal("username", u.Username),
		}); err != nil {
			return customerror.New(http.StatusInternalServerError, "user load", err.Error())
		}
		jsonBytes, err := json.Marshal(u)
		if err != nil {
			return customerror.New(http.StatusInternalServerError, "user parse to cache", err.Error())
		}
		if err := rdb.Set("instance:user:"+u.Username, string(jsonBytes), 0); err != nil {
			return customerror.New(http.StatusInternalServerError, "user parse save cache", err.Error())
		}
	}

	return nil
}

// Update updates object data in the database
func (u *User) Update(trs *db.Transaction, columns []string) error {
	opt := &db.Options{Conditions: builder.Equal("username", u.Username)}

	if len(columns) > 0 {
		if err := db.UpdateStructTx(trs.Tx, constants.TableCoreUsers, u, opt, strings.Join(columns, ",")); err != nil {
			return customerror.New(http.StatusInternalServerError, "user update", err.Error())
		}
	} else {
		return customerror.New(http.StatusBadRequest, "user update", "no columns to update")
	}

	return nil
}

// Delete deletes object from the database
func (u *User) Delete(trs *db.Transaction) error {
	if err := db.DeleteStructTx(trs.Tx, constants.TableCoreUsers, &db.Options{
		Conditions: builder.Equal("username", u.Username),
	}); err != nil {
		return customerror.New(http.StatusInternalServerError, "user delete", err.Error())
	}
	return nil
}
