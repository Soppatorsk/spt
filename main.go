package main

import (
	"context"
	"crypto/rand"
	"encoding/base64"
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
	clientID     = os.Getenv("SPOTIFY_ID")
	clientSecret = os.Getenv("SPOTIFY_SECRET")
	auth         = spotifyauth.New(spotifyauth.WithRedirectURL(redirectURI), spotifyauth.WithScopes(spotifyauth.ScopeUserReadPrivate))
	ch           = make(chan *spotify.Client)
	tk           = make(chan *oauth2.Token)
	state, err   = generateState()
)

type RequestData struct {
	ID    string `json:"id"`
	Token string `json:"token"`
}

type playlist struct {
	ID         string `json:"id"`
	Name       string `json:"name"`
	User       string `json:"user"`
	CollageURL string `json:"collageURL"`
	AI         string `json:"ai"`
}

// TODO database?
var playlists = []playlist{}

func main() {
	loadJSON()
	router := gin.Default()
	//static
	router.Use(static.Serve(hostDir+"/", static.LocalFile("./vue-front/dist", true)))
	//API
	router.GET(hostDir+"/callback", completeAuth)

	// router.GET(hostDir+"/auth/", getAuthOld)
	router.GET(hostDir+"/auth/", getAuth)
	router.GET(hostDir+"/playlists/", getPlaylists)
	router.GET(hostDir+"/playlists/:id", getPlaylistById)
	router.GET(hostDir+"/img/:filename", getImg)

	router.GET(hostDir+"/save", saveJSON)

	router.POST(hostDir+"/playlists/", createPlaylist)
	router.POST(hostDir+"/ai/", createAiResponse)

	router.Run("localhost" + port)
}

// GET
func getPlaylistById(c *gin.Context) {
	id := c.Param("id")
	playlist, err := findPlaylistById(id)
	if err != nil {
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": "Playlist not found"})
		return
	}
	c.IndentedJSON(http.StatusOK, playlist)
}

func getPlaylists(c *gin.Context) {
	c.IndentedJSON(http.StatusOK, playlists)
}

func getAuthOld(c *gin.Context) {
	url := auth.AuthURL(state)
	// c.Redirect(http.StatusTemporaryRedirect, url)
	c.String(http.StatusOK, url)
}
func getAuth(c *gin.Context) {
	token, err := newAuth()
	fmt.Print(err)
	setTokenInCookie(c, token)
}

func getImg(c *gin.Context) {
	id := c.Param("filename")
	imgPath := filepath.Join("img", id)
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

// POST
func createAiResponse(c *gin.Context) {
	var requestData RequestData

	if err := c.ShouldBindJSON(&requestData); err != nil {
		// Handle error
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request data"})
	}
	// Access the values from the request data
	id := requestData.ID
	token := requestData.Token

	oauthToken := &oauth2.Token{
		AccessToken: token,
	}
	//todo create a collage object
	client := spotify.New(auth.Client(c, oauthToken))
	c.String(http.StatusOK, ai.GenerateResponse(id, client))
}

func createPlaylist(c *gin.Context) {
	var requestData RequestData

	if err := c.ShouldBindJSON(&requestData); err != nil {
		// Handle error
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request data"})
		return
	}
	// Access the values from the request data
	id := requestData.ID
	token := requestData.Token

	p, err := findPlaylistById(id)
	if p == nil {

		fmt.Println(id)
		fmt.Println(token)

		// id := "2oLuXBMWzDQjbGNPThml5C"

		oauthToken := &oauth2.Token{
			AccessToken: token,
		}
		//todo create a collage object
		client := spotify.New(auth.Client(c, oauthToken))

		imgurl := collage.GenerateCollage(id, client)
		p, err := client.GetPlaylist(context.Background(), spotify.ID(id))
		if err != nil {
			log.Println(err)
		}
		ai := ai.GenerateResponse(id, client)

		var newPlaylist = playlist{
			ID: id, Name: p.Name, User: p.Owner.DisplayName, CollageURL: imgurl, AI: ai,
		}
		playlists = append([]playlist{newPlaylist}, playlists...)
		c.IndentedJSON(http.StatusCreated, newPlaylist)
	} else {
		c.IndentedJSON(http.StatusInternalServerError, "Playlist already generated")
		fmt.Println("Playlist already generated", err)
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

func completeAuth(c *gin.Context) {

	tok, err := auth.Token(c.Request.Context(), state, c.Request)
	if err != nil {
		log.Println(err)
	}
	if st := c.Request.FormValue("state"); st != state {
		log.Printf("State mismatch: %s != %s\n", st, state)
	}

	// use the token to get an authenticated client
	// client := spotify.New(auth.Client(c.Request.Context(), tok))
	// print(client)

	// fmt.Println(tok.AccessToken)
	setTokenInCookie(c, tok.AccessToken)
	c.Redirect(http.StatusTemporaryRedirect, hostDir+"/")
}

func newAuth() (string, error) {
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

func generateState() (string, error) {
	// Adjust the length of the state string as needed
	const stateLength = 32

	// Generate a random byte slice of the desired length
	stateBytes := make([]byte, stateLength)
	if _, err := rand.Read(stateBytes); err != nil {
		return "", err
	}

	// Encode the byte slice as a base64 string
	state := base64.URLEncoding.EncodeToString(stateBytes)

	return state, nil
}

func setTokenInCookie(c *gin.Context, token string) {
	cookie := &http.Cookie{
		Name:     "token",
		Value:    token,
		Path:     hostDir + "/",
		HttpOnly: false,
	}
	http.SetCookie(c.Writer, cookie)
}
