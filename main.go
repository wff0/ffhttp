package main

import (
	"ffhttp/ff"
	"net/http"
)

func main() {
	r := ff.New()

	v1 := r.Group("/v1")
	{
		v1.GET("/", func(c *ff.Context) {
			c.HTML(http.StatusOK, "<h1>Hello ff</h1>")
		})
		v1.GET("/hello", func(c *ff.Context) {
			c.String(http.StatusOK, "hello %s, you're at %s\n", c.Query("name"), c.Path)
		})
	}

	v2 := r.Group("/v2")
	{
		v2.POST("/login", func(c *ff.Context) {
			c.Json(http.StatusOK, ff.H{
				"username": c.PostForm("username"),
				"password": c.PostForm("password"),
			})
		})

		v2.GET("/hello/:name", func(c *ff.Context) {
			c.String(http.StatusOK, "hello %s, you're at %s\n", c.Param("name"), c.Path)
		})

		v2.GET("/assert/*filepath", func(c *ff.Context) {
			c.Json(http.StatusOK, ff.H{
				"filepath": c.Param("filepath"),
			})
		})
	}

	r.Run(":9999")
}
