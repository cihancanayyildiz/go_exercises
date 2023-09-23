package main

import (
	"encoding/json"
	"html/template"
	"log"
	"net/http"
	"os"
)

type StoryArc struct {
    Title string `json:"title"`
    Story []string `json:"story"`
    Options []struct{
        Text string `json:"text"`
        Arc string `json:"arc"`
    } `json:"options"`
}
type Story map[string]StoryArc

var tmplt *template.Template
func mapHandler(stories Story) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        if r.Method == "GET" {
            tmplt, _ = template.ParseFiles("cyoa.html")
            path := r.URL.Path[1:]
            if path == "" || path == "/" {
                err := tmplt.Execute(w, stories["intro"])
                if err != nil {
                    http.Error(w, "Something went wrong...", http.StatusInternalServerError)
                }
            }

            for k, v := range stories {
                if k == path {
                    err := tmplt.Execute(w, v)
                    if err != nil {
                        http.Error(w, "Something went wrong...", http.StatusInternalServerError)
                    }
                }
            }
        }
    }
}

func main() {
    var err error
    content, err := os.ReadFile("gopher.json")
    if err != nil {
        log.Fatal("error reading file: ", err)
    }

    var story Story
    err = json.Unmarshal(content, &story)

    if err != nil {
		log.Fatal("error unmarshaling json: ", err)
	}

	err = http.ListenAndServe("localhost:8080", mapHandler(story))
 
    if err != nil {
        log.Fatalln("There's an error with the server:", err)
    }
}