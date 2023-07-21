package color

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/zmb3/spotify/v2"
)

/*
energy R
dance	G
valence	B
*/

type song struct {
	Title string `json:"title"`
	Color string `json:"color"`
}

var songs = []song{}

func Generate(id string, client *spotify.Client) string {
	songs = []song{}
	p, err := client.GetPlaylistItems(context.Background(), spotify.ID(id))
	if err != nil {
		log.Println(err)
	}
	items := p.Items

	for i, item := range items {
		if i >= 50 {
			break
		}
		songColor(string(item.Track.Track.ID), client)
	}
	jsonData, err := json.Marshal(songs)
	if err != nil {
		log.Println(err)
	}
	fmt.Println("///Playlist color: " + PlaylistColor(id, client))
	return string(jsonData)
}

func songColor(id string, client *spotify.Client) {
	s, err := client.GetTrack(context.Background(), spotify.ID(id))
	if err == nil { //local file? skip

		t, err := client.GetAudioFeatures(context.Background(), spotify.ID(id))
		if err != nil {
			log.Println(err)
		} else {
			fmt.Println(t[0].Energy)
			r := int(t[0].Energy * 255)
			g := int(t[0].Danceability * 255)
			//b := int(t[0].Valence * 255)
			b := int(255 - (t[0].Valence * 255)) //Inverse. the blues. get it? haha
			hexCode := fmt.Sprintf("#%02x%02x%02x", r, g, b)
			fmt.Println(hexCode)
			var newSong = song{
				Title: s.Artists[0].Name + " - " + s.Name,
				Color: hexCode,
			}
			songs = append(songs, newSong)
		}
	}
}

func PlaylistColor(id string, client *spotify.Client) string {
	p, err := client.GetPlaylistItems(context.Background(), spotify.ID(id))
	if err != nil {
		log.Println(err)
	}
	r := 0
	g := 0
	b := 0

	items := p.Items
	iterated := 0
	for i, item := range items {
		f, err := client.GetAudioFeatures(context.Background(), item.Track.Track.ID)
		if err == nil {
			if f != nil {
				r = r + int(f[0].Energy*255)
				g = g + int(f[0].Danceability*255)
				b = b + int(255-(f[0].Valence*255))
				i++
				iterated = i
			}
		}
	}
	fmt.Println(iterated)
	r = r / iterated
	g = g / iterated
	b = b / iterated
	hexCode := fmt.Sprintf("#%02x%02x%02x", r, g, b)
	return hexCode
}
