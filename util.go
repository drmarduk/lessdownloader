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
