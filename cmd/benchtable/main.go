/*
Command 'benchtable' generates a markdown table from go bench results.
You can provide benchmark result via stdin or a file.

  $ go test -bench . -benchmem | benchtable

To install it,

  $ go get github.com/tcnksm/misc/cmd/benchtable

*/
package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
)

var items = []string{"name", "times", "speed", "allocs", "allocs"}

func main() {
	var rd io.Reader
	if len(os.Args) == 2 {
		file, err := os.Open(os.Args[1])
		if err != nil {
			log.Fatal(err)
		}
		defer file.Close()
		rd = file
	} else {
		rd = os.Stdin
	}
	sc := bufio.NewScanner(rd)

	fmt.Println("| " + strings.Join(items, " | ") + " |")
	str := "|"
	for _, item := range items {
		str += " " + strings.Repeat("-", len(item)) + " |"
	}
	fmt.Println(str)

	for sc.Scan() {
		l := sc.Text()
		if l == "PASS" {
			break
		}
		data := strings.Split(l, "\t")
		fmt.Println("|" + strings.Join(data, "|") + "|")
	}
}
