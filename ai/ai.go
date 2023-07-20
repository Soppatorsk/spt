package ai

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
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

	instructions := "Your task is to roast the users spotify playlist, tease and make fun of user and draw far-fetched conclusions on their music taste. User will provide a playlist title and some samples song - artist"
	//instructions := "Mock and make fun of the user, don't hold back and absolutely roast the user. Feel free to be rude and condecending and belittle user"
	// instructions := "Make a character portrayal of the user based on their playlist, be mean and judgemental and speak nothing but the most brutal truth. Talk about the user and not about the music in itself. Make guesses and far-fetched deductions based on the available information"
	playlistInfo := "TITLE: {TITLE} \nUSER: {USER} \nSAMPLES: {SAMPLES}"

	p, err := client.GetPlaylist(context.Background(), spotify.ID(id))
	if err != nil {
		log.Println(err)
	}
	playlistInfo = strings.Replace(playlistInfo, "{TITLE}", p.Name, 1)
	playlistInfo = strings.Replace(playlistInfo, "{USER}", p.Owner.DisplayName, 1)

	pi, err := client.GetPlaylistItems(context.Background(), spotify.ID(id))
	if err != nil {
		log.Println(err)
	}

	samples := ""

	items := pi.Items
	for _, item := range items {
		samples = samples + item.Track.Track.Artists[0].Name + " - " + item.Track.Track.Name + ",\n"
	}

	fmt.Println(samples)
	playlistInfo = strings.Replace(playlistInfo, "{SAMPLES}", samples, 1)

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
