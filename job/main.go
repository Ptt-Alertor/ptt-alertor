package main

import (
	"fmt"
	"log"
	"os/exec"

	"github.com/liam-lai/ptt-alertor/myutil"
)

func main() {
	fmt.Println("start job")
	filePath := myutil.ProjectRootPath() + "/job/sendmail/send_new_articles.go"
	cmd := exec.Command("go run " + filePath)
	err := cmd.Run()
	if err != nil {
		log.Fatal(err)
	}
	err = cmd.Wait()

	filePath = myutil.ProjectRootPath() + "/job/fetcharticles/fetch_newest_articles.go"
	cmd = exec.Command("go run " + filePath)
	err = cmd.Run()
	if err != nil {
		log.Fatal(err)
	}
	err = cmd.Wait()
	fmt.Println("finish")
}
