package models

type UserResponse struct {
	Authorized bool
	UserID     int64
	Error      string
}
