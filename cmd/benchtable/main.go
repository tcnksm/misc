package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
)

var items = []string{"name", "times", "speed", "allocs", "allocs"}

func main() {
	if len(os.Args) != 2 {
		log.Fatal("[Usage] benchtable FILE")
	}

	file, err := os.Open(os.Args[1])
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	fmt.Println("|" + strings.Join(items, "|") + "|")

	str := "|"
	for range items {
		str += " --- |"
	}
	fmt.Println(str)

	sc := bufio.NewScanner(file)
	for sc.Scan() {
		l := sc.Text()
		if l == "PASS" {
			break
		}
		data := strings.Split(l, "\t")
		fmt.Println("|" + strings.Join(data, "|") + "|")
	}
}
