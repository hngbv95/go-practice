package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"strings"
)

type Option struct {
	Text string `json:"text"`
	Arc  string `json:"arc"`
}

type Arc struct {
	Title   string   `json:"title"`
	Story   []string `json:"story"`
	Options []Option
}

type StoryLines map[string]Arc

func parseJsonFile(jsonData []byte) (StoryLines, error) {
	result := StoryLines{}
	//fmt.Println(string(jsonData))
	err := json.Unmarshal(jsonData, &result)
	return result, err
}

func readFile(name string) ([]byte, error) {
	file, err := os.Open("adventure.json")

	if err != nil {
		return nil, err
	}
	defer file.Close()
	stat, err := file.Stat()
	if err != nil {
		return nil, err
	}

	fileContent := make([]byte, stat.Size())
	_, errRead := file.Read(fileContent)

	return fileContent, errRead
}

func htmlHandler(htmlTemp *template.Template, stories StoryLines) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		url := r.RequestURI

		if url == "/" {
			http.Redirect(w, r, "/intro", http.StatusMovedPermanently)
		} else {
			arc := strings.TrimPrefix(url, "/")
			story, find := stories[arc]
			if find {
				htmlTemp.Execute(w, story)
			} else {
				fmt.Fprint(w, "404 - not found")
			}
		}
	}
}

func main() {
	//Readfile
	file, err := readFile("adventure.json")
	if err != nil {
		log.Fatal(err)
	}

	//Parse Json file to get the data
	stories, err1 := parseJsonFile(file)
	if err1 != nil {
		log.Fatal(err1)
	}

	htmlTemp, err2 := template.ParseFiles("arc.html")
	if err2 != nil {
		log.Fatal(err2)
	}

	http.ListenAndServe(":8080", htmlHandler(htmlTemp, stories))
}
