// ServiceMain.go
package main

import (
	"encoding/xml"
	"fmt"
	"net"
	"net/http"
	"strings"
)

var word_service WordService
var EndService bool

func lookup_handler(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		fmt.Println(err)
	}

	word := r.FormValue("Word")
	var response LookupResponse = word_service.Lookup(word)
	responseXML, _ := xml.Marshal(response)
	fmt.Fprintf(w, string(responseXML))
}

func add_handler(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		fmt.Println(err)
	}

	words_raw := r.FormValue("Words")
	seperator := r.FormValue("Seperator")
	if seperator == "" {
		seperator = ","
	}
	words := strings.Split(words_raw, seperator)

	var response AddWordsResponse = word_service.AddWords(words)
	responseXML, _ := xml.Marshal(response)
	fmt.Fprintf(w, string(responseXML))
}

func service_handler(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		fmt.Println(err)
	}
	action := r.FormValue("Action")
	switch strings.ToLower(action) {
	case "restart":
		restart_service()
	case "stop":
		stop_service()
	}
}

func restart_service() {
	err := word_service.Setup()
	if err != nil {
		fmt.Println(err)
		stop_service()
		return
	}
	http.DefaultServeMux = http.NewServeMux()
}

func stop_service() {
	EndService = true
	if word_service.listener != nil {
		word_service.listener.Close()
	}
}

func main() {
	word_service = WordService{}
	err := word_service.Setup()
	if err != nil {
		fmt.Println(err)
		return
	}

	http.HandleFunc("/words/lookup", lookup_handler)
	http.HandleFunc("/words/add", add_handler)
	http.HandleFunc("/service", service_handler)

	for !EndService {
		if word_service.listener, err = net.Listen("tcp", "localhost:8080"); err != nil {
			fmt.Println(err)
			if word_service.listener != nil {
				word_service.listener.Close()
			}
		}
		if err = http.Serve(word_service.listener, nil); err != nil {
			fmt.Println(err)
			if word_service.listener != nil {
				word_service.listener.Close()
			}
		}
	}
}
