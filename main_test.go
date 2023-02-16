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

func Test_count(t *testing.T) {
	// Test empty input
	emptyLines := make(chan string)
	emptyErrChan := make(chan error)
	go func() {
		close(emptyLines)
	}()
	res := count(emptyLines, emptyErrChan)
	if res.lineCount != 0 || res.wordCount != 0 || res.charCount != 0 || res.err != nil {
		t.Errorf(
			"count(emptyLines, emptyErrChan) = %v; expected lineCount=0, wordCount=0, charCount=0, err=nil",
			res,
		)
	}

	// Test input with one line
	oneLine := make(chan string)
	oneErrChan := make(chan error)
	go func() {
		oneLine <- "Hello, world!"
		close(oneLine)
	}()
	res = count(oneLine, oneErrChan)
	if res.lineCount != 1 || res.wordCount != 2 || res.charCount != 14 || res.err != nil {
		t.Errorf(
			"count(oneLine, oneErrChan) = %v; expected lineCount=1, wordCount=2, charCount=14, err=nil",
			res,
		)
	}

	// Test input with multiple lines
	multiLine := make(chan string)
	multiErrChan := make(chan error)
	go func() {
		multiLine <- "The quick brown fox"
		multiLine <- "jumps over the lazy dog."
		close(multiLine)
	}()
	res = count(multiLine, multiErrChan)
	if res.lineCount != 2 || res.wordCount != 9 || res.charCount != 45 || res.err != nil {
		t.Errorf(
			"count(multiLine, multiErrChan) = %v; expected lineCount=2, wordCount=9, charCount=42, err=nil",
			res,
		)
	}

	// Test input with error
	errorLine := make(chan string)
	errorErrChan := make(chan error)
	expectedErr := fmt.Errorf("test error")
	go func() {
		errorErrChan <- expectedErr
		close(errorLine)
	}()
	res = count(errorLine, errorErrChan)
	if res.lineCount != 0 || res.wordCount != 0 || res.charCount != 0 || res.err != expectedErr {
		t.Errorf(
			"count(errorLine, errorErrChan) = %v; expected lineCount=0, wordCount=0, charCount=0, err=%v",
			res,
			expectedErr,
		)
	}
}

func Test_result_generateOutput(t *testing.T) {
	type fields struct {
		lineCount int
		wordCount int
		charCount int
		filename  string
		err       error
	}

	tests := []struct {
		name      string
		fields    fields
		testFlags flagOptions
		want      string
		wantErr   bool
	}{
		{
			name: "success, all flags are set",
			fields: fields{
				lineCount: 3,
				wordCount: 6,
				charCount: 12,
				filename:  "test.txt",
				err:       nil,
			},
			testFlags: flagOptions{
				lineFlag: true,
				wordFlag: true,
				charFlag: true,
			},
			want:    "       3       6      12 test.txt\n",
			wantErr: false,
		},
		{
			name: "success, only lines flag is set",
			fields: fields{
				lineCount: 3,
				wordCount: 6,
				charCount: 12,
				filename:  "test.txt",
				err:       nil,
			},
			testFlags: flagOptions{
				lineFlag: true,
			},
			want:    "       3 test.txt\n",
			wantErr: false,
		},
		{
			name: "success, only lines and words flag are set",
			fields: fields{
				lineCount: 3,
				wordCount: 6,
				charCount: 12,
				filename:  "test.txt",
				err:       nil,
			},
			testFlags: flagOptions{
				lineFlag: true,
				wordFlag: true,
			},
			want:    "       3       6 test.txt\n",
			wantErr: false,
		},
		{
			name: "success, only words and chars flag is set",
			fields: fields{
				lineCount: 3,
				wordCount: 6,
				charCount: 12,
				filename:  "test.txt",
				err:       nil,
			},
			testFlags: flagOptions{
				lineFlag: true,
				charFlag: true,
			},
			want:    "       3      12 test.txt\n",
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := result{
				lineCount: tt.fields.lineCount,
				wordCount: tt.fields.wordCount,
				charCount: tt.fields.charCount,
				filename:  tt.fields.filename,
				err:       tt.fields.err,
			}

			flagSet = tt.testFlags

			got, err := r.generateOutput()
			if tt.wantErr {
				assert.EqualError(t, err, err.Error(), tt.fields.err)
			}
			assert.Equal(t, tt.want, got)
		})
	}
}
