package user

type User struct {
	Enable  bool `json:"enable"`
	Profile struct {
		Account string `json:"account"`
		Email   string `json:"email"`
		Line    string `json:"line"`
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
	Update() error
	Find() User
}
