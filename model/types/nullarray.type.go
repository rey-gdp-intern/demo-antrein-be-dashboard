package types

import (
	"antrein/bc-dashboard/internal/utils/parser"
	"errors"
)

type NullStringArray struct {
	StringArray []string
	Valid       bool
}

func (a *NullStringArray) Scan(value interface{}) error {
	if value == nil {
		a.StringArray, a.Valid = nil, false
		return nil
	}
	str, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}

	parsedArray, err := parser.ParseStringArray(string(str))
	if err != nil {
		return err
	}

	a.StringArray, a.Valid = parsedArray, true
	return nil
}
