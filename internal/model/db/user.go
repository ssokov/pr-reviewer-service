package db

import "time"

type User struct {
	ID        int64
	UserID    string
	Username  string
	TeamID    int64
	IsActive  bool
	CreatedAt time.Time
}
