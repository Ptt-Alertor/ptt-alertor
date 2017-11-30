//+build !test

package models

import "github.com/meifamily/ptt-alertor/models/user"

var User = user.NewUser(new(user.Redis))
