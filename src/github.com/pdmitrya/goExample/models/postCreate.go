package models

type PostCreate struct {
	Author		*string		`json:"author"`
	Message		*string		`json:"message"`
	Parent		*int		`json:"parent"`
}

