package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"strings"
)

func main() {
	var (
		linesFlag = flag.Bool("l", false, "count lines")
		wordsFlag = flag.Bool("w", false, "count words")
		charsFlag = flag.Bool("c", false, "count characters")
	)
	flag.Parse()

	args := flag.Args()
	if len(args) == 0 {
		args = []string{"-"}
	}

	for _, arg := range args {
		var file *os.File
		var err error
		if arg == "-" {
			file = os.Stdin
		} else {
			file, err = os.Open(arg)
			if err != nil {
				fmt.Fprintf(os.Stderr, "%s: %s\n", arg, err)
				continue
			}
			defer file.Close()
		}

		var (
			lines int
			words int
			chars int
		)
		scanner := bufio.NewScanner(file)
		for scanner.Scan() {
			lines++
			text := scanner.Text()
			chars += len(text) + 1
			if text == "" {
				continue
			}
			words += len(strings.Fields(text))
		}

		if *linesFlag {
			fmt.Printf("%8d", lines)
		}
		if *wordsFlag {
			fmt.Printf("%8d", words)
		}
		if *charsFlag {
			fmt.Printf("%8d", chars)
		}
		fmt.Printf(" %s\n", arg)
	}
}
