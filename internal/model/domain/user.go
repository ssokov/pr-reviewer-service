package domain

import "time"

type User struct {
	ID        int64
	UserID    string
	Username  string
	TeamID    int64
	TeamName  string
	IsActive  bool
	CreatedAt time.Time
}
