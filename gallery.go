package main

import (
	"fmt"
	"log"
	"regexp"
	"strconv"
	"strings"
)

type MotherlessGallery struct {
	Url string
	Id  string // die ID von der Gallery quasi

	ImageCount   int
	VideoCount   int
	GalleryCount int
	PostCount    int

	src string
}

func NewMotherlessGallery(url string) *MotherlessGallery {
	// da muss ich dem user (dir!!!) vertrauen
	x := strings.Split(url, "/")
	id := x[len(x)-1]
	return &MotherlessGallery{Url: url, Id: id}
}

// rm is a simple replace for slices
func rm(s string, list ...string) string {
	for _, x := range list {
		s = strings.Replace(s, x, "", -1)
	}
	return s
}

func (g *MotherlessGallery) Open() error {
	var err error
	g.src, err = httpGET(g.Url)
	if err != nil {
		return err
	}

	imgRegex := regexp.MustCompile(`Images [( |)[0-9,]{1,10}( |)]`)
	x := imgRegex.FindString(g.src)
	x = rm(x, " ", "Images[", "]", ",")
	y, err := strconv.Atoi(x)
	if err != nil {
		panic(err)
	}
	g.ImageCount = y

	vidRegex := regexp.MustCompile(`Videos [( |)[0-9,]{1,10}( |)]`)
	x = vidRegex.FindString(g.src)
	x = rm(x, " ", "Videos[", "]", ",")
	y, _ = strconv.Atoi(x)
	g.VideoCount = y

	galRegex := regexp.MustCompile(`Galleries [( |)[0-9,]{1,10}( |)]`)
	x = galRegex.FindString(g.src)
	x = rm(x, " ", "Galleries[", "]", ",")
	y, _ = strconv.Atoi(x)
	g.GalleryCount = y

	postRegex := regexp.MustCompile(`Posts [( |)[0-9,]{1,10}( |)]`)
	x = postRegex.FindString(g.src)
	x = rm(x, " ", "Posts[", "]", ",")
	y, _ = strconv.Atoi(x)
	g.PostCount = y

	return nil
}

func (g *MotherlessGallery) GetImages() ([]string, error) {
	var result []string
	var i int = 1

	for {
		//for i = 1; i <= pageCount; i++ {
		url := g.Url + fmt.Sprintf("?page=%d", i)
		src, err := httpGET(url)
		if err != nil {
			panic(err) // TODO: change
		}

		// got source for every page, extract image urls
		result = append(result, extractImagess(g.Url, src)...)
		// ende der liste
		if strings.Contains(src, `<span class="current" rel="0"> NEXT &raquo;</span>`) {
			break
		}
		i += 1
	}
	return result, nil
}

func extractImagess(prefix, src string) []string {
	// http://cdn.thumbs.motherlessmedia.com/thumbs/3533EFD-zoom.jpg?from_helper
	re := regexp.MustCompile(`data-codename="[a-zA-Z0-9]{3,13}" `)
	matches := re.FindAllString(src, -1)

	var result []string

	for _, tmp := range matches {
		tmp = strings.Replace(tmp, `data-codename="`, "", 1)
		tmp = strings.Replace(tmp, `" `, "", 1)
		result = append(result,
			"http://cdn.images.motherlessmedia.com/images/"+tmp+".jpg") // TODO: might be png, or else
	}
	log.Printf("Url: %s IoP: %d\n", prefix, len(result))

	return result
}

func extractImageUrl(prefix, src string) []string {
	re := regexp.MustCompile(`href="(/[a-zA-Z0-9]{3,13}|)/[a-zA-Z0-9]{3,13}" class="img-container"`)
	match := re.FindAllString(src, -1)

	var result []string

	for _, tmp := range match {
		tmp := strings.Replace(tmp, `href="`, "", 1)
		tmp = strings.Replace(tmp, ` class="img-container"`, "", 1)
		result = append(result, prefix+tmp)
	}
	return result
}
