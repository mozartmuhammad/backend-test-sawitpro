package common

import (
	"regexp"
	"strings"
	"unicode"
)

func ValidatePhone(phone string) (errMsg []string) {
	if !strings.HasPrefix(phone, "+62") {
		errMsg = append(errMsg, "Phone must start with +62")
	}
	if isMatch, _ := regexp.MatchString("[0-9]$", phone[3:]); !isMatch {
		errMsg = append(errMsg, "Phone must number only")
	}
	if len(phone) < 10 {
		errMsg = append(errMsg, "Phone must be greater than 10")
	}
	if len(phone) > 13 {
		errMsg = append(errMsg, "Phone must be less than 13")
	}
	return
}

func ValidateName(name string) (errMsg []string) {
	if len(name) < 3 {
		errMsg = append(errMsg, "Name must be greater than 3")
	}
	if len(name) > 60 {
		errMsg = append(errMsg, "Name must be less than 60")
	}
	return
}

func ValidatePassword(password string) (errMsg []string) {
	if len(password) < 6 {
		errMsg = append(errMsg, "Password must be greater than 6")
	}
	if len(password) > 64 {
		errMsg = append(errMsg, "Password must be less than 64")
	}

	var number, upper, special int
	for _, c := range password {
		switch {
		case unicode.IsNumber(c):
			number++
		case unicode.IsUpper(c):
			upper++
		case unicode.IsPunct(c) || unicode.IsSymbol(c):
			special++
		}
	}
	if number == 0 || upper == 0 || special == 0 {
		errMsg = append(errMsg, "Password must have at least 1 capital letter, 1 number, and 1 special character")
	}
	return
}
