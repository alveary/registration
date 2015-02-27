package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/go-martini/martini"
	validation "github.com/jamieomatthews/validation"
	"github.com/martini-contrib/binding"
	"github.com/martini-contrib/render"
)

// Registration is the minimal information a user has to provide
type Registration struct {
	Email    string `form:"email" json:"email" binding:"required"`
	Password string `form:"password" json:"password" binding:"required"`
}

// Result is a combination of a given registration and the related errors
type Result struct {
	Registration Registration
	Errors       map[string]string
}

func errorMap(errors binding.Errors) map[string]string {
	errorMap := map[string]string{}

	for _, error := range errors {
		if len(error.FieldNames) > 0 {
			errorMap[error.FieldNames[0]] = error.Message
		}
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

// RequestRegistration ...
func RequestRegistration(registration Registration, r render.Render) (target string, err error) {
	failure := make(chan error)
	success := make(chan string)
	defer func() {
		close(failure)
		close(success)
	}()

	go func() {
		json, err := json.Marshal(registration)
		if err != nil {
			failure <- err
			return
		}

		// TODO: handle target
		_, requestErr := http.Post("http://localhost:9001/", "application/json", bytes.NewBuffer(json))

		if requestErr != nil {
			failure <- requestErr
			return
		}

		success <- "target"
	}()

	select {
	case result := <-success:
		return result, nil
	case exception := <-failure:
		return "", exception
	case <-time.After(time.Second * 3):
		return "", errors.New("timeout")
	}
}

// AppEngine for web engine setup
func AppEngine() *martini.ClassicMartini {
	m := martini.Classic()
	m.Use(render.Renderer())

	m.Get("/", func(r render.Render) {
		r.HTML(200, "index", nil)
	})

	m.Head("/alive", func(resp http.ResponseWriter) {
		resp.WriteHeader(http.StatusOK)
	})

	m.Post("/", binding.Form(Registration{}), func(r render.Render, errors binding.Errors, registration Registration, resp http.ResponseWriter) {
		if errors.Len() > 0 {

			errorMap := errorMap(errors)
			r.HTML(200, "index", Result{Registration: registration, Errors: errorMap})

		} else {

			target, err := RequestRegistration(registration, r)
			fmt.Println(err)

			if err != nil {
				resp.WriteHeader(http.StatusInternalServerError)
				r.HTML(200, "failure", err)
				return
			}

			r.HTML(200, "success", target)
		}
	})

	return m
}

func init() {
	serviceentry := struct {
		Name  string
		Root  string
		Alive string
	}{
		"reg",
		"https://alveary-user-registration.herokuapp.com",
		"https://alveary-user-registration.herokuapp.com/alive",
	}

	json, _ := json.Marshal(serviceentry)

	http.Post("https://alveary-overseer.herokuapp.com/", "application/json", bytes.NewBuffer(json))
}

func main() {
	var port int
	flag.IntVar(&port, "p", 9000, "the port number")
	flag.Parse()

	m := AppEngine()
	m.RunOnAddr(":" + strconv.Itoa(port))
}
