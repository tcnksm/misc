/*
Command 'benchtable' generates a markdown table from go bench results.
You can provide benchmark result via stdin or a file.

  $ go test -bench . -benchmem | benchtable

See example output on https://gist.github.com/tcnksm/207e60f2e39c2f9b29d6082b1ea020e7

To install it,

  $ go get github.com/tcnksm/misc/cmd/benchtable

*/
package main

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"log"
	"os"
	"regexp"
	"strings"
)

var items = []string{"name", "times", "ns/op", "B/op", "allocs/op"}

var reNum = regexp.MustCompile(`\d+`)

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
	str := "| :"
	for i, _ := range items {
		if i == 0 {
			str += "---: |"
			continue
		}
		str += " ---: |"
	}
	fmt.Println(str)

	for sc.Scan() {
		l := sc.Text()
		if l == "PASS" {
			break
		}

		var buf bytes.Buffer
		fmt.Fprintf(&buf, "|")
		data := strings.Split(l, "\t")
		for j, d := range data {
			if j == 0 {
				fmt.Fprintf(&buf, "%s|", strings.TrimSpace(d))
				continue
			}

			fmt.Fprintf(&buf, "%s|", reNum.FindString(d))
		}
		fmt.Println(buf.String())
	}
}
