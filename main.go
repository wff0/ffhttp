package main

import (
	"ffhttp/ff"
	"fmt"
	"html/template"
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

type student struct {
	Name string
	Age  int
}

func FormatAsDate(t time.Time) string {
	year, month, day := t.Date()
	return fmt.Sprintf("%d-%02d-%02d", year, month, day)
}

func main() {
	r := ff.New()
	r.Use(ff.Logger())

	r.SetFuncMap(template.FuncMap{
		"FormatAsDate": FormatAsDate,
	})
	r.LoadHTMLGlob("templates/*")
	r.Static("/assets", "./static")

	stu1 := &student{Name: "fengfan", Age: 18}
	stu2 := &student{Name: "jack", Age: 22}

	r.GET("/", func(c *ff.Context) {
		c.HTML(http.StatusOK, "css.tmpl", nil)
	})

	r.GET("/students", func(c *ff.Context) {
		c.HTML(http.StatusOK, "arr.tmpl", ff.H{
			"title":  "ff",
			"stuArr": [2]*student{stu1, stu2},
		})
	})

	r.GET("/date", func(c *ff.Context) {
		c.HTML(http.StatusOK, "custom_func.tmpl", ff.H{
			"title": "ff",
			"now":   time.Date(2024, 1, 20, 0, 0, 0, 0, time.UTC),
		})
	})

	v1 := r.Group("/v1")
	{
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

		//v2.GET("/assert/*filepath", func(c *ff.Context) {
		//	c.Json(http.StatusOK, ff.H{
		//		"filepath": c.Param("filepath"),
		//	})
		//})
	}

	r.Run(":9999")
}
