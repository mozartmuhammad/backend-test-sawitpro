// This file contains types that are used in the repository layer.
package repository

import (
	"time"

	"github.com/SawitProRecruitment/UserService/common"
)

// A RegisterUser represents input for creating new users.
type RegisterUser struct {
	Name     string
	Password string
	Phone    string
}

// Validate validate user registration input.
// Return empty result if comply with given rules.
func (t *RegisterUser) Validate() (result []string) {
	result = append(result, common.ValidatePhone(t.Phone)...)
	result = append(result, common.ValidateName(t.Name)...)
	result = append(result, common.ValidatePassword(t.Password)...)

	return
}

// An User represents an user data in this application.
type User struct {
	ID        int64
	Phone     string
	Name      string
	Password  string
	CreatedAt time.Time
	UpdateAt  *time.Time
}

// Validate validate user update input.
// Return empty result if comply with given rules.
func (t *User) Validate() (result []string) {
	result = append(result, common.ValidatePhone(t.Phone)...)
	result = append(result, common.ValidateName(t.Name)...)

	return result
}
