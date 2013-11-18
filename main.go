package main

import (
	"fmt"
	"os"
	_ "time"
	"regexp"
	"strings"
)

func main() {
	// ID = len(7)
	// Gallery main = G + ID -> len(8)
	// Gallery Images = G + I + ID -> len(9)

	// extract ID

	var url string = os.Args[1]

	fmt.Println("Input URL:\t" + url) // debug
	if strings.Contains(url, "http://motherless.com/") {
		var validID bool = false
		var src, id, title string

		url = url[22:]
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
		src, _ = httpGET("http://motherless.com/G" + id)
		
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
			matches := reg.FindAllString/*Submatch*/(src, -1)
			for _, s := range matches {
				// Images  [ 15 ]
				var typ string = s[:strings.Index(s, " ")]
				var count string = strings.Replace(s[strings.Index(s, "[ ") + 2:], " ]", "", 1)
				fmt.Println("[!] " + typ + ": " + count)
			}


			// loop for pages and extract images/videos
			// regex:  <a href="/G+id/[0-9A-F]{7}"
			fmt.Println("[!] Search for Images in Gallery.")
			reg = regexp.MustCompile(`<a href="/G` + id + `/[0-9A-F]*"`)
			matches = reg.FindAllString(src, -1)
			for _, s := range matches {
				s = strings.Replace(s, "<a href=\"", "", -1)
				s = strings.Replace(s, "\"", "", -1)
				iURL := "http://motherless.com/" + s
				//fmt.Println(iURL)

				// process iURL
				srcI, _ := httpGET(iURL)

				// __fileurl = 'http://s16.motherlessmedia.com/dev301/0/549/152/0549152017.jpg';
				regI := regexp.MustCompile(`__fileurl = '(.*)'`)
				match := regI.FindString(srcI)
				match = strings.Replace(match, "__fileurl = '", "", -1)
				match = strings.Replace(match, "'", "", -1)
				
				//fmt.Println("\t- " + match)

				//downloadFromUrl(match, "./" + title)
			}


			// page 1 - extract images
			// current Page
		}

	}

}