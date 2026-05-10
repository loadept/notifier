// Package webhook proporciona logica de manejo de webhooks de notificaciones
package webhook

import (
	"encoding/json"
	"net/http"
)

type WebHook struct {
	httpClient *http.Client
	webhookURL string
}

type RequestBody struct {
	Source string          `json:"source"`
	Event  string          `json:"event"`
	Status string          `json:"status,omitempty"`
	Meta   json.RawMessage `json:"meta,omitempty"`
}

type DiscordEmbed struct {
	Username  string  `json:"username"`
	AvatarURL string  `json:"avatar_url"`
	Embeds    []Embed `json:"embeds"`
}

type Embed struct {
	Title       string       `json:"title"`
	Description string       `json:"description,omitempty"`
	Color       Color        `json:"color"`
	Fields      []EmbedField `json:"fields,omitempty"`
}

type EmbedField struct {
	Name   string `json:"name"`
	Value  string `json:"value"`
	Inline bool   `json:"inline"`
}

type Color int

const (
	ColorGithub  Color = 0x6E40C9 // #6E40C9
	ColorIpseek  Color = 0x3498DB // #3498DB
	ColorLoadept Color = 0xE67E22 // #E67E22

	ColorSuccess Color = 0x2ECC71 // #2ECC71
	ColorError   Color = 0xF5273F // #F5273F
)
