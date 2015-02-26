package main

import (
	"net/http"

	"github.com/go-martini/martini"
	validation "github.com/jamieomatthews/validation"
	"github.com/martini-contrib/binding"
	"github.com/martini-contrib/render"
)

// Registration is the minimal information a user has to provide
type Registration struct {
	Email    string `form:"email" binding:"required"`
	Password string `form:"password" binding:"required"`
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

// Validate password of a given registration
func (reg Registration) Validate(errors binding.Errors, req *http.Request) binding.Errors {
	v := validation.NewValidation(&errors, reg)

	v.Validate(&reg.Email).Classify("email-class").Email()
	v.Validate(&reg.Password).Key("password").MinLength(9)

	return *v.Errors.(*binding.Errors)
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
	m.RunOnAddr(":9000")
}
