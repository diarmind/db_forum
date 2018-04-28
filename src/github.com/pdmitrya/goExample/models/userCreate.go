package models

type UserCreate struct {
	About		*string		`json:"about"`
	Email		*string		`json:"email"`
	Fullname	*string		`json:"fullname"`
}
