package main

import (
	"reflect"
	"testing"
)

func Test_newAnnoDict(t *testing.T) {
	type args struct {
		words []string
	}
	tests := []struct {
		name string
		args args
		want map[string][]string
	}{
		{
			name: "values sorting",
			args: args{words: []string{"тяпка", "пятка", "пятак"}},
			want: map[string][]string{"акптя": {"пятак", "пятка", "тяпка"}},
		},
		{
			name: "Diff letter cases",
			args: args{words: []string{"тЯпкА", "ПЯТКА", "пЯтАК"}},
			want: map[string][]string{"акптя": {"пятак", "пятка", "тяпка"}},
		},
		{
			name: "With repeats",
			args: args{words: []string{"тяпка", "пятка", "пятак", "тяпка", "пятка", "пятак"}},
			want: map[string][]string{"акптя": {"пятак", "пятка", "тяпка"}},
		},
		{
			name: "Remove one line",
			args: args{words: []string{"тяпка", "пятка", "пятак", "афанасий", "никитиН", "НИКИТИН"}},
			want: map[string][]string{"акптя": {"пятак", "пятка", "тяпка"}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := newAnnoDict(tt.args.words); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("newAnnoDict() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_convertStringToKey(t *testing.T) {
	type args struct {
		word string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			args: args{word: "пятак"},
			want: "акптя",
		},
		{
			args: args{word: "пятка"},
			want: "акптя",
		},
		{
			args: args{word: "тяпка!"},
			want: "!акптя",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := convertStringToKey(tt.args.word); got != tt.want {
				t.Errorf("convertStringToKey() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_sortMapValues(t *testing.T) {
	type args struct {
		dict map[string][]string
	}
	tests := []struct {
		name string
		args args
		want map[string][]string
	}{
		{
			name: "Rus text default sort",
			args: args{map[string][]string{"абвя": {"бвфя", "ябая", "авбя"}}},
			want: map[string][]string{"абвя": {"авбя", "бвфя", "ябая"}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := sortMapValues(tt.args.dict); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("sortMapValues() = %v, want %v", got, tt.want)
			}
		})
	}
}
