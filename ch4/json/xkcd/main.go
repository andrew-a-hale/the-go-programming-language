package main

import (
	"database/sql"
	"encoding/json"
	"encoding/xml"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"

	_ "github.com/mattn/go-sqlite3"
)

type Feed struct {
	Title   string
	Link    string
	Id      string
	Updated string
	Entries []Entry `xml:"entry"`
}

type Entry struct {
	Title   string
	Link    string
	Updated string
	Id      string `xml:"id"`
	Summary string
}

type ComicMetadata struct {
	Num        int    `json:"num"`
	Year       string `json:"year"`
	Month      string `json:"month"`
	Day        string `json:"day"`
	News       string `json:"news"`
	Title      string `json:"title"`
	SafeTitle  string `json:"safe_title"`
	Transcript string `json:"transcript"`
	ImgUri     string `json:"img"`
	Alt        string `json:"alt"`
}

func main() {
	db, err := sql.Open("sqlite3", "./xkcd.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	switch os.Args[1] {
	case "seed":
		seed(db)
	case "search":
		term := os.Args[2]
		if term == "" {
			log.Fatal("cmd `search` requires a search term")
		}

		search(db, term)
	default:
		log.Fatal("expected seend or search command")
	}
}

func seed(db *sql.DB) {
	sqlStmt := `
	CREATE TABLE IF NOT EXISTS comics (
		num INTEGER,
		year INTEGER,
		month INTEGER,
		day INTEGER,
		news TEXT,
		title TEXT,
		safe_title TEXT,
		transcript TEXT,
		img_uri TEXT,
		alt TEXT
	);`

	_, err := db.Exec(sqlStmt)
	if err != nil {
		log.Fatal(err)
	}

	latest_id, err := read_xkcd_feed()
	if err != nil {
		log.Fatal(err)
	}

	for id := 1; id <= latest_id; id++ {
		comic, err := download_comic_metadata(db, id)
		if err != nil {
			continue
		}

		db, err := sql.Open("sqlite3", "./xkcd.db")
		if err != nil {
			log.Fatal(err)
		}

		write_to_db(db, comic)
	}

}

func read_xkcd_feed() (int, error) {
	feed_url := "https://xkcd.com/atom.xml"
	resp, err := http.Get(feed_url)
	if err != nil {
		return 0, errors.New("failed to get " + feed_url)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return 0, errors.New("failed to read " + feed_url)
	}

	feed := &Feed{}
	if err := xml.Unmarshal(body, feed); err != nil {
		log.Fatal(err)
	}

	latest := feed.Entries[0].Id
	latest_id, err := strconv.ParseInt(strings.Split(latest, "/")[3], 10, 0)
	if err != nil {
		return 0, errors.New("failed to get latest id from feed")
	}

	return int(latest_id), nil
}

func download_comic_metadata(db *sql.DB, id int) (*ComicMetadata, error) {
	if check_comic_already_in_db(db, id) {
		fmt.Printf("xkcd comic %d already indexed\n", id)
		return &ComicMetadata{}, nil
	}

	comic_url := fmt.Sprintf("https://xkcd.com/%d/info.0.json", id)
	resp, err := http.Get(comic_url)
	if err != nil {
		return &ComicMetadata{}, errors.New("failed to get " + comic_url)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return &ComicMetadata{}, errors.New("failed to read " + comic_url)
	}

	comic := &ComicMetadata{}
	if err := json.Unmarshal(body, comic); err != nil {
		return &ComicMetadata{}, errors.New("failed to read " + comic_url)
	}

	return comic, nil
}

func check_comic_already_in_db(db *sql.DB, num int) bool {
	stmt, err := db.Prepare(`select num from comics where num = ?`)
	if err != nil {
		log.Fatal(err)
	}

	var id int
	stmt.QueryRow(num).Scan(&id)
	return id != 0
}

func write_to_db(db *sql.DB, comic *ComicMetadata) error {
	tx, err := db.Begin()
	if err != nil {
		log.Fatal(err)
	}

	stmt, err := tx.Prepare(`
		insert into comics(num, year, month, day, news, title, safe_title, transcript, img_uri, alt)
		values(?, ?, ?, ?, ?, ?, ?, ?, ?, ?);
	`)
	if err != nil {
		log.Fatal(err)
	}
	defer stmt.Close()

	_, err = stmt.Exec(comic.Num, comic.Year, comic.Month, comic.Day, comic.News, comic.Title, comic.SafeTitle, comic.Transcript, comic.ImgUri, comic.Alt)
	if err != nil {
		log.Fatal(err)
	}

	err = tx.Commit()
	if err != nil {
		log.Fatal(err)
	}

	if comic.Num > 0 {
		fmt.Printf("index xkcd comic %d\n", comic.Num)
	} else {
	}
	return nil
}

func search(db *sql.DB, term string) {
	stmt := fmt.Sprintf(`select num, transcript from comics where transcript like '%%%s%%'`, term)
	rows, err := db.Query(stmt)
	if err != nil {
		log.Fatal(err)
	}

	i := 0
	for rows.Next() {
		var num int
		var transcript string
		err = rows.Scan(&num, &transcript)
		if err != nil {
			log.Fatal(err)
		}

		url := fmt.Sprintf("https://xkcd.com/%d", num)

		fmt.Printf("%s\n%s\n\n", url, transcript)
		i++
	}

	fmt.Printf("comics found: %d\n", i)
}
