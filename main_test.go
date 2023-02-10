//go:build all
// +build all

package main

import (
	"bytes"
	"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_main(t *testing.T) {
	tests := []struct {
		name string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			main()
		})
	}
}

func Test_readFile(t *testing.T) {
	type args struct {
		arg string
	}
	tests := []struct {
		name             string
		args             args
		wantFileContents []byte
		wantErr          string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotFileContents, err := readFile(tt.args.arg)
			assert.EqualError(t, err, tt.wantErr)
			assert.Equal(t, tt.wantFileContents, gotFileContents)
		})
	}
}

func Test_count(t *testing.T) {
	type args struct {
		fileContents []byte
	}
	tests := []struct {
		name          string
		args          args
		wantLineCount int
		wantWordCount int
		wantCharCount int
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotLineCount, gotWordCount, gotCharCount := count(tt.args.fileContents)
			assert.Equal(t, tt.wantLineCount, gotLineCount)
			assert.Equal(t, tt.wantWordCount, gotWordCount)
			assert.Equal(t, tt.wantCharCount, gotCharCount)
		})
	}
}

func Test_printResult(t *testing.T) {
	type args struct {
		lineCount int
		wordCount int
		charCount int
		file      string
	}
	tests := []struct {
		name string
		args args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			printResult(tt.args.lineCount, tt.args.wordCount, tt.args.charCount, tt.args.file)
		})
	}
}

func Test_printToStderr(t *testing.T) {
	type args struct {
		err error
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "Send a random error",
			args: args{
				err: fmt.Errorf("This is a random error"),
			},
		},
	}
	for _, tt := range tests {
		oldStderr := os.Stderr
		defer func() { os.Stderr = oldStderr }()
		r, w, _ := os.Pipe()
		os.Stderr = w
		t.Run(tt.name, func(t *testing.T) {
			printToStderr(tt.args.err)
			w.Close()
		})

		// Read the output from the pipe
		var buf bytes.Buffer
		_, _ = buf.ReadFrom(r)

		assert.Equal(t, tt.args.err.Error(), buf.String())
	}
}

func Test_printToStdout(t *testing.T) {
	type args struct {
		s string
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "Check if string that is passed gets printed to stdout",
			args: args{
				s: "This is a test string to check printing to os.stdout",
			},
		},
	}
	for _, tt := range tests {
		oldStdout := os.Stdout
		defer func() { os.Stdout = oldStdout }()
		r, w, _ := os.Pipe()
		os.Stdout = w
		t.Run(tt.name, func(t *testing.T) {
			printToStdout(tt.args.s)
			w.Close()
		})
		// Read the output from the pipe
		var buf bytes.Buffer
		_, _ = buf.ReadFrom(r)

		assert.Equal(t, tt.args.s, buf.String())
	}
}
