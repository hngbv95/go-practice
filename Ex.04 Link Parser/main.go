package main

import (
	"bytes"
	"flag"
	"fmt"
	"log"
	"os"

	"golang.org/x/net/html"
)

type Link struct {
	Href string
	Text string
}

func ReadFile(fileName string) (*[]byte, error) {
	file, err := os.Open(fileName)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	stat, err := file.Stat()
	if err != nil {
		return nil, err
	}

	content := make([]byte, stat.Size())
	_, errRead := file.Read(content)
	if errRead != nil {
		return nil, errRead
	}

	return &content, nil
}

func ExtractNodeLink(node *html.Node) *Link {
	if node.Type != html.ElementNode || node.Data != "a" {
		return nil
	}

	link := Link{}
	link.Text = node.FirstChild.Data
	link.Href = ""
	for _, attr := range node.Attr {
		if attr.Key == "href" {
			link.Href = attr.Val
			break
		}
	}

	return &link
}

func CrawLinkChan(node *html.Node) <-chan (*Link) {
	ch := make(chan *Link)

	go func() {
		defer close(ch)
		// Dept-First-Search
		stack := []*html.Node{node}

		for len(stack) > 0 {
			current := stack[len(stack)-1]
			stack = stack[:len(stack)-1]

			if current.Type == html.ElementNode && current.Data == "a" {

				link := ExtractNodeLink(current)
				ch <- (link)
			} else {
				for child := current.FirstChild; child != nil; child = child.NextSibling {
					stack = append(stack, child)
				}
			}
		}

	}()

	return ch
}

func ParseFlag() string {
	fileName := flag.String("f", "ext1.html", "Html file name")
	flag.Parse()

	return *fileName
}

func main() {
	fileName := ParseFlag()

	// Read file content
	content, err := ReadFile(fileName)
	if err != nil {
		log.Fatal(err)
	}

	rootNode, err := html.Parse(bytes.NewReader(*content))
	if err != nil {
		log.Fatal(err)
	}

	ch := CrawLinkChan(rootNode)

	for item := range ch {
		fmt.Println(*item)
	}
}
