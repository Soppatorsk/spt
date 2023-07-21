package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"

	"github.com/Soppatorsk/spt/ai"
	"github.com/Soppatorsk/spt/collage"
	"github.com/Soppatorsk/spt/color"
	"github.com/gin-contrib/static"
	"github.com/gin-gonic/gin"
	"github.com/zmb3/spotify/v2"
	spotifyauth "github.com/zmb3/spotify/v2/auth"
	"golang.org/x/oauth2"
)

const (
	hostDir = "" //ex. torsk.net/spt/ /spt as root //TODO global config file
	port    = ":3000"
	//redirectURI = "https://torsk.net/spt/callback"
	redirectURI = "http://localhost:3000" + hostDir + "/callback" //EDIT
)

var (
	auth = spotifyauth.New(spotifyauth.WithRedirectURL(redirectURI), spotifyauth.WithScopes(spotifyauth.ScopeUserReadPrivate))
)

type playlist struct {
	ID         string `json:"id"`
	Name       string `json:"name"`
	User       string `json:"user"`
	CollageURL string `json:"collageURL"`
	AI         string `json:"ai"`
	Color      string `json:"color"`
}

// TODO database?
var playlists = []playlist{}

func main() {
	loadJSON()
	router := gin.Default()
	//static/frontend
	router.Use(static.Serve(hostDir+"/", static.LocalFile("./vue-front/dist", true)))

	router.GET(hostDir+"/test/:id", test)
	//Get all public lists
	router.GET(hostDir+"/playlists/", getPlaylists)
	//Get image
	router.GET(hostDir+"/img/:filename", getImg)

	//Save to db.json
	router.GET(hostDir+"/save", saveJSON)

	//Get/Create full list and save to db
	router.GET(hostDir+"/playlist/:id", createPlaylist)
	//Get/Create response only
	router.GET(hostDir+"/ai/:id", getAiResponse)
	//Get/create collage only
	router.GET("/collage/:id", getCollage)

	router.Run("localhost" + port)
}

func test(c *gin.Context) {
	id := c.Param("id")
	client := getClient()
	s := color.Generate(id, client)
	c.String(http.StatusOK, s)
}

// GET
func getPlaylist(c *gin.Context) {
	id := c.Param("id")
	playlist, err := findPlaylistById(id)
	if err != nil {
		//does not exist
		log.Println(err)
		return
	}
	c.IndentedJSON(http.StatusOK, playlist)
}

func getPlaylists(c *gin.Context) {
	c.IndentedJSON(http.StatusOK, playlists)
}

func getImg(c *gin.Context) {
	id := c.Param("filename")
	imgPath := filepath.Join("img", id)
	c.File(imgPath)
}

func getCollage(c *gin.Context) {
	id := c.Param("id")
	//validate
	client := getClient()
	imgPath := collage.GenerateCollage(id, client)
	c.File(imgPath)
}

func saveJSON(c *gin.Context) {
	jsonData, err := json.Marshal(playlists)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to marshal JSON"})
		return
	}
	err = ioutil.WriteFile("db.json", jsonData, 0644)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to write JSNO file"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "JSON file saved."})
}

func loadJSON() {
	data, err := ioutil.ReadFile("db.json")
	if err != nil {
		log.Println(err)
	}

	var p []playlist
	err = json.Unmarshal(data, &p)
	if err != nil {
		log.Println(err)
	}

	playlists = p

}

func getAiResponse(c *gin.Context) {
	id := c.Param("id")
	client := getClient()

	c.String(http.StatusOK, ai.GenerateResponse(id, client))
}

func createPlaylist(c *gin.Context) {
	id := c.Param("id")
	p, err := findPlaylistById(id)
	if err != nil {
		log.Println(err)
	}
	if p == nil { //no entry in db
		client := getClient()

		imgurl := collage.GenerateCollage(id, client)
		p, err := client.GetPlaylist(context.Background(), spotify.ID(id))
		if err != nil {
			log.Println(err)
		}
		ai := ai.GenerateResponse(id, client)
		color := color.PlaylistColor(id, client)

		var newPlaylist = playlist{
			ID:         id,
			Name:       p.Name,
			User:       p.Owner.DisplayName,
			CollageURL: imgurl,
			AI:         ai,
			Color:      color,
		}
		playlists = append([]playlist{newPlaylist}, playlists...)
		saveJSON(c)
		c.IndentedJSON(http.StatusCreated, newPlaylist)
	} else { //entry in db, read it
		getPlaylist(c)
	}
}

// POST/GET Helpers
func findPlaylistById(id string) (*playlist, error) {
	for i, p := range playlists {
		if p.ID == id {
			return &playlists[i], nil
		}
	}
	return nil, errors.New("playlist not found")
}

func spotifyAuth() (string, error) {
	//base64 encoded client_id:client_secret
	spotifyKey := os.Getenv("SPOTIFY_KEY")
	data := url.Values{}
	data.Set("grant_type", "client_credentials")
	payload := strings.NewReader(data.Encode())
	req, err := http.NewRequest("POST", "https://accounts.spotify.com/api/token", payload)
	if err != nil {
		return "", err
	}
	authHeader := fmt.Sprintf("Basic %s", spotifyKey)
	req.Header.Set("Authorization", authHeader)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	// Read the response body
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	// Parse the JSON response
	var response struct {
		AccessToken string `json:"access_token"`
	}
	err = json.Unmarshal(body, &response)
	if err != nil {
		return "", err
	}

	return response.AccessToken, nil
}

func getClient() *spotify.Client {
	token, err := spotifyAuth()
	if err != nil {
		log.Println(err)
	}

	oauthToken := &oauth2.Token{
		AccessToken: token,
	}

	return spotify.New(auth.Client(context.Background(), oauthToken))

}
