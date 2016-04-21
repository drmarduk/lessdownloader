package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
	"time"
)

// basic GET function
func httpGET(url string) (string, error) {
	// returns the source of webpage
	resp, err := http.Get(url)
	if err != nil {
		l(err.Error())
		return "", err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		l(err.Error())
		return "", err
	}
	return string(body), nil
}

func httpHEAD(url string) (string, error) {
	cl := &http.Client{}
	req, _ := http.NewRequest("HEAD", url, nil)

	resp, err := cl.Do(req)
	if err != nil {
		l(err.Error())
		return "", err
	}

	return resp.Request.URL.String(), nil
}

func downloadFromUrl(url string, folder string) {
	tokens := strings.Split(url, "/")
	fileName := tokens[len(tokens)-1]

	output, error := os.Create(folder + "/" + fileName)
	if error != nil {
		fmt.Println("Error while creating", fileName, "-", error)
		return
	}
	defer output.Close()

	response, error := http.Get(url)
	if error != nil {
		fmt.Println("Error while downloading", url, "-", error)
		return
	}
	defer response.Body.Close()

	_, error = io.Copy(output, response.Body)
	if error != nil {
		fmt.Println("Error while downloading", url, "-", error)
		return
	}
}

type WriteCounter struct {
	Total int64
}

func (wc *WriteCounter) Write(p []byte) (int, error) {
	n := len(p)
	wc.Total += int64(n)
	fmt.Printf("Read %d bytes for a total of %d\n", n, wc.Total)
	return n, nil
}

func download(url, file string) {
	out, err := os.Create(file)
	if err != nil {
		panic(err)
	}

	defer out.Close()

	resp, err := http.Get(url)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	src := io.TeeReader(resp.Body, &WriteCounter{})
	_, err = io.Copy(out, src)
	if err != nil {
		panic(err)
	}
}

func getFilename(url string) string {
	x := strings.Split(url, "/")
	return x[len(x)-1]
}

func downloadprogress(url, file string) {
	dst, _ := os.Create(file)
	defer dst.Close()
	src, _ := http.Get(url)
	defer src.Body.Close()

	buf := make([]byte, 1000*1024)

	var err error = nil
	var start time.Time
	var end time.Duration
	for err == nil {
		start = time.Now()
		_, err = io.ReadFull(src.Body, buf)
		end = time.Now().Sub(start)
		fmt.Printf("\033[H\033[2J")
		fmt.Printf("Time:\t%.4f/%d\t", end.Seconds(), len(buf))
		fmt.Printf("Speed:\t%.4f\n", float64(1)/float64(end.Seconds())*float64(len(buf)))

		// 146ms -> 100k
		// 1000 / 146 * len(buf)
		// 1000ms-> xxxK
		dst.Write(buf)
	}
}

func writeMetaData(url string, folder string) {
	d1 := []byte(url)
	ioutil.WriteFile(folder+"/metadata.txt", d1, 0644)
}

// "main" logger, maybe to file, default to stdout
func l(message string) {
	if len(message) > 1 {
		log.Printf("%s %s\n", timestamp(), message)
	}
}

// simple timestamp, no year
func timestamp() string {
	t := time.Now()
	layout := "[31.12 - 24:59:59]"
	return t.Format(layout)
}
