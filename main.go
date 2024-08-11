package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"io/ioutil"
	//"io/fs"
	"log"
	"strings"
	"regexp"
	"sort"
)

var codebaseDir = flag.String("d", "", "directory of the codebase to crawl")

var fileGraph = map[string][]string{}

func topoSort(m map[string] []string) []string {
	var order []string
	seen := make(map[string]bool)
	var visitAll func(items []string)
	visitAll = func(items []string) {
		for _, item := range items {
			if !seen[item] {
				seen[item] = true
				visitAll(m[item])
				order = append(order, item)
			}
		}
	}
	var keys []string
	for key := range m {
		keys = append(keys, key)
	}

	sort.Strings(keys)
	visitAll(keys)
	return order
}

func getImports(path string, file os.FileInfo, err error) error {
	if !file.IsDir() {
		b, err := ioutil.ReadFile(path)
		if err != nil {
			log.Fatal(err)
		}
		s := []string{}
		data := strings.Split(string(b), "\n")
		for _, line := range data {
			matched, _ := regexp.MatchString(`import \{[A-Za-z]+\} from ["'][^"]*["'];`, line)
			matched2, _ := regexp.MatchString(`import ["'][^"]*["'];`, line)
			if (matched || matched2) && !strings.Contains(line, "@") {
				re := regexp.MustCompile(`(?s)["'].*?["']`);
				match := re.FindStringSubmatch(line)
				match[0] = strings.TrimSuffix(match[0], "\"")
				match[0] = strings.TrimPrefix(match[0], "\"")
				match[0] = strings.TrimSuffix(match[0], "'")
				match[0] = strings.TrimPrefix(match[0], "'")
				dir := path[:strings.LastIndex(path, "/")+1]
				n := filepath.Clean(dir + "/" + match[0])
				//fmt.Println(n)
				s = append(s, n)
			}
		}
		n, _ := filepath.Abs(path)
		//fmt.Println(n)
		fileGraph[n] = s
	}
	return err
}

func relPath(path string) string {
	abspath, err := filepath.Abs(*codebaseDir)
	res, err := filepath.Rel(abspath, path)
	if err != nil {
		return ""
	}
	return res
}

func main() {
	flag.Parse()
	if *codebaseDir == "" {
		log.Fatal("no dir provided")
	}
	err := filepath.Walk(*codebaseDir, getImports)
	//files, err := os.ReadDir(*codebaseDir)
	if err != nil {
		log.Fatal(err)
	}

	res := topoSort(fileGraph)

	for _, val := range res {
		fmt.Println(relPath(val))
		fmt.Println("\t\t|")
		fmt.Println("\t\tv")
		/*
		for _, val2 := range fileGraph[val] {
			fmt.Println("|")
			fmt.Println("\\----> ", relPath(val2))
		}
		*/
	}
	/*
	for _, file := range files {
	}
	*/
}
