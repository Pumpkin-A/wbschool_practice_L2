package main

import (
	"bufio"
	"bytes"
	"os"
	"strings"
	"testing"
)

func readFileOrPanic(filename string) []string {
	file, err := os.Open(filename)
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
			name: "default Find USD",
			args: args{
				data:   readFileOrPanic("testFiles/default_1.txt"),
				params: parametres{pattern: "USD"},
			},
			wantOut: strings.Join(readFileOrPanic("testFiles/default_1_ans.txt"), "\n"),
		},
		{
			name: "default Find last eur (eur$)",
			args: args{
				data:   readFileOrPanic("testFiles/default_2.txt"),
				params: parametres{pattern: "eur$"},
			},
			wantOut: strings.Join(readFileOrPanic("testFiles/default_2_ans.txt"), "\n"),
		},
		{
			name: "default_1 Count",
			args: args{
				data:   readFileOrPanic("testFiles/default_1.txt"),
				params: parametres{pattern: "u", toCount: true},
			},
			wantOut: "3\n",
		},
		{
			name: "default_1 Count_ignore_case",
			args: args{
				data:   readFileOrPanic("testFiles/default_2.txt"),
				params: parametres{pattern: "usd", toCount: true, ignoreCase: true},
			},
			wantOut: "4\n",
		},
		{
			name: "default_1 Count_ignore_case",
			args: args{
				data:   readFileOrPanic("testFiles/default_2.txt"),
				params: parametres{pattern: "usd", toCount: true, ignoreCase: true, invert: true},
			},
			wantOut: "2\n",
		},
		{
			name: "fixed no pattern",
			args: args{
				data:   readFileOrPanic("testFiles/default_2.txt"),
				params: parametres{pattern: "byn$", fixed: true},
			},
			wantOut: "",
		},
		// check output flags -A -B -C -n
		{
			name: "out -c ",
			args: args{
				data:   readFileOrPanic("testFiles/out.txt"),
				params: parametres{pattern: "bitoc$", after: 2, before: 2},
			},
			wantOut: strings.Join(readFileOrPanic("testFiles/out_c_ans.txt"), "\n"),
		},
		{
			name: "out -C ",
			args: args{
				data:   readFileOrPanic("testFiles/out.txt"),
				params: parametres{pattern: "bitoc$", after: 2, before: 2},
			},
			wantOut: strings.Join(readFileOrPanic("testFiles/out_c_ans.txt"), "\n"),
		},
		{
			name: "out -B first line",
			args: args{
				data:   readFileOrPanic("testFiles/out.txt"),
				params: parametres{pattern: "USD", before: 2},
			},
			wantOut: strings.Join(readFileOrPanic("testFiles/out_b_ans.txt"), "\n"),
		},
		{
			name: "out -A last line",
			args: args{
				data:   readFileOrPanic("testFiles/out.txt"),
				params: parametres{pattern: "bitoc byn", after: 2},
			},
			wantOut: strings.Join(readFileOrPanic("testFiles/out_a_ans.txt"), "\n"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			out := &bytes.Buffer{}
			doGrep(tt.args.data, tt.args.params, out)
			if gotOut := out.String(); gotOut != tt.wantOut {
				t.Errorf("doGrep() = %v, want %v", gotOut, tt.wantOut)
			}
		})
	}
}
