package main

import (
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"strings"
	"time"
)

var (
	linkList map[string]string
)

func init() {
	seed := time.Now().UnixNano()
	src := rand.NewSource(seed)
	rand.New(src)
}

func main() {
	linkList = map[string]string{}

	http.HandleFunc("/add-link", AddLink)
	http.HandleFunc("/short", GetLink)

	log.Fatal(http.ListenAndServe(":9099", nil))
}

func AddLink(w http.ResponseWriter, r *http.Request) {
	values := r.URL.Query()
	key, ok := values["link"]
	if !ok {
		w.WriteHeader(http.StatusConflict)
		fmt.Fprintln(w, "Failed to add link")
		return
	}

	if _, ok := linkList[key[0]]; ok {
		w.WriteHeader(http.StatusConflict)
		fmt.Fprintln(w, "Already have this link")
		return
	}

	genString := fmt.Sprint(rand.Int63n(1000))
	linkList[genString] = key[0]
	w.Header().Set("Content-Type", "text/html")
	w.WriteHeader(http.StatusAccepted)
	linkString := fmt.Sprintf("<a href=\"http://localhost:9099/short/%s\">http://localhost:9099/short/%s</a>", genString, genString)
	fmt.Fprintln(w, "Added shortlink")
	fmt.Fprintln(w, linkString)
}

func GetLink(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Path
	paths := strings.Split(path, "/")
	log.Printf("Redirected to %s\n", linkList[paths[2]])
	http.Redirect(w, r, linkList[paths[2]], http.StatusPermanentRedirect)
}
