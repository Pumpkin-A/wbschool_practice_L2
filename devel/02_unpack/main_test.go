package main

import "testing"

func Test_unpakcing(t *testing.T) {
	type args struct {
		str string
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{
			name: "OK - Default",
			args: args{str: "a4bc2d5e"},
			want: "aaaabccddddde",
		},
		{
			name: "Ok - long numbers",
			args: args{str: "a100b50"},
			want: "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaabbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbb",
		},
		{
			name: "Ok - No repeat",
			args: args{str: "abcd"},
			want: "abcd",
		},
		{
			name: "Ok - with unicode",
			args: args{str: "a5한3b2"},
			want: "aaaaa한한한bb",
		},
		{
			name: "Empty string",
			args: args{str: ""},
			want: "",
		},
		{
			name:    "Error string - only numbers",
			args:    args{str: "123"},
			want:    "",
			wantErr: true,
		},
		{
			name:    "Error string - first number",
			args:    args{str: "12ab5d5"},
			want:    "",
			wantErr: true,
		},
		{
			name:    "Error string - non-number or word chars",
			args:    args{str: "a55d5f2.d3"},
			want:    "",
			wantErr: true,
		},
		{
			name: "Escape - 1",
			args: args{str: `qwe\4\5`},
			want: "qwe45",
		},
		// {
		// 	name: "Escape - 2",
		// 	args: args{str: `qwe\45`},
		// 	want: "qwe44444",
		// },
		// {
		// 	name: "Escape - 3",
		// 	args: args{str: `qwe\\5`},
		// 	want: `qwe\\\\\`,
		// },
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := unpakcing(tt.args.str)
			if (err != nil) != tt.wantErr {
				t.Errorf("unpakcing() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("unpakcing() = %v, want %v", got, tt.want)
			}
		})
	}
}
