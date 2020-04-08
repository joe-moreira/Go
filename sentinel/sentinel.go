// to build the binary ("go build") have to run:
// env GOOS=linux GOARCH=amd64 GOARM=7 go build sentinel.go

package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"
)

func main() {
	files, err := ioutil.ReadDir("/etc/redis/")
	if err != nil {
		log.Fatal(err)
	}
	for _, f := range files {
		if strings.Contains(string(f.Name()), "sentinel") {
			fmt.Println(strings.TrimRight(f.Name(), ".conf"))
			t, err := os.OpenFile("/etc/redis/"+f.Name(), os.O_APPEND|os.O_RDWR, 0644)
			if err != nil {
				panic(err)
			}
			defer t.Close()

			//check if logfile word alredy exists on the .conf file
			scanner := bufio.NewScanner(t)
			line := 1

			for scanner.Scan() {
				if strings.Contains(scanner.Text(), "logfile") {
					return
				}
				line++
			}

			if _, err = t.WriteString("logfile /var/log/redis/" + strings.TrimRight(f.Name(), ".conf") + ".log"); err != nil {
				panic(err)
			}

		}
	}
}
