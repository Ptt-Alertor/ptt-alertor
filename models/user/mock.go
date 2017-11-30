//+build test

package user

type Mock struct{}

func (Mock) List() (accounts []string) {
	return []string{"dinos80152@gmail.com"}
}

func (Mock) Exist(account string) bool {
	if account == "dinos80152@gmail.com" {
		return true
	}
	return false
}

func (Mock) Save(account string, user interface{}) error {
	return nil
}

func (Mock) Update(account string, user interface{}) error {
	return nil
}

func (Mock) Find(account string, user *User) {
	user.Profile.Account = "dinos80152@gmail.com"
	return
}
