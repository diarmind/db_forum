package models

type UserUpdate struct {
	About		*string		`json:"about"`
	Email		*string		`json:"email"`
	Fullname	*string		`json:"fullname"`
}
