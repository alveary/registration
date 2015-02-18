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
}

func registrationFromRequest(request *http.Request) (registration Registration) {
	request.ParseForm()

	firstname := request.Form.Get("firstname")
	lastname := request.Form.Get("lastname")

	registration = Registration{firstname, lastname}

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

		// message := "Hello " + firstname + " " + lastname + " (RegistrationCount is now " + fmt.Sprintf("%v", regs.Len()) + ")"
		obj := gin.H{"registrations": registrations}
		c.HTML(200, "success.tmpl", obj)
	})

	return engine
}

func main() {
	AppEngine().Run(":3000")
}
