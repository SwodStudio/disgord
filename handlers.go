package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)

type IncomingPayload struct {
	Source  string `json:"source" binding:"required"`
	Content string `json:"content" binding:"required"`
}

type DiscordPayload struct {
	Username string `json:"username"`
	Content  string `json:"content"`
}

func HandleSend(c *gin.Context) {
	var input IncomingPayload

	// check json format
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid format"})
		return
	}

	payload := DiscordPayload{
		Username: fmt.Sprintf("DisGord-%s", input.Source),
		Content:  input.Content,
	}

	// forward to dc
	if err := forwardToDiscord(payload); err != nil {
		log.Printf("Forward error: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Forward failed"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "Sent"})
}

func forwardToDiscord(payload DiscordPayload) error {
	url := os.Getenv("DISCORD_WEBHOOK")
	if url == "" {
		return fmt.Errorf("missing DISCORD_WEBHOOK")
	}

	b, _ := json.Marshal(payload)
	resp, err := http.Post(url, "application/json", bytes.NewBuffer(b))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// rate limits and errors
	if resp.StatusCode >= 400 {
		return fmt.Errorf("discord api error status: %d", resp.StatusCode)
	}

	return nil
}
