package article

import (
	"fmt"
	"time"
)

type Comment struct {
	Tag      string
	UserID   string
	Content  string
	DateTime time.Time
}

func (c Comment) String() string {
	// 推 ChoDino: 推文推文
	return fmt.Sprintf("%s %s%s", c.Tag, c.UserID, c.Content)
}

type Comments []Comment

func (cs Comments) String() string {
	var content string
	for _, p := range cs {
		content += "\n" + p.String()
	}
	return content
}
