package main

import "testing"

func TestProcessVideo(t *testing.T) {
	tests := []struct {
		in, out string
	}{
		{"http://motherless.com/G1BAABCE/B4879D4", "http://cdn.videos.motherlessmedia.com/videos/B4879D4.mp4"},
	}

	for _, tt := range tests {
		src, err := httpGET(tt.in)
		if err != nil {
			panic(err)
		}
		got := processVideo(src)
		if got != tt.out {
			t.Fatalf("processVideo(%s): got: %s, expected: %s\n", tt.in, got, tt.out)
		}
	}
}
