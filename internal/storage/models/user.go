package models

type User struct {
	ID       int64
	Username string
	Password string
}

type UserProfile struct {
	Username    string
	Description string
}
