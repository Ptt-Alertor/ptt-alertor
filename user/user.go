package user

type User interface {
	All() []*User
	Save() error
	Find() User
}
