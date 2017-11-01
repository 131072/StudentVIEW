package main

import (
	"github.com/gin-contrib/gzip"
	"github.com/gin-gonic/autotls"
	"github.com/gin-gonic/gin"
)

var (
	//Domain cannot be an IP Address, unless you are willing to sacrifice HTTPS
	domain    = "localhost"
	subdomain = ""
)

func main() {
	if domain == "" || domain == "localhost" {
		http := gin.Default()
		http.Static("/studentview", "./public")
		http.Use(gzip.Gzip(gzip.DefaultCompression))
		http.Run(":80")
	} else {
		http := gin.Default()
		https := gin.Default()
		https.Static("/studentview", "./public")
		http.Use(gzip.Gzip(gzip.DefaultCompression))
		https.Use(gzip.Gzip(gzip.DefaultCompression))
		http.GET("/*path", func(c *gin.Context) {
			c.Redirect(302, "https://"+domain+subdomain+c.Param("variable"))
		})

		go autotls.Run(https, domain)
		http.Run(":80")
	}
}
