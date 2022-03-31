package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

func main() {
	url := "https://www.countryflags.com/"
	resp, err := http.Get(url)
	if err != nil {
		fmt.Println(err)
	}
	defer resp.Body.Close()

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	imageLinks := make([]string, 0)

	doc.Find(".thumb img").Each(func(i int, s *goquery.Selection) {
		imageLinks = append(imageLinks, s.AttrOr("src", ""))
	})

	// Create Flags folder if not exists
	if _, err := os.Stat("Flags"); os.IsNotExist(err) {
		os.Mkdir("Flags", 0o755)
	}

	imageLinksLength := len(imageLinks)

	// Save the images to Flag Folder
	for i, link := range imageLinks {
		resp, err := http.Get(link)
		if err != nil {
			fmt.Println(err)
		}
		defer resp.Body.Close()

		countryName := strings.Split(link, "/")[4]

		file, err := os.Create("Flags/" + countryName + ".png")
		if err != nil {
			log.Fatal(err)
		}
		defer file.Close()

		_, err = io.Copy(file, resp.Body)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println("Downloaded: "+countryName, i, "/", imageLinksLength)
	}
}
