package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
	"sync"

	"github.com/PuerkitoBio/goquery"
)

const url string = "https://www.countryflags.com/"

func getImageUrlList(url string) ([]string, int) {
	var imageUrls []string
	if doc, err := goquery.NewDocument(url); err != nil {
		log.Fatal(err)
	} else {
		doc.Find(".thumb img").Each(func(i int, s *goquery.Selection) {
			imageUrls = append(imageUrls, s.AttrOr("src", ""))
		})
	}
	return imageUrls, len(imageUrls)
}

func downloadImage(url string, fileName string, waitGroup *sync.WaitGroup) {
	if resp, err := http.Get(url); err != nil {
		log.Fatal(err)
	} else {
		defer resp.Body.Close()
		if file, err := os.Create("Flags/" + fileName + ".png"); err != nil {
			log.Fatal(err)
		} else {
			defer file.Close()
			if _, err := io.Copy(file, resp.Body); err != nil {
				log.Fatal(err)
			}
		}
	}
	waitGroup.Done()
}

func main() {
	imageUrlList, imageUrlListLen := getImageUrlList(url)

	if _, err := os.Stat("Flags"); os.IsNotExist(err) {
		os.Mkdir("Flags", 0o755)
	}

	waitGroup := new(sync.WaitGroup)

	// Make concurrent requests to download images
	for i, image := range imageUrlList {
		fileName := strings.Split(image, "/")[4]
		fmt.Println("Downloading "+fileName, "...", i+1, "of", imageUrlListLen)

		waitGroup.Add(1)
		go downloadImage(image, fileName, waitGroup)
	}

	waitGroup.Wait()
}
