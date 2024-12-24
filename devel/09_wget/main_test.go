package main

import (
	"os"
	"testing"
)

func Test_counter_download(t *testing.T) {
	type args struct {
		uri      string
		maxLevel int
	}
	tests := []struct {
		name string
		c    *counter
		args args
	}{
		{
			c: newDownloadCounter(),
			args: args{
				uri:      "https://sun9-7.userapi.com/impg/P5qwKdi3VEJvVXv6KvTrTs4ij9s9hHeyAMWr6w/d0lgFdt9UOY.jpg?size=806x1080&quality=96&sign=8d98653ac22de3ebf2b2f7864776df88&type=album",
				maxLevel: 0,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.c.download(tt.args.uri, tt.args.maxLevel)
			expectedFilename := "d0lgFdt9UOY.jpg"
			info1, _ := os.Stat(expectedFilename)
			info2, _ := os.Stat("test.jpg")
			if info1.Size() != info2.Size() {
				t.Error("not equal!!")
			}
			os.Remove(expectedFilename)
		})
	}
}
