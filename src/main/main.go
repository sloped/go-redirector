package main

import (
	"bufio"
	"net/http"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()
	router.GET("/healthcheck", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "ok",
		})
	})
	
	file, err := os.Open("redirects")
    if err != nil {
        panic(err)
    }
    defer file.Close()

    scanner := bufio.NewScanner(file)
    for scanner.Scan() {
        fields := strings.Fields(scanner.Text())
        if len(fields) == 2 {
            targetURL := fields[1]
            router.GET(fields[0], func(c *gin.Context) {
                c.Redirect(http.StatusFound, targetURL)
            })
        }
    }

	if err := scanner.Err(); err != nil {
        panic(err)
    }

	router.NoRoute(func(c *gin.Context) {
        c.Redirect(http.StatusFound, "https://en.wikipedia.org/wiki/Special:Random")
    })

	router.Run() // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}
