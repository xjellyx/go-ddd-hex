package vo

import "time"

type BaseVO struct {
	ID        string    `json:"-"`
	CreatedAt time.Time `json:"created_at,omitempty"`
	UpdatedAt time.Time `json:"updated_at,omitempty"`
}
