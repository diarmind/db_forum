package models

type ForumCreate struct {
	Slug		*string		`json:"slug"`
	Title		*string		`json:"title"`
	User		*string		`json:"user"`
}

