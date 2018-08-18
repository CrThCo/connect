package main

import (
	"fmt"
	"os"
	"strconv"

	"github.com/gin-gonic/gin"
)

func main() {
	gin.DisableConsoleColor()
	r := gin.Default()
	r.Use(gin.Logger())

	routes(r)

	port, err := strconv.Atoi(os.Getenv("PORT"))
	if err != nil {
		port = 5050
	}
	addr := fmt.Sprintf(":%d", port)
	r.Run(addr)
}
