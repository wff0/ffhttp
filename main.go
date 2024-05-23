package main

import (
	"ffhttp/ff"
	"net/http"
)

func main() {
	r := ff.New()
	r.GET("/", func(c *ff.Context) {
		c.HTML(http.StatusOK, "<h1>Hello ff</h1>")
	})

	r.GET("/hello", func(c *ff.Context) {
		c.String(http.StatusOK, "hello %s, you're at %s\n", c.Query("name"), c.Path)
	})

	r.POST("/login", func(c *ff.Context) {
		c.Json(http.StatusOK, ff.H{
			"username": c.PostForm("username"),
			"password": c.PostForm("password"),
		})
	})

	r.Run(":9999")
}
