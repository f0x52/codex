package main

import (
	//"bytes"
	"fmt"
	"github.com/go-macaron/pongo2"
	"gopkg.in/macaron.v1"
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
	m.Run()
}
