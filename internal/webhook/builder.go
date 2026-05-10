package webhook

import (
	"encoding/json"
	"fmt"
)

type EmbedBuilder func(r RequestBody) (DiscordEmbed, error)

var builders = map[string]EmbedBuilder{
	"github:deploy":  buildGithubDeploy,
	"ipseek:health":  buildIpseekHealth,
	"loadept:visit":  buildLoadeptVisit,
	"loadept:health": buildLoadeptHealth,
}

const avatar = "https://assets.loadept.com/loggy.webp"

func getBuilder(source, event string) (EmbedBuilder, bool) {
	key := fmt.Sprintf("%s:%s", source, event)
	builder, ok := builders[key]
	return builder, ok
}

func buildGithubDeploy(r RequestBody) (DiscordEmbed, error) {
	var meta struct {
		Repo    string `json:"repo"`
		Commit  string `json:"commit"`
		Branch  string `json:"branch"`
		Actor   string `json:"actor"`
		Message string `json:"message"`
	}
	if err := json.Unmarshal(r.Meta, &meta); err != nil {
		return DiscordEmbed{}, err
	}

	color := ColorGithub
	title := "✅ Deploy exitoso"
	if r.Status != "success" {
		color = ColorError
		title = "❌ Deploy fallido"
	}

	return DiscordEmbed{
		Username:  "Loadept Notifier",
		AvatarURL: avatar,
		Embeds: []Embed{
			{
				Title: title,
				Color: color,
				Fields: []EmbedField{
					{Name: "Repo", Value: meta.Repo, Inline: true},
					{Name: "Branch", Value: meta.Branch, Inline: true},
					{Name: "Actor", Value: meta.Actor, Inline: true},
					{Name: "Commit", Value: fmt.Sprintf("`%s`", meta.Commit[:7]), Inline: true},
					{Name: "Mensaje", Value: meta.Message, Inline: false},
				},
			},
		},
	}, nil
}

func buildIpseekHealth(r RequestBody) (DiscordEmbed, error) {
	if r.Status == "ok" {
		return DiscordEmbed{
			Username:  "Loadept Notifier",
			AvatarURL: avatar,
			Embeds: []Embed{
				{
					Title: "✅ IP Seek — online",
					Color: ColorSuccess,
				},
			},
		}, nil
	}

	var meta struct {
		Reason   string `json:"reason"`
		Endpoint string `json:"endpoint"`
	}
	json.Unmarshal(r.Meta, &meta)

	return DiscordEmbed{
		Username:  "Loadept Notifier",
		AvatarURL: avatar,
		Embeds: []Embed{
			{
				Title: "❌ IP Seek — caído",
				Color: ColorError,
				Fields: []EmbedField{
					{Name: "Endpoint", Value: meta.Endpoint, Inline: true},
					{Name: "Razón", Value: meta.Reason, Inline: true},
				},
			},
		},
	}, nil
}

func buildLoadeptVisit(r RequestBody) (DiscordEmbed, error) {
	var meta struct {
		Path    string `json:"path"`
		IP      string `json:"ip"`
		Country string `json:"country"`
		City    string `json:"city"`
	}
	json.Unmarshal(r.Meta, &meta)

	return DiscordEmbed{
		Username:  "Notifier",
		AvatarURL: avatar,
		Embeds: []Embed{
			{
				Title:       "👤 Nueva visita",
				Description: "¡Alguien ha visitado el short url de loadept!",
				Color:       ColorLoadept,
				Fields: []EmbedField{
					{Name: "IP", Value: meta.IP, Inline: true},
					{Name: "Short URL", Value: meta.Path, Inline: true},
					{Name: "País", Value: meta.Country, Inline: true},
					{Name: "Ciudad", Value: meta.City, Inline: true},
				},
			},
		},
	}, nil
}

func buildLoadeptHealth(r RequestBody) (DiscordEmbed, error) {
	if r.Status == "ok" {
		return DiscordEmbed{
			Username:  "Loadept Notifier",
			AvatarURL: avatar,
			Embeds: []Embed{
				{
					Title: "✅ loadept.com — online",
					Color: ColorSuccess,
				},
			},
		}, nil
	}

	var meta struct {
		Reason   string `json:"reason"`
		Endpoint string `json:"endpoint"`
	}
	json.Unmarshal(r.Meta, &meta)

	return DiscordEmbed{
		Username:  "Loadept Notifier",
		AvatarURL: avatar,
		Embeds: []Embed{
			{
				Title: "❌ loadept.com — caído",
				Color: ColorError,
				Fields: []EmbedField{
					{Name: "Endpoint", Value: meta.Endpoint, Inline: true},
					{Name: "Razón", Value: meta.Reason, Inline: true},
				},
			},
		},
	}, nil
}
