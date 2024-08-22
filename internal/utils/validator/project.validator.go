package validator

import (
	"antrein/bc-dashboard/model/dto"
	"errors"
)

func ValidateCreateProject(req dto.CreateProjectRequest) error {
	if !IsUsername(req.ID) {
		return errors.New("ID project minimal 5 karakter, terdiri dari huruf kecil, angka, underscore(_) dan strip(-)")
	}
	return nil
}
