package main

import (
	"github.com/go-martini/martini"
	"github.com/martini-contrib/binding"
	"github.com/martini-contrib/render"
)

// Registration is the minimal information a user has to provide
type Registration struct {
	Firstname string `form:"firstname" binding:"required"`
	Lastname  string `form:"lastname" binding:"required"`
	Email     string `form:"email" binding:"required"`
}

// Result is a combination of a given registration and the related errors
type Result struct {
	Registration Registration
	Errors       map[string]string
}

func errorMap(errors binding.Errors) map[string]string {
	errorMap := map[string]string{}

	for _, error := range errors {
		errorMap[error.FieldNames[0]] = error.Message
	}

	return errorMap
}

// AppEngine for web engine setup
func AppEngine() *martini.ClassicMartini {
	m := martini.Classic()
	m.Use(render.Renderer())

	m.Get("/", func(r render.Render) {
		r.HTML(200, "index", nil)
	})

	m.Post("/", binding.Form(Registration{}), func(r render.Render, errors binding.Errors, registration Registration) {
		errorMap := errorMap(errors)
		r.HTML(200, "index", Result{Registration: registration, Errors: errorMap})
	})

	return m
}

func main() {
	m := AppEngine()
	m.Run()
}
