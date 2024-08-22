package validator

import (
	"regexp"
	"strings"
)

func IsEmail(email string) bool {
	emailRegex := `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`
	regex := regexp.MustCompile(emailRegex)
	return regex.MatchString(email)

}

func IsPhoneNumber(phoneNumber string) bool {
	phoneNumberRegex := `^\+\d{9,15}$`
	regex := regexp.MustCompile(phoneNumberRegex)
	return regex.MatchString(phoneNumber)
}

func IsAlphanumeric(input string) bool {
	alphanumericRegex := "^[a-zA-Z0-9]+$"
	regex := regexp.MustCompile(alphanumericRegex)
	return regex.MatchString(input)
}

func IsAlphanumericWithSpace(input string) bool {
	alphanumericRegex := "^[a-zA-Z0-9 ]+$"
	regex := regexp.MustCompile(alphanumericRegex)
	return regex.MatchString(input)
}

func IsPersonName(input string) bool {
	nameRegex := `^[a-zA-Z\s'.]+$`
	regex := regexp.MustCompile(nameRegex)
	return regex.MatchString(input)
}

func ContainsSpace(input string) bool {
	return strings.Contains(input, " ")
}

func IsUsername(input string) bool {
	usernameRegex := "^[a-z0-9_-]{5,}$"
	regex := regexp.MustCompile(usernameRegex)
	return regex.MatchString(input)
}
