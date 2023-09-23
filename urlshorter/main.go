package main

import (
	"fmt"
	"net/http"

	"gopkg.in/yaml.v2"
)

type YAMLMap struct {
	Path string `yaml:"path"`
	Url string `yaml:"url"`
}

func MapHandler(pathsToUrls map[string]string, fallback http.Handler) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request){
		url:= pathsToUrls[request.URL.Path]
		if url != "" {
			http.Redirect(writer, request, url, http.StatusFound)
			return
		}
		fallback.ServeHTTP(writer, request)
	}
}

func YAMLHandler(yml []byte, fallback http.Handler) (http.HandlerFunc, error) {
	var ymlMap []YAMLMap
	err := yaml.Unmarshal(yml, &ymlMap)
	if err != nil {
		return nil, err
	}

	paths := make(map[string]string)
	for _, v := range ymlMap {
		paths[v.Path] = v.Url
	}

	return MapHandler(paths, fallback), nil
}

func main() {
	mux := defaultMux()

	// Build the MapHandler using the mux as the fallback
	pathsToUrls := map[string]string{
		"/urlshort-godoc": "https://godoc.org/github.com/gophercises/urlshort",
		"/yaml-godoc":     "https://godoc.org/gopkg.in/yaml.v2",
	}

 	mapHandler := MapHandler(pathsToUrls, mux)

	// Build the YAMLHandler using the mapHandler as the
	// fallback
	yaml := `
- path: /urlshort
  url: https://github.com/gophercises/urlshort
- path: /urlshort-final
  url: https://github.com/gophercises/urlshort/tree/solution
`
	yamlHandler, err := YAMLHandler([]byte(yaml), mapHandler)
	if err != nil {
		panic(err)
	}
	fmt.Println("Starting the server on :8080")
	http.ListenAndServe(":8080", yamlHandler)
}

func defaultMux() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("/", hello)
	return mux
}

func hello(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Hello, Cihan!")
}