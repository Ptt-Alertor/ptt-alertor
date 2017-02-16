package main

import (
	"io/ioutil"
	"log"

	"github.com/liam-lai/ptt-alertor/myutil"
	"github.com/liam-lai/ptt-alertor/pttboard"
)

func main() {
	//TODO: Read database or setting file
	boards := []string{"free_box", "lol"}
	for _, board := range boards {
		articlesJSON := pttboard.Index(board)
		err := ioutil.WriteFile(myutil.ProjectRootPath()+"/storage/articles/"+board+".json", articlesJSON, 0644)
		if err != nil {
			log.Fatal(err)
		}
	}
}
