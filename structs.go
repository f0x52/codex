package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
	"time"
)

type Site struct {
	Name   string            `json:"name"`
	Boards map[string]int    `json:"boards"`
	Alias  map[string]string `json:"alias"`
	Redis  map[string]string `json:"redis"`
}

type Board struct {
	Name    string `json:"name"`
	Uri     string `json:"uri"`
	Desc    string `json:"desc"`
	Count   int    `json:"count"`
	Threads []int  `json:"threads"`
}

type Thread struct {
	Id      int       `json:"id"`
	Board   string    `json:"board"`
	Name    string    `json:"name"`
	Content string    `json:"content"`
	Posts   []Post    `json:"posts"`
	Bump    time.Time `json:"bump"`
}

type Post struct {
	Id        int       `json:"id"`
	Content   string    `json:"content"`
	Timestamp time.Time `json:"time"`
}

type Submit struct {
	Board     string `binding:"Required"`
	Content   string `binding:"Required"`
	Thread    int
	NewThread bool `binding:"Required"`
	Name      string
}

func (site *Site) Load() {
	_, err := os.Stat("site.json")
	if err != nil {
		str := "{\"name\": \"Codex\", \"boards\": {\"meta\": 1}, \"redis\": {\"address\": \"localhost:6379\"}}"
		ioutil.WriteFile("site.json", []byte(str), 0644)
		fmt.Println("Edit site details in site.json")
		os.Exit(1)
	}
	file, _ := os.Open("site.json")
	decoder := json.NewDecoder(file)
	err = decoder.Decode(&site)
	if err != nil {
		fmt.Println(err.Error())
	}
}

func (board *Board) Save() {
	json, err := json.Marshal(board)
	if err != nil {
		fmt.Println(err)
	}
	if err = storage.Set("board: "+board.Name, json, 0).Err(); err != nil {
		fmt.Println("Storage error: " + err.Error())
	}
}

func (board *Board) Load(name string) {
	redisData, err := storage.Get("board: " + name).Result()
	if err != nil {
		fmt.Println("Loading error: " + err.Error())
	}
	json.Unmarshal([]byte(redisData), &board)
}

func (thread *Thread) Save() {
	json, err := json.Marshal(thread)
	if err != nil {
		fmt.Println(err)
	}
	if err = storage.Set("thread: "+thread.Board+"/"+strconv.Itoa(thread.Id), json, 0).Err(); err != nil {
		fmt.Println("Storage error: " + err.Error())
	}
}

func (thread *Thread) Load(board string, id string) error {
	redisData, err := storage.Get("thread: " + board + "/" + id).Result()
	if err != nil {
		return err
	}
	json.Unmarshal([]byte(redisData), &thread)
	return nil
}
