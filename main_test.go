package main

import (
	"bytes"
	"io"
	"os"
	"testing"
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

// TODO: Tests timeout after 30s. Debug this. Tests for : Printing to os.Stdin and Stdout.
func Test_readFile(t *testing.T) {
	type args struct {
		arg string
	}
	tests := []struct {
		name             string
		args             args
		wantFileContents string
		wantErr          bool
		wantStderr       string
	}{
		{
			name: "success, read the contents of test/file.txt",
			args: args{
				arg: "test/file.txt",
			},
			wantFileContents: "A\nB\nC\nD\nE\nF\n",
			wantErr:          false,
		},
		{
			name: "failure, no read permission on the file test/no_permissions.txt",
			args: args{
				arg: "test/no_permissions.txt",
			},
			wantFileContents: "",
			wantErr:          true,
			wantStderr:       "wc: test/no_permissions.txt: no such file or directory\n",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.wantErr {
				oldStderr := os.Stderr
				r, w, _ := os.Pipe()
				os.Stderr = w

				readFile(tt.args.arg)

				var buf bytes.Buffer
				io.Copy(&buf, r)
				defer func() {
					w.Close()
					os.Stderr = oldStderr
				}()

				if buf.String() != tt.wantStderr {
					t.Errorf(
						"readFile() = \n%v, want : \n%v",
						string(buf.String()),
						tt.wantStderr,
					)
				}

			} else {
				gotFileContents, err := readFile(tt.args.arg)
				if (err != nil) != tt.wantErr {
					t.Errorf("readFile() error = %v, wantErr %v", err, tt.wantErr)
					return
				}
				if string(gotFileContents) != string(tt.wantFileContents) {
					t.Errorf(
						"readFile() = \n%v, want : \n%v",
						string(gotFileContents),
						tt.wantFileContents,
					)
				}
			}
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
			if gotLineCount != tt.wantLineCount {
				t.Errorf("count() gotLineCount = %v, want %v", gotLineCount, tt.wantLineCount)
			}
			if gotWordCount != tt.wantWordCount {
				t.Errorf("count() gotWordCount = %v, want %v", gotWordCount, tt.wantWordCount)
			}
			if gotCharCount != tt.wantCharCount {
				t.Errorf("count() gotCharCount = %v, want %v", gotCharCount, tt.wantCharCount)
			}
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
