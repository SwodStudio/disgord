# disgord

`disgord` is a Go-based API and SDK for forwarding messages to Discord Webhooks 

## Usage

### Get

```bash
go get github.com/SwodStudio/disgord
```

### Environment Variables

You need set `DISCORD_WEBHOOK` environment variable.

### Running as an API Server

```go
package main

import "github.com/SwodStudio/disgord"

func main() {
    disgord.RunServer()
}
```

### Data Structures

#### `DiscordPayload`

```go
type DiscordPayload struct {
	Username  string  `json:"username,omitempty"`
	AvatarURL string  `json:"avatar_url,omitempty"`
	Content   string  `json:"content,omitempty"`
	Embeds    []Embed `json:"embeds,omitempty"`
}
```

#### `Embed Example`

```go
embed := disgord.Embed{
    Title:       "Embed Test",
    Description: "Embed Description :3",
    Color:       0xE74C3C,
    Fields: []disgord.EmbedField{
        {Name: "Field Test", Value: "Value", Inline: true},
    },
    Footer: &disgord.EmbedFooter{
        Text: "Embed Footer",
    },
}

payload := disgord.DiscordPayload{
    Content: "Testing message",
    Embeds:  []disgord.Embed{embed},
}

disgord.Send(payload)
```

```json
{
  "content": "Testing message",
  "embeds": [
    {
      "title": "Embed Test",
      "description": "Embed Description :3",
      "color": 15158332,
      "fields": [
        {
          "name": "Field Test",
          "value": "Value",
          "inline": true
        }
      ],
      "footer": {
        "text": "Embed Footer"
      }
    }
  ]
}
```

```
