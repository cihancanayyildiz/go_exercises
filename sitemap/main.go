package main

import (
	"bytes"
	"encoding/xml"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"

	"golang.org/x/net/html"
)

type Link struct {
	Href string
	Text string
}

type Urlset struct {
	XMLName xml.Name `xml:"urlset"`
	Url     []*Url   `xml:"url"`
}

type Url struct {
	Loc string `xml:"loc"`
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
				tempLink.Text = ""
				linkFlag = true
				for i := range t.Attr {
					if t.Attr[i].Key == "href" {
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
	url := flag.String("url", "https://www.sitemaps.org", "Url")

	flag.Parse()

	resp, err := http.Get(*url)
	if err != nil {
		log.Fatalln(err)
	}

	body, err := io.ReadAll(resp.Body)

	if err != nil {
		log.Fatalln(err)
	}

	links := Html_link_parser(body)

	fmt.Println(links[1].Href)

	urlset := &Urlset{}
	for _, v := range links {
		if v.Href != "" {
			singleUrl := &Url{
				Loc: v.Href,
			}
			urlset.Url = append(urlset.Url, singleUrl)
		}
	}

	out, err := xml.Marshal(urlset)
	if err != nil {
		log.Fatalln(err)
	}

	w := &bytes.Buffer{}
	w.WriteString(xml.Header)
	w.Write(out)

	// writer, _ := os.Open("/output.xml")
	// encoder := xml.NewEncoder(writer)
	// error := encoder.Encode(string(out))
	// if error != nil {
	// 	log.Fatalln(error)
	// }

	_ = os.WriteFile("output.xml", w.Bytes(), 0644)

}
