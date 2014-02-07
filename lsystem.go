package main

import (
	"github.com/holizz/terrapin"
	"github.com/martini-contrib/render"
	"image"
	"image/png"
	"net/http"
	"regexp"
	"strconv"
	"strings"
)

type LSystem struct {
	Definitions [][3]string
	Rules [][2]string
	Iterations int
}

func handleLSystem(r *http.Request, rr render.Render) {
	err := r.ParseForm()
	if err != nil {
		panic(err)
	}

	rr.HTML(200, "lsystem", map[string]string{
		"definitions": r.Form.Get("definitions"),
		"rules": r.Form.Get("rules"),
		"iterations": r.Form.Get("iterations"),
	})
}

func handleLSystemPng(w http.ResponseWriter, r *http.Request, rr render.Render) {
	err := r.ParseForm()
	if err != nil {
		panic(err)
	}

	// Parse form into an LSystem
	sys := LSystem{}

	defReg := regexp.MustCompile("^(.) = (.*)\\((.*)\\)$")

	for _, line := range strings.Split(r.Form.Get("definitions"), "\r\n") {
		s := defReg.FindAllStringSubmatch(line, -1)
		if len(s) != 1 || len(s[0]) != 4 {
			rr.HTML(400, "error", "could not parse definitions")
			return
		}
		sys.Definitions = append(sys.Definitions, [3]string{s[0][1], s[0][2], s[0][3]})
	}

	ruleReg := regexp.MustCompile("^(.) -> (.*)$")

	for _, line := range strings.Split(r.Form.Get("rules"), "\r\n") {
		s := ruleReg.FindAllStringSubmatch(line, -1)
		if len(s) != 1 || len(s[0]) != 3 {
			rr.HTML(400, "error", "could not parse definitions")
			return
		}
		sys.Rules = append(sys.Rules, [2]string{s[0][1], s[0][2]})
	}

	iterations, err := strconv.Atoi(r.Form.Get("iterations"))
	if err != nil {
		panic(err)
	}
	sys.Iterations = iterations

	// Execute lsystem

	i := image.NewRGBA(image.Rect(0, 0, 300, 300))

	t := terrapin.NewTerrapin(i, terrapin.Position{150.0, 150.0})

	t.Forward(20)

	png.Encode(w, i)
}
