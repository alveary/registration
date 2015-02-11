package main

import (
	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()

	r.Use(gin.Logger())
	r.Use(gin.Recovery())

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

		message := "Hello " + firstname + " " + lastname
		c.String(200, message)
	})

	r.Run(":3000")
}
