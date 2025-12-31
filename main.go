package disgord

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func RunServer() {
	_ = godotenv.Load()

	r := gin.Default()
	r.POST("/send", HandleSend)

	if err := r.Run(":8080"); err != nil {
		log.Fatal(err)
	}
}

