package models

type User struct {
	ID       uint
	Name     string
	Login    string
	HashPass string
	Role     string
}
