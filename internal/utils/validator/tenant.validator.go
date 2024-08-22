package validator

import (
	"antrein/bc-dashboard/model/dto"
	"errors"
)

func ValidateCreateAccount(req dto.CreateTenantRequest) error {
	if !IsEmail(req.Email) {
		return errors.New("Email tidak valid")
	}
	if req.Password != req.RetypePassword {
		return errors.New("Password tidak sama")
	}
	return nil
}
