package main

import (
	"ffhttp/ff"
	"log"
	"net/http"
	"time"
)

func onlyForV2() ff.HandlerFunc {
	return func(c *ff.Context) {
		t := time.Now()

		c.Next()

		log.Printf("[%d] %s in %v for group v2", c.StatusCode, c.Req.RequestURI, time.Since(t))
	}
}

func main() {
	r := ff.New()
	r.Use(ff.Logger())

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
	v2.Use(onlyForV2())
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
