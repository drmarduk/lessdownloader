package main

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

type MotherlessGallery struct {
	Url string

	ImageCount   int
	VideoCount   int
	GalleryCount int
	PostCount    int

	src string
}

func NewMotherlessGallery(url string) *MotherlessGallery {
	return &MotherlessGallery{Url: url}
}

func (g *MotherlessGallery) Open() error {
	var err error
	g.src, err = httpGET(g.Url)
	if err != nil {
		return err
	}

	imgRegex := regexp.MustCompile(`Images [( |)[0-9,]{1,10}( |)]`)
	x := imgRegex.FindString(g.src)
	x = strings.Replace(x, " ", "", -1)
	x = strings.Replace(x, "Images[", "", 1)
	x = strings.Replace(x, "]", "", 1)
	x = strings.Replace(x, ",", "", -1)
	y, err := strconv.Atoi(x)
	if err != nil {
		panic(err)
	}
	g.ImageCount = y

	vidRegex := regexp.MustCompile(`Videos [( |)[0-9,]{1,10}( |)]`)
	x = vidRegex.FindString(g.src)
	x = strings.Replace(x, " ", "", -1)
	x = strings.Replace(x, "Videos[", "", 1)
	x = strings.Replace(x, "]", "", 1)
	x = strings.Replace(x, ",", "", -1)
	y, _ = strconv.Atoi(x)
	g.VideoCount = y

	galRegex := regexp.MustCompile(`Galleries [( |)[0-9,]{1,10}( |)]`)
	x = galRegex.FindString(g.src)
	x = strings.Replace(x, " ", "", -1)
	x = strings.Replace(x, "Galleries[", "", 1)
	x = strings.Replace(x, "]", "", 1)
	x = strings.Replace(x, ",", "", -1)
	y, _ = strconv.Atoi(x)
	g.GalleryCount = y

	postRegex := regexp.MustCompile(`Posts [( |)[0-9,]{1,10}( |)]`)
	x = postRegex.FindString(g.src)
	x = strings.Replace(x, " ", "", -1)
	x = strings.Replace(x, "Posts[", "", 1)
	x = strings.Replace(x, "]", "", 1)
	x = strings.Replace(x, ",", "", -1)
	y, _ = strconv.Atoi(x)
	g.PostCount = y
	return nil
}

// Images [15,392]
// Videos [4,831]
// Images  [ 43 ]
// Videos  [ 103 ]

func (g *MotherlessGallery) GetImages() ([]string, error) {
	var result []string

	var pageCount int = g.ImageCount / 80 // ca. total pages
	var i int = 1

	for i = 1; i <= pageCount; i++ {
		src, err := httpGET(g.Url + fmt.Sprintf("?page=%d", i))
		if err != nil {
			panic(err)
		}

		// got source for every page, extract image urls
		result = append(result, extractImageUrl(g.Url, src)...)
	}

	fmt.Println(result)
	return result, nil
}

func extractImageUrl(prefix, src string) []string {
	re := regexp.MustCompile(`href="/[a-zA-Z0-9]{3,13}" class="img-container"`)
	match := re.FindAllString(src, -1)

	var result []string

	for _, tmp := range match {
		tmp := strings.Replace(tmp, `href="`, "", 1)
		tmp = strings.Replace(tmp, ` class="img-container"`, "", 1)
		result = append(result, prefix+tmp)
	}

	return result
}
