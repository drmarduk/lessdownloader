package main

import (
	"flag"
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"
	_ "time"
)

func main() {
	var arg *string = flag.String("url", "", "url to parse")
	var typ *string = flag.String("typ", "img", "what typ to process")
	//var action *string = flag.String("action", "download", "the action to perform")
	flag.Parse()
	src, _ := httpGET(*arg)
	url := ""
	switch *typ {
	case "img":
		url = processImage(src)
		downloadprogress(url, "data/"+getFilename(url))
	case "vid":
		url = processVideo(src)
		downloadprogress(url, "data/"+getFilename(url))
	case "gal":
		img := NewMotherlessGallery(url)
		img.Open()
		images, err := img.GetImages()
		if err != nil {
			panic(err)
		}
		for _, i := range images {
			downloadprogress(i, "data/"+getFilename(i))
		}
		break
	}
}

func processGallery(src string) string {
	// return Gallery ID
	paths := strings.Split(src, "/")
	return paths[len(paths)-1]
}

func processVideo(src string) string {
	reg := regexp.MustCompile(`__fileurl = 'http://cdn.videos.motherlessmedia.com/videos/[a-zA-Z0-9]{3,9}.(mp4|webm|avi|mpg|mpeg|gif|gifv)`)
	match := reg.FindString(src)
	match = strings.Replace(match, `__fileurl = '`, "", 1)
	return match
}

func processImage(src string) string {
	reg := regexp.MustCompile(`__fileurl = 'http://cdn.images.motherlessmedia.com/images/[a-zA-Z0-9]{3,9}.(jpg|jpeg|png|gif)`)
	match := reg.FindString(src)
	match = strings.Replace(match, `__fileurl = '`, "", 1)
	return match
}

// alte main, ugly, might be usefull
func dmain() {
	// ID = len(7)
	// Gallery main = G + ID -> len(8)
	// Gallery Images = G + I + ID -> len(9)

	// extract ID
	// http://motherless.com/G0471F1B

	var url string = os.Args[1]

	fmt.Println("Input URL:\t" + url) // debug
	if strings.Contains(url, "http://motherless.com/") {
		var validID bool = false
		var src, id, title string
		var maximg int

		url = url[22:]                // id, die uebergeben wurde
		fmt.Println("ID2:\t\t" + url) // debug

		if string(url[0]) == "G" {
			fmt.Println("[+] Gallery found.")

			// Main Gallery Site or Image, Video
			var scase string = string(url[1])
			switch {
			case scase == "M":
				fmt.Println("[+] Gallery - All Uploads")
				id = url[2:]
				validID = false

			case scase == "I":
				fmt.Println("[+] Gallery - Image Uploads")
				id = url[2:]
				validID = true

			case scase == "V":
				fmt.Println("[+] Gallery - Video Uploads")
				id = url[2:]
				validID = false

			case scase == "G":
				fmt.Println("[+] Gallery - Galleries")
				id = url[2:]
				validID = false

			case scase == "P":
				fmt.Println("[+] Gallery - Board Posts")
				id = url[2:]
				validID = false

			case len(url[1:]) == 7:
				fmt.Println("[+] Gallery - Home")
				id = url[1:]
				validID = true
			}
		} else {
			fmt.Println("[-] No Gallery found, might be a single Picture or Video.")
			//id = url
			validID = false // do special.
		}

		// get HTML Code
		src, _ = httpGET("http://motherless.com/G" + id) // muesste fuer bilder auf GI statt G

		// start downloading MetaData
		if validID {
			// get Gallery Name
			reg := regexp.MustCompile(`<title>(.*) - MOTHERLESS.COM</title>`)
			match := reg.FindString(src)
			title = strings.Replace(match, "<title>", "", -1)
			title = strings.Replace(title, " - MOTHERLESS.COM</title>", "", -1)
			// create folder
			os.Mkdir(title, 0777)
			fmt.Println("[!] Gallery: " + title)

			// Image count
			// &nbsp; Images  [ 55 ] &nbsp;
			// &nbsp; Images  [ 4,800 ] &nbsp;
			reg = regexp.MustCompile(`(Images|Videos)  \[ (?P<count>.*) \]`)
			matches := reg.FindAllString /*Submatch*/ (src, -1)
			for _, s := range matches {
				// Images  [ 15 ]
				var typ string = s[:strings.Index(s, " ")]
				var count string = strings.Replace(s[strings.Index(s, "[ ")+2:], " ]", "", 1)

				if strings.Contains(typ, "mages") {
					maximg, _ = strconv.Atoi(count)
				}

				fmt.Println("[!] " + typ + ": " + count)
			}

			// ab hier laedt es nur die erste Seite, das is muell
			// loop for pages and extract images/videos
			// regex:  <a href="/G+id/[0-9A-F]{7}"
			fmt.Println("[!] Search for Images in Gallery.")

			// page 1 - extract images
			// current Page
			fmt.Println("Downlad: max. " + strconv.Itoa(maximg) + " Images.")
			var imgagecount int = 1
			for i := 1; !strings.Contains(src, "There are no images in this gallery."); { // get max
				fmt.Println("Seite: " + strconv.Itoa(i))
				// fuer images immer GI
				src, _ = httpGET("http://motherless.com/GI" + id + "?page=" + strconv.Itoa(i)) // muesste fuer bilder auf GI statt G

				reg = regexp.MustCompile(`<a href="/G` + id + `/[0-9A-F]*"`)
				matches := reg.FindAllString(src, -1)
				for _, s := range matches {
					s = strings.Replace(s, "<a href=\"", "", -1)
					s = strings.Replace(s, "\"", "", -1)
					iURL := "http://motherless.com/" + s

					// process iURL
					srcI, _ := httpGET(iURL)

					// __fileurl = 'http://s16.motherlessmedia.com/dev301/0/549/152/0549152017.jpg';
					regI := regexp.MustCompile(`__fileurl = '(.*)'`)
					match := regI.FindString(srcI)
					match = strings.Replace(match, "__fileurl = '", "", -1)
					match = strings.Replace(match, "'", "", -1)

					downloadFromUrl(match, "./"+title)

					fmt.Println("[!] [" + strconv.Itoa(imgagecount) + " | " + strconv.Itoa(maximg) + "]")
					imgagecount++
				}
				i++
			}

			go writeMetaData("http://motherless.com/G"+id, title)

		}

	}

}
