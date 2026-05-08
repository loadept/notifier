package webhook

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
)

var ErrDiscordNoResponse = errors.New("discord has not responded to the request")

func NewHandler(httpClient *http.Client, webhookURL string) *WebHook {
	return &WebHook{httpClient: httpClient, webhookURL: webhookURL}
}

func (wh *WebHook) Handler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var body RequestBody
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			http.Error(w, "internal error", http.StatusInternalServerError)
			return
		}

		builder, ok := getBuilder(body.Source, body.Event)
		if !ok {
			http.Error(w, "unknown event", http.StatusBadRequest)
			return
		}

		embed, err := builder(body)
		if err != nil {
			http.Error(w, "internal error", http.StatusInternalServerError)
			return
		}

		if err := wh.sendToDiscord(embed); err != nil {
			// if errors.Is(err, ErrDiscordNoResponse) {
			// 	http.Error(w, err.Error(), http.StatusInternalServerError)
			// 	return
			// }
			http.Error(w, "internal error", http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusNoContent)
	})
}

func (wh *WebHook) sendToDiscord(embed DiscordEmbed) error {
	jsonData, err := json.Marshal(embed)
	if err != nil {
		return err
	}

	req, err := http.NewRequest("POST", wh.webhookURL, bytes.NewBuffer(jsonData))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := wh.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode > 299 {
		return ErrDiscordNoResponse
	}
	return nil
}
