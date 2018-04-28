package models

type PostFull struct {
	Author		*User		`json:"author"`
	Forum		*Forum		`json:"forum"`
	Post		*Post		`json:"post"`
	Thread		*Thread		`json:"thread"`
}
