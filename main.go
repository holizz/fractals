package main

import (
	"github.com/codegangsta/martini"
	"github.com/martini-contrib/render"
)

func main() {
	m := martini.Classic()

	m.Use(render.Renderer(render.Options{
		Layout: "layout",
	}))

	m.Get("/", func (rr render.Render) {
		rr.HTML(200, "index", nil)
	})

	m.Get("/lsystem", handleLSystem)
	m.Get("/lsystem.png", handleLSystemPng)

	m.Run()
}
