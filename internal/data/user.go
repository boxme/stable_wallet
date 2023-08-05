package data

type User struct {
	Id    uint64
	Email string
	Token string `json:"-"` // Use the - directive
}
