package domain

import "time"

type Team struct {
	ID        int64
	TeamName  string
	Members   []User
	CreatedAt time.Time
}
