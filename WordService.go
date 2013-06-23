// WordService.go
package main

import (
	"fmt"
	"github.com/ThatRubenAguilar/words"
	"net"
)

type WordService struct {
	listener        net.Listener
	word_dictionary *words.WordDictionary
}

type LookupResponse struct {
	Info  words.WordInfo
	Error error
}

type AddWordsResponse struct {
	Error error
}

func (ws *WordService) Setup() (err error) {
	if ws.word_dictionary == nil {
		filename := "words.txt"
		word_dict := new(words.WordDictionary)

		err = word_dict.AddWordsFromFile(filename)
		if err != nil {
			fmt.Println(err)
			return
		}
		ws.word_dictionary = word_dict
	} else {
		copy_dict := ws.word_dictionary.Copy()
		ws.word_dictionary = copy_dict
	}
	return
}

func (ws *WordService) Lookup(Word string) LookupResponse {
	response := LookupResponse{}
	var lookup_dict *words.WordDictionary
	lookup_dict = ws.word_dictionary

	if response.Info, response.Error = lookup_dict.Lookup(Word); response.Error != nil {
		fmt.Println(response.Error)
	}
	return response
}

func (ws *WordService) AddWords(Words []string) AddWordsResponse {
	response := AddWordsResponse{}
	copy_dict := ws.word_dictionary.Copy()
	if response.Error = copy_dict.AddWords(Words); response.Error != nil {
		return response
	}
	ws.word_dictionary = copy_dict
	return response
}
