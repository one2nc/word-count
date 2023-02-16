package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"sync"

	"github.com/spf13/cobra"
)

type flagOptions struct {
	lineFlag bool
	wordFlag bool
	charFlag bool
}

var (
	flagSet                                        flagOptions
	totalLineCount, totalWordCount, totalCharCount int
)

// const maxOpenFileLimit = 1024

func init() {
	// Add flags to count lines, words, and characters
	rootCmd.Flags().BoolVarP(&flagSet.lineFlag, "lines", "l", false, "Count number of lines")
	rootCmd.Flags().BoolVarP(&flagSet.wordFlag, "words", "w", false, "Count number of words")
	rootCmd.Flags().BoolVarP(&flagSet.charFlag, "chars", "c", false, "Count number of characters")
}

func main() {
	// Execute the cobra command
	if err := rootCmd.Execute(); err != nil {
		printToStderr(err)
		os.Exit(1)
	}
}

var rootCmd = &cobra.Command{
	Use:   "wc",
	Short: "wc is a word, line, and character count tool",
	Long:  `wc is a word, line, and character count tool that reads from the standard input or from a file and outputs the number of lines, words, and characters`,
	Run: func(cmd *cobra.Command, args []string) {

		// If length args is equal to '0' then set args as "-" to be identified as os.Stdin
		if len(args) == 0 {
			args = []string{"-"}
		}

		var wg sync.WaitGroup
		for _, arg := range args {
			go worker(arg, &wg)
			wg.Add(1)
		}

		wg.Wait()

		// print total only if more than one file is passed
		if len(args) > 1 {
			totalResult, err := result{
				lineCount: totalLineCount,
				wordCount: totalWordCount,
				charCount: totalCharCount,
				filename:  "total",
			}.generateOutput()
			if err != nil {
				printToStderr(err)
			}
			printToStdout(totalResult)
		}
	},
}

type result struct {
	lineCount int
	wordCount int
	charCount int
	filename  string
	err       error
}

func worker(arg string, wg *sync.WaitGroup) {
	lines := make(chan string)
	errChan := make(chan error)
	defer wg.Done()

	//keep reading lines from files and let the count function run as it happens.
	go readLinesInFile(arg, lines, errChan)

	result := count(lines, errChan)
	totalLineCount += result.lineCount
	totalWordCount += result.wordCount
	totalCharCount += result.charCount
	result.filename = arg

	output, err := result.generateOutput()
	if err != nil {
		printToStderr(err)
		return
	}

	printToStdout(output)
}

func readLinesInFile(filename string, lines chan<- string, errChan chan<- error) {
	var scanner *bufio.Scanner

	const chunkSize = 1024 * 1024 // 1 MB
	defer close(lines)
	defer close(errChan)

	if filename == "-" {
		scanner = bufio.NewScanner(os.Stdin)
		scanner.Buffer(make([]byte, chunkSize), chunkSize)
	} else {
		file, err := os.Open(filename)
		if err != nil {
			err = fmt.Errorf(
				"gowc: " + strings.Replace(err.Error(), "open ", "", 1) + "\n",
			)
			errChan <- err
		}
		defer file.Close()
		scanner = bufio.NewScanner(file)
		scanner.Buffer(make([]byte, chunkSize), chunkSize)
	}

	for scanner.Scan() {
		line := scanner.Text()
		lines <- line // send line to the channel
	}

	if err := scanner.Err(); err != nil {
		err = fmt.Errorf(
			"gowc: " + strings.Replace(err.Error(), "read ", "", 1) + "\n",
		)
		errChan <- err
	}
}

func count(lines <-chan string, errChan <-chan error) result {
	var r result

	for {
		select {
		case err := <-errChan:
			if err != nil {
				r.err = err
				errChan = nil
				return r
			}
		case line, ok := <-lines:
			if !ok {
				return r
			}
			r.lineCount++
			words := strings.Fields(line)
			r.wordCount += len(words)
			r.charCount += len(line) + 1
		}
	}
}

func (r result) generateOutput() (string, error) {
	var output string

	if r.err != nil {
		return "", r.err
	}

	// append only if lineFlag is set
	if flagSet.lineFlag {
		output += fmt.Sprintf("%8d", r.lineCount)
	}

	// append only if wordFlag is set
	if flagSet.wordFlag {
		output += fmt.Sprintf("%8d", r.wordCount)
	}

	// append only if charFlag is set
	if flagSet.charFlag {
		output += fmt.Sprintf("%8d", r.charCount)
	}

	if !flagSet.lineFlag && !flagSet.wordFlag && !flagSet.charFlag {
		output += fmt.Sprintf("%8d", r.lineCount)
		output += fmt.Sprintf("%8d", r.wordCount)
		output += fmt.Sprintf("%8d", r.charCount)
	}

	// append the filename only if reading from a file intead of os.Stdin after appending the count
	if r.filename == "-" {
		output += "\n"
	} else {
		output += fmt.Sprint(" " + r.filename + "\n")
	}

	return output, nil
}

func printToStderr(err error) {
	fmt.Fprint(os.Stderr, err.Error())
}

func printToStdout(s string) {
	fmt.Fprint(os.Stdout, s)
}
