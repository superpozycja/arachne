package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"strings"
)

var codebaseDir = flag.String("d", "", "directory of the codebase to crawl")

func main() {
	flag.Parse()
	if *codebaseDir == "" {
		log.Fatal("no dir provided")
	}
	files, err := ioutil.ReadDir(*codebaseDir)
	if err != nil {
		log.Fatal(err)
	}
	for _, file := range files {
		if !file.IsDir() {
			fmt.Println(file.Name())
			fmt.Println("")
			b, err := ioutil.ReadFile(*codebaseDir + file.Name())
			if err != nil {
				log.Fatal(err)
			}
			data := strings.Split(string(b), "\n")
			for _, line := range data {
				if strings.Contains(line, "import") {
					fmt.Println(line)
				}
			}
			fmt.Println("")
			fmt.Println("")
		}
	}
}
