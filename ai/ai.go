package ai

import (
	"bytes"
	"context"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/zmb3/spotify/v2"
)

type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type Response struct {
	Choices []struct {
		Message struct {
			Content string `json:"content"`
		} `json:"message"`
	} `json:"choices"`
}

func GenerateResponse(id string, client *spotify.Client) string {
	url := "https://api.openai.com/v1/chat/completions"
	apiKey := os.Getenv("OPENAI_API_KEY")

	instructions := "Your task is to roast the users spotify playlist, tease and make fun of user and draw far-fetched conclusions on their music taste. User will provide a playlist title and some sample artists"
	playlistInfo := "TITLE: {TITLE} \nARTISTS: {ARTISTS}"

	p, err := client.GetPlaylist(context.Background(), spotify.ID(id))
	if err != nil {
		log.Println(err)
	}
	playlistInfo = strings.Replace(playlistInfo, "{TITLE}", p.Name, 1)

	pi, err := client.GetPlaylistItems(context.Background(), spotify.ID(id))
	if err != nil {
		log.Println(err)
	}

	artists := ""

	items := pi.Items
	for _, item := range items {
		artists = artists + item.Track.Track.Artists[0].Name + ",\n"
	}

	playlistInfo = strings.Replace(playlistInfo, "{ARTISTS}", artists, 1)

	messages := []Message{
		{Role: "system", Content: instructions},
		{Role: "user", Content: playlistInfo},
	}

	params := map[string]interface{}{
		"messages":    messages,
		"model":       "gpt-3.5-turbo",
		"max_tokens":  500,
		"temperature": 0.9,
	}

	jsonParams, _ := json.Marshal(params)

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonParams))
	if err != nil {
		log.Fatal(err)
	}

	req.Header.Set("Authorization", "Bearer "+apiKey)
	req.Header.Set("Content-Type", "application/json")

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Fatal(err)
	}

	var respJSON Response
	errk := json.Unmarshal([]byte(body), &respJSON)
	if errk != nil {
		log.Fatal(errk)
	}
	return string(respJSON.Choices[0].Message.Content)
}
