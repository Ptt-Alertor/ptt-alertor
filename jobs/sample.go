package jobs

import (
	"fmt"
	"time"
)

type Sample struct {
}

func (s Sample) Run() {
	fmt.Println(time.Now())
}
