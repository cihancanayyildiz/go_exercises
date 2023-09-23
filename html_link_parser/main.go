package html_link_parser

import (
	"flag"
	"fmt"
	"log"
	"os"
	"strings"

	"golang.org/x/net/html"
)

type Link struct {
	Href string
	Text string
}	

func Html_link_parser(file []byte) []Link {
	tkn := html.NewTokenizer(strings.NewReader(string(file)))

	var link []Link

	var tempLink = Link{
		Href: "",
		Text: "",
	}

	var linkFlag bool = false
	var htmlEndFlag bool = false
	for {
		tt := tkn.Next()
		switch {
			case tt == html.ErrorToken:
				htmlEndFlag = true

			case tt == html.StartTagToken:
				t := tkn.Token()
				if t.Data == "a" {
					tempLink.Href = ""
					tempLink.Text= ""
					linkFlag = true
					for i := range t.Attr {
						if(t.Attr[i].Key == "href"){
							tempLink.Href = t.Attr[i].Val
							break
						}
					}
				}

        	case tt == html.TextToken:
				t := tkn.Token()
				if linkFlag {
					t.Data = strings.TrimSpace(t.Data)
					tempLink.Text += t.Data
				}

			case tt == html.EndTagToken:
				t := tkn.Token()
				if t.Data == "a" {
					linkFlag = false
					link = append(link, tempLink)
				}	
		}
		if htmlEndFlag {
			break
		}
	}
	return link
}

func main() {
	filepath := flag.String("file", "ex1.html", "html file path")
	flag.Parse()
	file, err := os.ReadFile(*filepath)

	if err != nil {
		log.Fatalln(err)
	}

	link := Html_link_parser(file)

	fmt.Println(link)
}