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

type result struct {
	lineCount int
	wordCount int
	charCount int
	filename  string
	err       error
}

var (
	flagSet                                        flagOptions
	totalLineCount, totalWordCount, totalCharCount int
)

const maxOpenFileLimit = 1024

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
			lines := make(chan string, maxOpenFileLimit)
			errChan := make(chan error)
			go worker(lines, errChan, arg, &wg)
			wg.Add(1)
		}

		wg.Wait()

		// print total only if more than one file is passed
		if len(args) > 1 {
			totalResult, err := generateResult(
				result{
					lineCount: totalLineCount,
					wordCount: totalWordCount,
					charCount: totalCharCount,
					filename:  "total",
				},
			)
			if err != nil {
				printToStderr(err)
			}
			printToStdout(totalResult)
		}
	},
}

func worker(lines chan string, errChan chan error, arg string, wg *sync.WaitGroup) {
	defer wg.Done()
	go readLinesInFile(arg, lines, errChan)

	countResult := count(lines, errChan)
	totalLineCount += countResult.lineCount
	totalWordCount += countResult.wordCount
	totalCharCount += countResult.charCount
	countResult.filename = arg

	result, err := generateResult(countResult)
	if err != nil {
		printToStderr(err)
	}
	printToStdout(result)
}

func readLinesInFile(arg string, lines chan<- string, errChan chan<- error) {
	var scanner *bufio.Scanner

	const chunkSize = 1024 * 1024 // 1 MB
	defer close(lines)
	defer close(errChan)

	if arg == "-" {
		scanner = bufio.NewScanner(os.Stdin)
		scanner.Buffer(make([]byte, chunkSize), chunkSize)

		for scanner.Scan() {
			line := scanner.Text()
			lines <- line
		}

		if err := scanner.Err(); err != nil {
			err = fmt.Errorf(
				"gowc: " + strings.Replace(err.Error(), "read ", "", 1) + "\n",
			)
			errChan <- err
		}
	} else {
		file, err := os.Open(arg)
		if err != nil {
			err = fmt.Errorf(
				"gowc: " + strings.Replace(err.Error(), "open ", "", 1) + "\n",
			)
			errChan <- err
		}
		defer file.Close()
		scanner = bufio.NewScanner(file)
		scanner.Buffer(make([]byte, chunkSize), chunkSize)

		for scanner.Scan() {
			line := scanner.Text()
			lines <- line
		}

		if err := scanner.Err(); err != nil {
			err = fmt.Errorf(
				"gowc: " + strings.Replace(err.Error(), "read ", "", 1) + "\n",
			)
			errChan <- err
		}
	}
}

func count(lines <-chan string, errChan <-chan error) result {
	var res result

	for err := range errChan {
		if err != nil {
			res.err = err
			return res
		}
	}

	for line := range lines {
		res.lineCount++
		words := strings.Fields(line)
		res.wordCount += len(words)
		res.charCount += len(line) + 1
	}

	return res
}

func generateResult(r result) (string, error) {
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
