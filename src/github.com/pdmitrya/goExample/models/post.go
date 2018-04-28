package models

import "time"

type Post struct {
	Author		*string			`json:"author"`
	Created		*time.Time		`json:"created"`
	Forum		*string			`json:"forum"`
	Id			*int			`json:"id"`
	IsEdited	*bool			`json:"isEdited"`
	Message		*string			`json:"message"`
	Parent		*int			`json:"parent"`
	Thread		*int			`json:"thread"`
}
