package main

import (
	"net/http"
	"registration/stack"

	"github.com/gin-gonic/gin"
)

func mainEngine() (engine *gin.Engine) {
	return gin.Default()
}

func withAppSetup(engine *gin.Engine) *gin.Engine {
	// setup basic gin default middlewares
	engine.Use(gin.Logger())
	engine.Use(gin.Recovery())

	// add access to assets and templates
	engine.Static("/assets", "./assets")
	engine.LoadHTMLGlob("templates/*")

	return engine
}

// Registration is a user registration, that is not submitted
type Registration struct {
	Firstname string
	Lastname  string
	Email     string
}

func registrationFromRequest(request *http.Request) (registration Registration) {
	request.ParseForm()

	firstname := request.Form.Get("firstname")
	lastname := request.Form.Get("lastname")
	email := request.Form.Get("email")

	registration = Registration{firstname, lastname, email}

	return registration
}

// AppEngine for registrations service
func AppEngine() *gin.Engine {
	engine := mainEngine()
	engine = withAppSetup(engine)

	var registrations stack.Stack

	engine.GET("/", func(c *gin.Context) {
		obj := gin.H{"title": "Main website"}
		c.HTML(200, "index.tmpl", obj)
	})

	engine.POST("/", func(c *gin.Context) {
		registration := registrationFromRequest(c.Request)

		registrations.Push(registration)

		obj := gin.H{"registrations": registrations}
		c.HTML(200, "success.tmpl", obj)
	})

	return engine
}

func main() {
	AppEngine().Run(":3000")
}
