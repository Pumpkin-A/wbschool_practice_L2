package main

import "testing"

func Test_unpacking(t *testing.T) {
	type args struct {
		str string
	}
	tests := []struct {
		name     string
		unpacker *Unpacker
		args     args
		want     string
		wantErr  bool
	}{
		{
			name:     "OK - Default",
			unpacker: NewUnpacker(),
			args:     args{str: "a4bc2d5e"},
			want:     "aaaabccddddde",
		},
		{
			name:     "Ok - long numbers",
			unpacker: NewUnpacker(),
			args:     args{str: "a100b50"},
			want:     "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaabbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbb",
		},
		{
			name:     "Ok - No repeat",
			unpacker: NewUnpacker(),
			args:     args{str: "abcd"},
			want:     "abcd",
		},
		{
			name:     "Ok - with unicode",
			unpacker: NewUnpacker(),
			args:     args{str: "a5한3b2"},
			want:     "aaaaa한한한bb",
		},
		{
			name:     "Empty string",
			unpacker: NewUnpacker(),
			args:     args{str: ""},
			want:     "",
		},
		{
			name:     "Error string - only numbers",
			unpacker: NewUnpacker(),
			args:     args{str: "123"},
			want:     "",
			wantErr:  true,
		},
		{
			name:     "Error string - first number",
			unpacker: NewUnpacker(),
			args:     args{str: "12ab5d5"},
			want:     "",
			wantErr:  true,
		},
		{
			name:     "non-number or word chars",
			unpacker: NewUnpacker(),
			args:     args{str: "a5d5f2.3d"},
			want:     "aaaaadddddff...d",
		},
		{
			name:     "Escape - 1",
			unpacker: NewUnpacker(),
			args:     args{str: `qwe\4\5`},
			want:     "qwe45",
		},
		{
			name:     "Escape - 2",
			unpacker: NewUnpacker(),
			args:     args{str: `qwe\45`},
			want:     "qwe44444",
		},
		{
			name:     "Escape - 3",
			unpacker: NewUnpacker(),
			args:     args{str: `qwe\\5`},
			want:     `qwe\\\\\`,
		},
		{
			name:     "Escape - 4",
			unpacker: NewUnpacker(),
			args:     args{str: `qwe\`},
			want:     "",
			wantErr:  true,
		},
		{
			name:     "Escape - 5",
			unpacker: NewUnpacker(),
			args:     args{str: `\\\`},
			want:     "",
			wantErr:  true,
		},
		{
			name:     "Escape - 6",
			unpacker: NewUnpacker(),
			args:     args{str: `\\\\a`},
			want:     `\\a`,
		},
		{
			name:     "Escape - 7",
			unpacker: NewUnpacker(),
			args:     args{str: `\f`},
			want:     "",
			wantErr:  true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.unpacker.Do(tt.args.str)
			if (err != nil) != tt.wantErr {
				t.Errorf("unpacking() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("unpacking() = %v, want %v", got, tt.want)
			}
		})
	}
}
