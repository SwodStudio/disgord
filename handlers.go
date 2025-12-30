package disgord

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
	Source   string `json:"source" binding:"required"`
	Content  string `json:"content" binding:"required"`
	IconURL  string `json:"icon_url"`
	Nickname string `json:"nickname"`
}

type DiscordPayload struct {
	Username  string `json:"username"`
	AvatarURL string `json:"avatar_url"`
	Content   string `json:"content"`
}

func HandleSend(c *gin.Context) {
	var input IncomingPayload
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid format"})
		return
	}

	name := input.Nickname
	if name == "" {
		name = fmt.Sprintf("DisGord-%s", input.Source)
	}

	payload := DiscordPayload{
		Username:  name,
		AvatarURL: input.IconURL,
		Content:   input.Content,
	}

	if err := Send(payload); err != nil {
		log.Printf("Forward error: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Forward failed"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "Sent"})
}

func Send(payload DiscordPayload) error {
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

	if resp.StatusCode >= 400 {
		return fmt.Errorf("discord api error: %d", resp.StatusCode)
	}

	return nil
}
