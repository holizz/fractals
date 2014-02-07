package main

import (
	"errors"
	"github.com/holizz/terrapin"
	"github.com/martini-contrib/render"
	"image"
	"image/png"
	"math"
	"net/http"
	"net/url"
	"regexp"
	"strconv"
	"strings"
)

type LSystem struct {
	Definitions map[string][2]string
	Rules map[string]string
	StartState string
	Iterations int
}

func (sys *LSystem) ParseForm(form url.Values) error {
	defReg := regexp.MustCompile("^(.) = (.*)\\((.*)\\)$")
	sys.Definitions = make(map[string][2]string)

	for _, line := range strings.Split(form.Get("definitions"), "\r\n") {
		s := defReg.FindAllStringSubmatch(line, -1)
		if len(s) != 1 || len(s[0]) != 4 {
			return errors.New("could not parse definitions")
		}
		sys.Definitions[s[0][1]] = [2]string{s[0][2], s[0][3]}
	}

	ruleReg := regexp.MustCompile("^(.) -> (.*)$")
	sys.Rules = make(map[string]string)

	for _, line := range strings.Split(form.Get("rules"), "\r\n") {
		s := ruleReg.FindAllStringSubmatch(line, -1)
		if len(s) != 1 || len(s[0]) != 3 {
			return errors.New("could not parse rules")
		}
		sys.Rules[s[0][1]] = s[0][2]
	}

	iterations, err := strconv.Atoi(form.Get("iterations"))
	if err != nil {
		panic(err)
	}
	sys.Iterations = iterations

	sys.StartState = form.Get("startstate")

	return nil
}

func (sys *LSystem) Execute(t *terrapin.Terrapin) {
	// Rewrite

	state := sys.StartState

	for i := 0; i < sys.Iterations; i++ {
		for from, to := range sys.Rules {
			state = strings.Replace(state, from, strings.ToLower(to), -1)
		}

		state = strings.ToUpper(state)
	}

	// Run turtle

	for _, a := range state {
		pair := sys.Definitions[string(a)]
		fn := pair[0]
		val := pair[1]
		v, err := strconv.Atoi(val)
		if err != nil {
			panic(err)
		}

		vv := float64(v)
		vRad := vv * (math.Pi / 180)

		switch fn {
		case "fwd":
			t.Forward(vv)
		case "left":
			t.Left(vRad)
		case "right":
			t.Right(vRad)
		}
	}
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
		"startstate": r.Form.Get("startstate"),
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

	sys.Execute(t)

	png.Encode(w, i)
}
