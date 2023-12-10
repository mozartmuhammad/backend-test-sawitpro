// This file contains types that are used in the repository layer.
package repository

import (
	"time"

	"github.com/SawitProRecruitment/UserService/common"
)

type RegisterUser struct {
	Name     string `json:"name"`
	Password string `json:"password"`
	Phone    string `json:"phone"`
}

func (t *RegisterUser) Validate() (result []string) {
	result = append(result, common.ValidatePhone(t.Phone)...)
	result = append(result, common.ValidateName(t.Name)...)
	result = append(result, common.ValidatePassword(t.Password)...)

	return
}

type User struct {
	ID        int64
	Phone     string
	Name      string
	Password  string
	CreatedAt time.Time
	UpdateAt  *time.Time
}

func (t *User) Validate() (result []string) {
	result = append(result, common.ValidatePhone(t.Phone)...)
	result = append(result, common.ValidateName(t.Name)...)

	return result
}
