package db

import "time"

type Team struct {
	ID        int64
	TeamName  string
	CreatedAt time.Time
}
