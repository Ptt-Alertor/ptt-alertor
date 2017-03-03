package user

type User interface {
	All() Users
	Save() error
	Find() User
}
