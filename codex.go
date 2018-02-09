package main

import (
	"fmt"
	"github.com/go-macaron/binding"
	"github.com/go-macaron/pongo2"
	"gopkg.in/macaron.v1"
	"strconv"
)

var (
	site      = loadSiteJSON()
	authLevel = 1
)

func index(ctx *macaron.Context) {
	ctx.HTML(200, "index")
}

func board(ctx *macaron.Context) {
	board := ctx.Params("board")
	ctx.Data["board"] = board
	if site.Boards[board] != 0 && site.Boards[board] <= authLevel { //higher authLevels for mod/admin only boards
		fmt.Println("visit " + board)
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
		threadJSON, err := loadThreadJSON(board, thread)
		ctx.Data["thread"] = threadJSON
		if err != nil {
			fmt.Println(err.Error())
			ctx.Data["error"] = "404"
			ctx.HTML(404, "error")
			return
		}
		ctx.HTML(200, "thread")
		return
	}
	ctx.Data["error"] = "404"
	ctx.HTML(404, "error")
}

func catalog(ctx *macaron.Context) {
	board := ctx.Params("board")
	if site.Boards[board] != 0 && site.Boards[board] <= authLevel { //higher authLevels for mod/admin only boards
		ctx.HTML(200, "catalog")
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
		boardJSON, err := loadBoardJSON(board)
		if err != nil {
			fmt.Println(err.Error())
			ctx.Data["error"] = "That board doesn't exist"
			ctx.HTML(404, "error")
			return
		}
		boardJSON.Count++
		saveBoardJSON(boardJSON, board)
		if !data.NewThread { //Posting to a thread
			thread = strconv.Itoa(data.Thread)
			threadJSON, err := loadThreadJSON(board, thread)
			if err != nil {
				fmt.Println(err.Error())
				ctx.Data["error"] = "That thread doesn't exist"
				ctx.HTML(404, "error")
				return
			}
			threadJSON.Posts = append(threadJSON.Posts, Post{boardJSON.Count, data.Content})
			saveThreadJSON(threadJSON, board, thread)
		} else { //Creating a new thread
			threadJSON := Thread{Id: boardJSON.Count, Name: data.Name, Content: data.Content}
			thread = strconv.Itoa(boardJSON.Count)
			saveThreadJSON(threadJSON, board, thread)
		}
		ctx.Redirect("/" + board + "/" + thread)
	}
}

func main() {
	fmt.Println("Codex 0.01")
	m := macaron.New()
	m.Use(macaron.Logger())
	m.Use(macaron.Static("static"))
	m.Use(pongo2.Pongoer())
	m.Use(func(ctx *macaron.Context) {
		ctx.Data["site"] = site
	})

	m.Get("/", index)
	m.Get("/:board/", board)
	m.Get("/:board/catalog", catalog)
	m.Get("/:board/:thread", thread)
	m.Post("/post", binding.BindIgnErr(Submit{}), post)
	m.Run()
}
