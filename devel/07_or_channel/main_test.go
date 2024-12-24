package main

import (
	"testing"
	"time"
)

func Test_main(t *testing.T) {
	tests := []struct {
		name string
	}{
		{name: "default test"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			start := time.Now()
			main()
			if time.Since(start) > time.Millisecond*1100 {
				t.Fatalf("no passed!")
			}
		})
	}
}
