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
