package models

import (
	"time"
)

type Message struct {
	ID         int64
	FromID     int64
	Text       string
	CreateTime time.Time
}
