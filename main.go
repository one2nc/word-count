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

	var (
		lines      int
		words      int
		chars      int
		totalLines int
		totalWords int
		totalChars int
	)

	for _, arg := range args {
		var file *os.File
		var err error
		if arg == "-" {
			file = os.Stdin
		} else {
			file, err = os.Open(arg)
			if err != nil {
				fmt.Fprintf(os.Stderr, "./wc: %s\n", err)
				return
			}
			defer file.Close()

			fileInfo, err := file.Stat()
			if err != nil {
				fmt.Fprintf(os.Stderr, "./wc: %s\n", err)
				continue
			}
			if fileInfo.IsDir() {
				fmt.Fprintf(os.Stderr, "./wc %s: read: is a directory\n", arg)
				continue
			}
		}

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

		if len(args) > 1 {
			totalLines += lines
			totalWords += words
			totalChars += chars
			lines, words, chars = 0, 0, 0
		}
	}

	if len(args) > 1 {
		if *linesFlag {
			fmt.Printf("%8d", totalLines)
		}
		if *wordsFlag {
			fmt.Printf("%8d", totalWords)
		}
		if *charsFlag {
			fmt.Printf("%8d", totalChars)
		}
		fmt.Printf(" total\n")
	}

}
