package user

type User struct {
	Profile struct {
		Account string `json:"account"`
		Email   string `json:"email"`
	}
	Subscribes []Subscribe
}

type Subscribe struct {
	Board    string
	Keywords []string
}

type UserAction interface {
	All() []*User
	Save() error
	Find() User
}
