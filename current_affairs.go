package main

import (
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
		fmt.Println("Error while downloading:", url)
	}

	temp := strings.Split(url, "/")
	fileName := temp[len(temp)-1]
	fmt.Println(fileName)
	defer response.Body.Close()
	dir := os.Getenv("TNPSC_DIR")
	file, err := os.Create(dir + fileName)
	if err != nil {
		log.Fatal("Error while creating file ", fileName, " ", err)
	}
	_, err = io.Copy(file, response.Body)

	if err != nil {
		log.Fatal("Error while downloading file ", fileName, " ", err)
	}

	file.Close()
	fmt.Println("File downloaded: ==>", fileName)
}

func main() {
	url := "http://www.tnpscportal.in/2014/06/tnpsc-current-affairs-in-tamil-june-2014.html"

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

			ok, url := getHref(token)
			if !ok {
				continue
			}

			if strings.Index(url, "pdf") > -1 {
				go savePdf(url)
			}

		}
	}
}
