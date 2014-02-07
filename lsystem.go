package main

import (
	"errors"
	"github.com/holizz/terrapin"
	"github.com/martini-contrib/render"
	"image"
	"image/png"
	"net/http"
	"net/url"
	"regexp"
	"strconv"
	"strings"
)

type LSystem struct {
	Definitions [][3]string
	Rules [][2]string
	Iterations int
}

func (sys LSystem) ParseForm(form url.Values) error {
	defReg := regexp.MustCompile("^(.) = (.*)\\((.*)\\)$")

	for _, line := range strings.Split(form.Get("definitions"), "\r\n") {
		s := defReg.FindAllStringSubmatch(line, -1)
		if len(s) != 1 || len(s[0]) != 4 {
			return errors.New("could not parse definitions")
		}
		sys.Definitions = append(sys.Definitions, [3]string{s[0][1], s[0][2], s[0][3]})
	}

	ruleReg := regexp.MustCompile("^(.) -> (.*)$")

	for _, line := range strings.Split(form.Get("rules"), "\r\n") {
		s := ruleReg.FindAllStringSubmatch(line, -1)
		if len(s) != 1 || len(s[0]) != 3 {
			return errors.New("could not parse rules")
		}
		sys.Rules = append(sys.Rules, [2]string{s[0][1], s[0][2]})
	}

	iterations, err := strconv.Atoi(form.Get("iterations"))
	if err != nil {
		panic(err)
	}
	sys.Iterations = iterations

	return nil
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
	err = sys.ParseForm(r.Form)
	if err != nil {
		panic(err)
	}

	// Execute lsystem

	i := image.NewRGBA(image.Rect(0, 0, 300, 300))

	t := terrapin.NewTerrapin(i, terrapin.Position{150.0, 150.0})

	t.Forward(20)

	png.Encode(w, i)
}