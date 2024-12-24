package main

import (
	"reflect"
	"strings"
	"testing"
)

func Test_readDataFromFiles(t *testing.T) {
	type args struct {
		filenames []string
	}
	tests := []struct {
		name    string
		args    args
		want    []string
		wantErr bool
	}{
		{
			name: "File exist",
			args: args{filenames: []string{"file.txt"}},
			want: []string{"abhishek 44", "satish 1", "rajan 22", "zvisehn 6", "naveen 2", "divyam 11", "harsh"},
		},
		{
			name:    "File no exist",
			args:    args{filenames: []string{"error.txt"}},
			want:    []string{},
			wantErr: true,
		},
		{
			name:    "File exist and no exist",
			args:    args{filenames: []string{"file.txt", "error.txt"}},
			want:    []string{},
			wantErr: true,
		},
		{
			name: "Some exists files merged",
			args: args{filenames: []string{"file.txt", "numeric.txt"}},
			want: []string{"abhishek 44", "satish 1", "rajan 22", "zvisehn 6", "naveen 2", "divyam 11", "harsh", "10", "1", "13", "2", "22", "X", "Y"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := readDataFromFiles(tt.args.filenames)
			if (err != nil) != tt.wantErr {
				t.Errorf("readDataFromFiles() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("readDataFromFiles() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_makeStringsUnique(t *testing.T) {
	type args struct {
		strs []string
	}
	tests := []struct {
		name string
		args args
		want []string
	}{
		{
			name: "Default",
			args: args{strs: []string{"a", "b", "12", "a", "12", "c"}},
			want: []string{"a", "b", "12", "c"},
		},
		{
			name: "All unique",
			args: args{strs: []string{"a", "b", "c", "d"}},
			want: []string{"a", "b", "c", "d"},
		},
		{
			name: "All same",
			args: args{strs: []string{"a", "a", "a", "a"}},
			want: []string{"a"},
		},
		{
			name: "Empty slice",
			args: args{strs: []string{}},
			want: []string{},
		},
		{
			name: "Nil slice",
			args: args{strs: nil},
			want: []string{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := makeStringsUnique(tt.args.strs); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("makeStringsUnique() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_doSort_numeric(t *testing.T) {
	base := []string{
		"abhishek 44",
		"satish 1",
		"rajan 22",
		"zvisehn 6",
		"naveen 2",
		"divyam 11",
		"harsh",
	}

	var baseDouble []string
	baseDouble = append(baseDouble, base...)
	baseDouble = append(baseDouble, base...)

	// out, err := exec.Command("sort", "-r", "file.txt").Output()
	// fmt.Println("Out:", out, err)
	// k2 := strings.Split(string(out), string(rune(13))+string(rune(10)))
	// out, _ = exec.Command("sort", "file.txt").Output()
	// k1 := strings.Split(string(out), string(rune(13))+string(rune(10)))

	type args struct {
		data   []string
		params parametres
	}
	tests := []struct {
		name string
		args args
		want []string
	}{
		{
			name: "k1",
			args: args{data: base, params: parametres{column: 1, byNumeric: false, isReverse: false, isUnique: false}},
			want: []string{
				"abhishek 44",
				"divyam 11",
				"harsh",
				"naveen 2",
				"rajan 22",
				"satish 1",
				"zvisehn 6",
			},
		},
		{
			name: "k2",
			args: args{data: base, params: parametres{column: 2, byNumeric: false, isReverse: false, isUnique: false}},
			want: []string{
				"harsh",
				"satish 1",
				"divyam 11",
				"naveen 2",
				"rajan 22",
				"abhishek 44",
				"zvisehn 6",
			},
		},
		{
			name: "k3",
			args: args{data: base, params: parametres{column: 3, byNumeric: false, isReverse: false, isUnique: false}},
			want: []string{
				"abhishek 44",
				"divyam 11",
				"harsh",
				"naveen 2",
				"rajan 22",
				"satish 1",
				"zvisehn 6",
			},
		},
		{
			name: "k1 reversed",
			args: args{data: base, params: parametres{column: 1, byNumeric: false, isReverse: true, isUnique: false}},
			want: []string{
				"zvisehn 6",
				"satish 1",
				"rajan 22",
				"naveen 2",
				"harsh",
				"divyam 11",
				"abhishek 44",
			},
		},
		{
			name: "k1 numeric",
			args: args{data: base, params: parametres{column: 1, byNumeric: true, isReverse: false, isUnique: false}},
			want: []string{
				"abhishek 44",
				"divyam 11",
				"harsh",
				"naveen 2",
				"rajan 22",
				"satish 1",
				"zvisehn 6",
			},
		},
		{
			name: "k2 numeric",
			args: args{data: base, params: parametres{column: 2, byNumeric: true, isReverse: false, isUnique: false}},
			want: []string{
				"harsh",
				"satish 1",
				"naveen 2",
				"zvisehn 6",
				"divyam 11",
				"rajan 22",
				"abhishek 44",
			},
		},
		{
			name: "k2 numeric reversed",
			args: args{data: base, params: parametres{column: 2, byNumeric: true, isReverse: true, isUnique: false}},
			want: []string{
				"abhishek 44",
				"rajan 22",
				"divyam 11",
				"zvisehn 6",
				"naveen 2",
				"satish 1",
				"harsh",
			},
		},
		{
			name: "k1 unique",
			args: args{data: baseDouble, params: parametres{column: 1, byNumeric: false, isReverse: false, isUnique: true}},
			want: []string{
				"abhishek 44",
				"divyam 11",
				"harsh",
				"naveen 2",
				"rajan 22",
				"satish 1",
				"zvisehn 6",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := doSort(tt.args.data, tt.args.params); !reflect.DeepEqual([]byte(strings.Join(got, "")), []byte(strings.Join(tt.want, ""))) {
				t.Errorf("doSort() = %v, want %v", got, tt.want)
			}
		})
	}
}
