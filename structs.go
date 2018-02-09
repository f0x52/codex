package main

import (
	"encoding/json"
	"fmt"
	"os"
)

func load(path string) *os.File {
	fmt.Println("loading: " + path)
	file, err := os.Open(path)
	if err != nil {
		fmt.Println(err.Error())
	}
	return file
}

func loadSiteJSON() Site {
	path := "json/site.json"
	file := load(path)
	decoder := json.NewDecoder(file)
	var JSON Site
	err := decoder.Decode(&JSON)
	if err != nil {
		fmt.Println(err.Error())
	}
	return JSON
}

func loadBoardJSON(board string) Board {
	path := fmt.Sprintf("json/boards/%s/%s.json", board, board)
	file := load(path)
	decoder := json.NewDecoder(file)
	var JSON Board
	err := decoder.Decode(&JSON)
	if err != nil {
		fmt.Println(err.Error())
	}
	return JSON
}

func loadThreadJSON(board string, thread string) (Thread, error) {
	path := fmt.Sprintf("json/boards/%s/%s.json", board, thread)
	var JSON Thread
	_, err := os.Stat(path)
	if err == nil {
		file := load(path)
		decoder := json.NewDecoder(file)
		err = decoder.Decode(&JSON)
		if err != nil {
			fmt.Println(err.Error())
		}
	}
	return JSON, err
}

type Site struct {
	Name   string            `json:"name"`
	Boards map[string]int    `json:"boards"`
	Alias  map[string]string `json:"alias"`
}

type Board struct {
	Name string `json:"name"`
	Uri  string `json:"uri"`
	Desc string `json:"desc"`
}

type Thread struct {
	Id      int    `json:"id"`
	Name    string `json:"name"`
	Content string `json:"content"`
	Posts   []Post `json:"posts"`
}

type Post struct {
	Id      int    `json:"id"`
	Content string `json:"content"`
}
