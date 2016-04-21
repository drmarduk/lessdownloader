package main

import "testing"

func TestMotherlessGalleryOpen(t *testing.T) {
	tests := []struct {
		in                               string
		images, videos, galleries, posts int
	}{
		{"http://motherless.com/G43457A1", 67, 13, 0, 0},
		{"http://motherless.com/G4D50E8B", 19, 0, 0, 0},
		{"http://motherless.com/GD61E4D4", 192, 0, 4, 0},
		{"http://motherless.com/GD1CD49F", 26, 0, 0, 0},
		{"http://motherless.com/GC2857B6", 2, 0, 4, 0},
		{"http://motherless.com/G8BA035B", 32, 0, 0, 0},
		{"http://motherless.com/G35631A7", 44, 0, 0, 0},
		{"http://motherless.com/G6E3028B", 19, 0, 0, 0},    // nice :)
		{"http://motherless.com/G016CE01", 6, 68, 116, 0},  // mal nicht familiy
		{"http://motherless.com/GAA7F7DC", 1173, 51, 0, 0}, // over 1000
		{"http://motherless.com/gi/geile_sammlung___hot_pics___vids", 3103, 471, 0, 0},
	}

	for _, tt := range tests {
		g := NewMotherlessGallery(tt.in)
		if err := g.Open(); err != nil {
			t.Fatalf("Open(%s): nope: %s\n", tt.in, err.Error())

		}
		if tt.images != g.ImageCount {
			t.Fatalf("ImageCount: got %d, expected: %d\n", g.ImageCount, tt.images)
		}
		if tt.videos != g.VideoCount {
			t.Error("video count not good")
		}
		//if tt.galleries != g.GalleryCount {
		//	t.Error("gallery count not good")
		//}
		//if tt.posts != g.PostCount {
		//	t.Error("post count not good")
		//}
	}
}

func TestGalleryImages(t *testing.T) {
	tests := []struct {
		in  string
		c   int
		out []string
	}{
		{"http://motherless.com/G4D50E8B", 19, []string{""}},
	}

	for _, tt := range tests {
		g := NewMotherlessGallery(tt.in)
		images, err := g.GetImages()
		if err != nil {
			t.Fatalf("Error while getting images: " + err.Error())
		}
		if len(images) != tt.c {
			t.Errorf("GetImages(%s): got: %d, expected: %d", tt.in, len(images), tt.c)
		}
	}
}
