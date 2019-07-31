package user

import (
	"net/http"

	"github.com/agile-work/srv-mdl-shared/models/customerror"
	"github.com/agile-work/srv-shared/constants"
	"github.com/agile-work/srv-shared/sql-builder/builder"
	"github.com/agile-work/srv-shared/sql-builder/db"
	"github.com/agile-work/srv-shared/token"
)

// Login validate credentials and return user token
func (u *User) Login() error {
	if u.Email == "" || u.Password == "" {
		return customerror.New(http.StatusBadRequest, "user login", "invalid credentials body")
	}

	password := u.Password
	if err := db.SelectStruct(constants.TableCoreUsers, u, &db.Options{
		Conditions: builder.Equal("email", u.Email),
	}); err != nil {
		return customerror.New(http.StatusInternalServerError, "user login load user", err.Error())
	}

	if u.ID == "" {
		return customerror.New(http.StatusNotFound, "user login", "user not found with this email")
	}

	if u.Password != password {
		return customerror.New(http.StatusUnauthorized, "user login", "invalid password")
	}

	if !u.Active {
		return customerror.New(http.StatusUnauthorized, "user login", "deactivated user")
	}

	u.Password = ""
	u.Security = nil
	u.SecurityInstances = nil

	payload := make(map[string]interface{})
	payload["code"] = u.Username
	payload["language_code"] = u.LanguageCode

	tokenString, err := token.New(payload, 2*constants.Hour)
	if err != nil {
		return customerror.New(http.StatusInternalServerError, "user login", err.Error())
	}

	u.Token = tokenString

	return nil
}
