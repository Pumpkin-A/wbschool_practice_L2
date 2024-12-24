package main

import (
	"bufio"
	"bytes"
	"os"
	"strings"
	"testing"
)

var testDataDir string = "testdata"

func readFileOrPanic(filename string) []string {
	file, err := os.Open(testDataDir + "/" + filename)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	lines := []string{}
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	if err := scanner.Err(); err != nil {
		panic(err)
	}
	return lines
}

func Test_doGrep(t *testing.T) {
	type args struct {
		data   []string
		params parametres
	}
	tests := []struct {
		name    string
		args    args
		wantOut string
	}{
		{
			name: "Delim space First column",
			args: args{
				data:   readFileOrPanic("file.txt"),
				params: parametres{delim: " ", fields: []int{0}},
			},
			wantOut: strings.Join(readFileOrPanic("file_ans_1.txt"), "\n"),
		},
		{
			name: "Delim space First column",
			args: args{
				data:   readFileOrPanic("file.txt"),
				params: parametres{delim: ":", fields: []int{0, 1}, withDelimOnly: true},
			},
			wantOut: strings.Join(readFileOrPanic("file_ans_2.txt"), "\n"),
		},
		{
			name: "Delim space First column",
			args: args{
				data:   readFileOrPanic("file.txt"),
				params: parametres{delim: " ", fields: []int{0, 1, 2}, withDelimOnly: true},
			},
			wantOut: strings.Join(readFileOrPanic("file_ans_3.txt"), "\n"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			out := &bytes.Buffer{}
			doCut(tt.args.data, tt.args.params, out)
			if gotOut := out.String()[:len(out.String())-1]; gotOut != tt.wantOut {
				t.Errorf("doGrep() = %v, want %v", gotOut, tt.wantOut)
			}
		})
	}
}
