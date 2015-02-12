package main

import (
	"registration-service/registrations"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()

	r.Use(gin.Logger())
	r.Use(gin.Recovery())

	var regs registrations.Registrations

	r.Static("/assets", "./assets")

	r.LoadHTMLGlob("templates/*")

	r.GET("/registrations/new", func(c *gin.Context) {
		obj := gin.H{"title": "Main website"}
		c.HTML(200, "index.tmpl", obj)
	})

	r.POST("/registrations", func(c *gin.Context) {
		c.Request.ParseForm()

		firstname := c.Request.Form.Get("firstname")
		lastname := c.Request.Form.Get("lastname")

		reg := registrations.Registration{firstname, lastname}

		regs.Push(reg)

		// message := "Hello " + firstname + " " + lastname + " (RegistrationCount is now " + fmt.Sprintf("%v", regs.Len()) + ")"
		obj := gin.H{"firstname": firstname, "lastname": lastname, "registrations": regs}
		c.HTML(200, "success.tmpl", obj)
	})

	r.Run(":3000")
}
