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
	Source   string  `json:"source" binding:"required"`
	Content  string  `json:"content" binding:"required"`
	IconURL  string  `json:"icon_url"`
	Nickname string  `json:"nickname"`
	Embeds   []Embed `json:"embeds"`
}

type DiscordPayload struct {
	Username  string  `json:"username,omitempty"`
	AvatarURL string  `json:"avatar_url,omitempty"`
	Content   string  `json:"content,omitempty"`
	Embeds    []Embed `json:"embeds,omitempty"`
}

// Embed
type Embed struct {
	Title       string          `json:"title,omitempty"`
	Description string          `json:"description,omitempty"`
	URL         string          `json:"url,omitempty"`
	Color       int             `json:"color,omitempty"`
	Footer      *EmbedFooter    `json:"footer,omitempty"`
	Image       *EmbedImage     `json:"image,omitempty"`
	Thumbnail   *EmbedThumbnail `json:"thumbnail,omitempty"`
	Author      *EmbedAuthor    `json:"author,omitempty"`
	Fields      []EmbedField    `json:"fields,omitempty"`
}

type EmbedField struct {
	Name   string `json:"name,omitempty"`
	Value  string `json:"value,omitempty"`
	Inline bool   `json:"inline,omitempty"`
}

type EmbedAuthor struct {
	Name    string `json:"name,omitempty"`
	URL     string `json:"url,omitempty"`
	IconURL string `json:"icon_url,omitempty"`
}

type EmbedThumbnail struct {
	URL string `json:"url,omitempty"`
}

type EmbedImage struct {
	URL string `json:"url,omitempty"`
}

type EmbedFooter struct {
	Text    string `json:"text,omitempty"`
	IconURL string `json:"icon_url,omitempty"`
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
		Embeds:    input.Embeds,
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
		return fmt.Errorf("DISCORD_WEBHOOK not found")
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

