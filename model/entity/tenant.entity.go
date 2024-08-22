package entity

import (
	"database/sql"
	"time"
)

type Tenant struct {
	ID        string       `db:"id"`
	Email     string       `db:"email"`
	Password  string       `db:"password"`
	Name      string       `db:"name"`
	CreatedAt time.Time    `db:"created_at"`
	UpdatedAt sql.NullTime `db:"updated_at,omitempty"`
}
