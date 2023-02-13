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
			lines := make(chan string)
			go worker(lines, arg, &wg)
			wg.Add(1)
		}

		wg.Wait()

		// print total only if more than one file is passed
		if len(args) > 1 {
			totalResult := generateResult(
				totalLineCount,
				totalWordCount,
				totalCharCount,
				"total",
			)
			printToStdout(totalResult)
		}
	},
}

func worker(lines chan string, arg string, wg *sync.WaitGroup) {
	defer wg.Done()
	go readLinesInFile(arg, lines)

	lineCount, wordCount, charCount := count(lines)
	totalLineCount += lineCount
	totalWordCount += wordCount
	totalCharCount += charCount

	result := generateResult(lineCount, wordCount, charCount, arg)
	printToStdout(result)
}

func readLinesInFile(arg string, lines chan<- string) {
	var scanner *bufio.Scanner

	const chunkSize = 1024 * 1024 // 1 MB
	defer close(lines)

	if arg == "-" {
		scanner = bufio.NewScanner(os.Stdin)
	} else {
		file, err := os.Open(arg)
		if err != nil {
			err = fmt.Errorf(
				"gowc: " + strings.Replace(err.Error(), "open ", "", 1) + "\n",
			)
			printToStderr(err)
		}
		defer file.Close()
		scanner = bufio.NewScanner(file)
	}

	scanner.Buffer(make([]byte, chunkSize), chunkSize)

	for scanner.Scan() {
		line := scanner.Text()
		lines <- line
	}

	if err := scanner.Err(); err != nil {
		err = fmt.Errorf(
			"gowc: " + strings.Replace(err.Error(), "read ", "", 1) + "\n",
		)
		printToStderr(err)
	}
}

func count(lines <-chan string) (int, int, int) {
	var lineCount, wordCount, charCount int
	for line := range lines {
		lineCount++
		words := strings.Fields(line)
		wordCount += len(words)
		charCount += len(line) + 1
	}

	return lineCount, wordCount, charCount
}

func generateResult(lineCount, wordCount, charCount int, file string) string {
	var result string
	// append only if lineFlag is set
	if flagSet.lineFlag {
		result += fmt.Sprintf("%8d", lineCount)
	}

	// append only if wordFlag is set
	if flagSet.wordFlag {
		result += fmt.Sprintf("%8d", wordCount)
	}

	// append only if charFlag is set
	if flagSet.charFlag {
		result += fmt.Sprintf("%8d", charCount)
	}

	if !flagSet.lineFlag && !flagSet.wordFlag && !flagSet.charFlag {
		result += fmt.Sprintf("%8d", lineCount)
		result += fmt.Sprintf("%8d", wordCount)
		result += fmt.Sprintf("%8d", charCount)
	}

	// append the filename only if reading from a file intead of os.Stdin after appending the count
	if file == "-" {
		result += "\n"
	} else {
		result += fmt.Sprint(" " + file + "\n")
	}

	return result
}

func printToStderr(err error) {
	fmt.Fprint(os.Stderr, err.Error())
}

func printToStdout(s string) {
	fmt.Fprint(os.Stdout, s)
}
