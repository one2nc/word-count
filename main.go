package main

import (
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "wc",
	Short: "wc is a word, line, and character count tool",
	Long:  `wc is a word, line, and character count tool that reads from the standard input or from a file and outputs the number of lines, words, and characters`,
	Run: func(cmd *cobra.Command, args []string) {
		var totalLineCount, totalWordCount, totalCharCount int

		// If length args is equal to '0' then set args as "-" to be identified as os.Stdin
		if len(args) == 0 {
			args = []string{"-"}
		}

		for _, arg := range args {
			fileContents, err := readFile(arg)
			if err != nil {
				printToStderr(err)
				return
			}

			lineCount, wordCount, charCount := count(fileContents)
			totalLineCount += lineCount
			totalWordCount += wordCount
			totalCharCount += charCount

			printResult(lineCount, wordCount, charCount, arg)
		}

		// print total only if more than one file is passed
		if len(args) > 1 {
			printResult(
				totalLineCount,
				totalWordCount,
				totalCharCount,
				"total",
			)
		}
	},
}

type flagOptions struct {
	lineFlag bool
	wordFlag bool
	charFlag bool
}

var flagSet flagOptions

func init() {
	// Add flags to count lines, words, and characters
	rootCmd.Flags().BoolVarP(&flagSet.lineFlag, "lines", "l", false, "Count number of lines")
	rootCmd.Flags().BoolVarP(&flagSet.charFlag, "words", "w", false, "Count number of words")
	rootCmd.Flags().BoolVarP(&flagSet.wordFlag, "chars", "c", false, "Count number of characters")
}

func main() {
	// Execute the cobra command
	if err := rootCmd.Execute(); err != nil {
		printToStderr(err)
		os.Exit(1)
	}
}

func readFile(arg string) (fileContents []byte, err error) {
	if arg == "-" {
		fileContents, err = io.ReadAll(os.Stdin)
		if err != nil {
			err = fmt.Errorf(
				"wc: " + strings.Replace(err.Error(), "open ", "", 1) + "\n",
			)
			return
		}
		return
	}

	fileContents, err = os.ReadFile(arg)
	if err != nil {
		err = fmt.Errorf(
			"wc: " + strings.Replace(err.Error(), "open ", "", 1) + "\n",
		)
		return
	}
	return
}

func count(fileContents []byte) (lineCount, wordCount, charCount int) {
	charCount = len(fileContents)
	lineCount = strings.Count(string(fileContents), "\n")
	wordCount = len(strings.Fields(string(fileContents)))
	return
}

func printResult(lineCount, wordCount, charCount int, file string) {
	// print only if lineFlag is set
	if flagSet.lineFlag {
		printToStdout(fmt.Sprintf("%8d", lineCount))
	}

	// print only if wordFlag is set
	if flagSet.wordFlag {
		printToStdout(fmt.Sprintf("%8d", wordCount))
	}

	// print only if charFlag is set
	if flagSet.charFlag {
		printToStdout(fmt.Sprintf("%8d", charCount))
	}

	// print the filename only if reading from a file intead of os.Stdin after printing the count
	if file == "-" {
		printToStdout("\n")
	} else {
		printToStdout(fmt.Sprint(" " + file + "\n"))
	}
}

func printToStderr(err error) {
	fmt.Fprint(os.Stderr, err.Error())
}

func printToStdout(s string) {
	fmt.Fprint(os.Stdout, s)
}
