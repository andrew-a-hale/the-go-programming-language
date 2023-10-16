package main

import (
	"encoding/json"
	"fmt"
	"image"
	"image/jpeg"
	"io"
	"log"
	"net/http"
	"os"
)

type Movie struct {
	Title  string `json:"title"`
	Poster string `json:"poster"`
}

func main() {
	apikey_file, err := os.Open("apikey")
	if err != nil {
		log.Fatal(err)
	}

	apikey, err := io.ReadAll(apikey_file)
	if err != nil {
		log.Fatal(err)
	}

	movie_title := os.Args[1]
	url := fmt.Sprintf("http://www.omdbapi.com/?apikey=%s&t=%s", apikey, movie_title)
	resp, err := http.Get(url)
	if err != nil {
		log.Fatal(err)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	movie := &Movie{}
	if err := json.Unmarshal(body, movie); err != nil {
		log.Fatal(err)
	}

	resp, err = http.Get(movie.Poster)
	if err != nil {
		log.Fatal(err)
	}

	poster, _, err := image.Decode(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	out, err := os.Create(movie_title+".jpg")
	if err != nil {
		log.Fatal(err)
	}
	defer out.Close()

	err = jpeg.Encode(out, poster, &jpeg.Options{Quality: 100})
	if err != nil {
		log.Fatal(err)
	}
}
