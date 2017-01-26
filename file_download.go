package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"

	"golang.org/x/net/html"
)

func getHref(t html.Token) (ok bool, href string) {
	for _, a := range t.Attr {
		if a.Key == "href" {
			href = a.Val
			ok = true
		}
	}
	return
}

func savePdf(url string) {
	response, e := http.Get(url)
	if e != nil {
		fmt.Println("Error while downloading ==>", url)
	}
	defer response.Body.Close()

	temp := strings.Split(url, "/")
	fileName := temp[len(temp)-1]

	dir := os.Getenv("TNPSC_DIR")
	file, err := os.Create(dir + fileName)
	if err != nil {
		log.Fatal("Error while creating file ==>", fileName, " ", err)
	}

	_, err = io.Copy(file, response.Body)
	if err != nil {
		log.Fatal("Error while downloading file ==>", fileName, " ", err)
	}

	file.Close()

	fmt.Println("File downloaded ==>", fileName)
}

func main() {
	var url, fileType string
	flag.StringVar(&url, "url", "url", "URL to crawl")
	flag.StringVar(&fileType, "ftype", "pdf", "File type to download")
	flag.Parse()

	//remaining := make(chan int)
	//done := make(chan bool, false)

	resp, err := http.Get(url)
	if err != nil {
		fmt.Println("ERROR: ", err)
		return
	}

	b := resp.Body
	defer b.Close()

	tokens := html.NewTokenizer(b)
	//quit := false

	for {
		token := tokens.Next()
		switch {
		case token == html.StartTagToken:
			token := tokens.Token()

			if token.Data != "a" {
				continue
			}

			ok, downloadURL := getHref(token)
			if !ok {
				continue
			}

			if !strings.Contains(downloadURL, "www") {
				strippedURL := strings.Replace(downloadURL, "http://", "", -1)
				baseURL := strings.Split(strippedURL, "/")[0]
				downloadURL = baseURL + downloadURL
				fmt.Println("downloadURL ==>", downloadURL)
			}

			if strings.Index(downloadURL, fileType) > -1 {
				go savePdf(downloadURL)
			}
			/*case token == html.ErrorToken:
			if done {
				return
			}*/
		}
	}
}
