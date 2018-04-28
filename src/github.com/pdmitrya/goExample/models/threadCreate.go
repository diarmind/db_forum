package models

import "time"

type ThreadCreate struct {
	Author		*string		`json:"author"`
	Created		*time.Time	`json:"created"`
	Message		*string		`json:"message"`
	Title		*string		`json:"title"`
	Slug		*string		`json:"slug"`
}
