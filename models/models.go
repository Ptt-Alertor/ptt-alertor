package models

import (
	"github.com/meifamily/ptt-alertor/models/article"
	"github.com/meifamily/ptt-alertor/models/user"
)

var User = user.NewUser(new(user.Redis))
var Article = article.NewArticle(new(article.Redis))
