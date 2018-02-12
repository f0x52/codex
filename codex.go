package main

import (
	"fmt"
	"github.com/go-macaron/binding"
	"github.com/go-macaron/pongo2"
	"github.com/go-redis/redis"
	"gopkg.in/macaron.v1"
	"strconv"
	"time"
)

var (
	storage   *redis.Client
	authLevel = 1
	site      Site
)

func index(ctx *macaron.Context) {
	ctx.HTML(200, "index")
}

func board(ctx *macaron.Context) {
	board := ctx.Params("board")
	if site.Boards[board] != 0 && site.Boards[board] <= authLevel { //higher authLevels for mod/admin only boards
		fmt.Println("visit " + board)
		var boardData Board
		boardData.Load(board)
		ctx.Data["board"] = boardData
		ctx.HTML(200, "board")
		return
	}
	ctx.Data["error"] = "404"
	ctx.HTML(404, "error")
}

func thread(ctx *macaron.Context) {
	board := ctx.Params("board")
	thread := ctx.Params("thread")
	ctx.Data["board"] = board
	ctx.Data["threadId"] = thread
	if site.Boards[board] != 0 && site.Boards[board] <= authLevel {
		var threadData Thread
		err := threadData.Load(board, thread)
		if err != nil {
			fmt.Println(err.Error())
			ctx.Data["error"] = "404"
			ctx.HTML(404, "error")
			return
		}
		ctx.Data["thread"] = threadData
		ctx.HTML(200, "thread")
		return
	}
	ctx.Data["error"] = "404"
	ctx.HTML(404, "error")
}

func post(data Submit, ctx *macaron.Context) { //TODO: DO SANITATION!!!!1!
	fmt.Printf("%+v", data)
	board := data.Board
	var thread string
	if site.Boards[board] != 0 && site.Boards[board] <= authLevel {
		var boardData Board
		boardData.Load(board)
		boardData.Count++
		if data.NewThread { //Creating a new thread
			threadData := Thread{Id: boardData.Count, Board: board, Name: data.Name, Content: data.Content, Bump: time.Now()}
			thread = strconv.Itoa(boardData.Count)
			threadData.Save()
			boardData.Threads = append(boardData.Threads, boardData.Count)
		} else { //Posting to an existing thread
			thread = strconv.Itoa(data.Thread)
			var threadData Thread
			err := threadData.Load(board, thread)
			if err != nil {
				fmt.Println(err.Error())
				ctx.Data["error"] = "That thread doesn't exist"
				ctx.HTML(404, "error")
				return
			}
			threadData.Posts = append(threadData.Posts, Post{boardData.Count, data.Content, time.Now()})
			threadData.Bump = time.Now()
			threadData.Save()
		}
		boardData.Save()
		ctx.Redirect("/" + board + "/" + thread)
	}
}

func main() {
	fmt.Println("Codex 0.01")
	site.Load()
	storage = redis.NewClient(&redis.Options{
		Addr:     site.Redis["address"],
		Password: "",
		DB:       0,
	})
	if _, err := storage.Ping().Result(); err != nil {
		fmt.Errorf("Can't ping redis: %s \n", err)
		return
	}

	m := macaron.New()
	m.Use(macaron.Logger())
	m.Use(macaron.Static("static"))
	m.Use(pongo2.Pongoer())
	m.Use(func(ctx *macaron.Context) {
		ctx.Data["site"] = site
	})

	m.Get("/", index)
	m.Get("/:board/", board)
	//m.Get("/:board/catalog", catalog)
	m.Get("/:board/:thread", thread)
	m.Post("/post", binding.BindIgnErr(Submit{}), post)

	m.Run()
}
