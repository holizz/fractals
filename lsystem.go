package main

import (
	"errors"
	"github.com/holizz/terrapin"
	"github.com/martini-contrib/render"
	"image"
	"image/png"
	"io"
	"math"
	"net/http"
	"net/url"
	"regexp"
	"strconv"
	"strings"
)

type LSystem struct {
	Definitions map[string]Definition
	Rules       map[string]string
	StartState  string
	Iterations  int
}

type Definition struct {
	Function string
	Value    int
}

func (sys *LSystem) ParseForm(form url.Values) error {
	defReg := regexp.MustCompile("^(.) = (.*)\\((.*)\\)$")
	sys.Definitions = make(map[string]Definition)

	for _, line := range strings.Split(form.Get("definitions"), "\r\n") {
		s := defReg.FindAllStringSubmatch(line, -1)
		if len(s) != 1 || len(s[0]) != 4 {
			return errors.New("could not parse definitions")
		}

		val, err := strconv.Atoi(s[0][3])
		if err != nil {
			return errors.New("could not parse integer")
		}

		sys.Definitions[s[0][1]] = Definition{
			Function: s[0][2],
			Value:    val,
		}
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
		def := sys.Definitions[string(a)]

		value := float64(def.Value)
		valueRad := value * (math.Pi / 180)

		switch def.Function {
		case "fwd":
			t.Forward(value)
		case "left":
			t.Left(valueRad)
		case "right":
			t.Right(valueRad)
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
		"rules":       r.Form.Get("rules"),
		"iterations":  r.Form.Get("iterations"),
		"startstate":  r.Form.Get("startstate"),
	})
}

func handleLSystemPng(w http.ResponseWriter, r *http.Request, rr render.Render) {
	err := r.ParseForm()
	if err != nil {
		turtleError(w)
		return
	}

	// Parse form into an LSystem
	sys := LSystem{}
	err = sys.ParseForm(r.Form)
	if err != nil {
		turtleError(w)
		return
	}

	// Execute lsystem

	size := 780

	i := image.NewRGBA(image.Rect(0, 0, size, size))
	t := terrapin.NewTerrapin(i, terrapin.Position{
		float64(size) / 2,
		float64(size) / 2,
	})

	sys.Execute(t)

	png.Encode(w, i)
}

func turtleError(w io.Writer) {
	size := 780
	x := 50.0

	i := image.NewRGBA(image.Rect(0, 0, size, size))
	t := terrapin.NewTerrapin(i, terrapin.Position{
		float64(size) / 2 - x * 5,
		float64(size) / 2 - x * 5,
	})

	for _, c := range "ERROR" {
		switch c {
		case 'E':
			t.Right(math.Pi / 2)
			t.Forward(x)
			t.Right(math.Pi)
			t.Forward(x)
			t.Left(math.Pi / 2)
			t.Forward(x)
			t.Left(math.Pi / 2)
			t.Forward(x)
			t.Right(math.Pi)
			t.Forward(x)
			t.Left(math.Pi / 2)
			t.Forward(x)
			t.Left(math.Pi / 2)
			t.Forward(x * 2)
			t.Left(math.Pi / 2)
		case 'R':
			t.Forward(x * 2)
			for i := 0; i < 3; i++ {
				t.Right(math.Pi / 2)
				t.Forward(x)
			}
			t.Left(math.Pi * 3 / 4)
			t.Forward(math.Hypot(x, x))
			t.Left(math.Pi / 4)
			t.Forward(x)
			t.Left(math.Pi / 2)
		case 'O':
			for i := 0; i < 2; i++ {
				t.Forward(x * 2)
				t.Right(math.Pi / 2)
				t.Forward(x)
				t.Right(math.Pi / 2)
			}
			t.Right(math.Pi / 2)
			t.Forward(x * 2)
			t.Left(math.Pi / 2)
		}
	}

	png.Encode(w, i)
}
