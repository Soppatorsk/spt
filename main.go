package main

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/Soppatorsk/spt/collage"
	"github.com/gin-contrib/static"
	"github.com/gin-gonic/gin"
	"github.com/zmb3/spotify/v2"
	spotifyauth "github.com/zmb3/spotify/v2/auth"
	"golang.org/x/oauth2"
)

const (
	port        = ":3000"
	redirectURI = "http://localhost:3000/callback"
)

var (
	clientID     = os.Getenv("SPOTIFY_ID")
	clientSecret = os.Getenv("SPOTIFY_SECRET")
	auth         = spotifyauth.New(spotifyauth.WithRedirectURL(redirectURI), spotifyauth.WithScopes(spotifyauth.ScopeUserReadPrivate))
	ch           = make(chan *spotify.Client)
	tk           = make(chan *oauth2.Token)
	state        = "abc123" //TODO generate
)

type RequestData struct {
	ID    string `json:"id"`
	Token string `json:"token"`
}

type playlist struct {
	ID         string `json:"id"`
	CollageURL string `json:"collageURL"`
}

/*
	var playlists = []playlist{
		{
			ID:         "0",
			CollageURL: "/img/test.jpg",
		},
		{
			ID:         "3I3qR4YotZyh9R5BFSUx89",
			CollageURL: "/img/3I3qR4YotZyh9R5BFSUx89.jpg",
		},
		{
			ID:         "3okg2NywBjVFkFM9LNrWA2",
			CollageURL: "/img/3okg2NywBjVFkFM9LNrWA2.jpg",
		},
		{
			ID:         "2oLuXBMWzDQjbGNPThml5C",
			CollageURL: "/img/2oLuXBMWzDQjbGNPThml5C.jpg",
		},
	}
*/
//TODO database?
var playlists = []playlist{}

func main() {

	router := gin.Default()
	//static
	router.Use(static.Serve("/", static.LocalFile("./vue-front/dist", true)))
	//API
	router.GET("/callback", completeAuth)
	router.GET("/playlists/", getPlaylists)
	router.GET("/playlists/:id", getPlaylistById)
	router.GET("/img/:filename", getImg)
	router.GET("/auth/", getAuth)
	router.POST("/playlists/", createPlaylist)
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

func getAuth(c *gin.Context) {
	url := auth.AuthURL(state)
	// c.Redirect(http.StatusTemporaryRedirect, url)
	c.String(http.StatusOK, url)
}

func getImg(c *gin.Context) {
	id := c.Param("filename")
	imgPath := filepath.Join("img", id)
	c.File(imgPath)
}

// POST
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

	ck, err := c.Request.Cookie("token")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("server side ck" + ck.Value)

	fmt.Println(id)
	fmt.Println(token)

	// id := "2oLuXBMWzDQjbGNPThml5C"

	oauthToken := &oauth2.Token{
		AccessToken: token,
	}
	//todo input validation
	// create a collage object
	client := spotify.New(auth.Client(c, oauthToken))
	imgurl := collage.GenerateCollage(id, client)

	var newPlaylist = playlist{
		ID: id, CollageURL: imgurl,
	}
	playlists = append(playlists, newPlaylist)
	c.IndentedJSON(http.StatusCreated, newPlaylist)
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
		log.Fatal(err)
	}
	if st := c.Request.FormValue("state"); st != state {
		log.Fatalf("State mismatch: %s != %s\n", st, state)
	}

	// use the token to get an authenticated client
	// client := spotify.New(auth.Client(c.Request.Context(), tok))
	// print(client)

	// fmt.Println(tok.AccessToken)
	setTokenInCookie(c, tok.AccessToken)
	c.Redirect(http.StatusTemporaryRedirect, "/")
}

func setTokenInCookie(c *gin.Context, token string) {
	cookie := &http.Cookie{
		Name:     "token",
		Value:    token,
		Path:     "/",
		HttpOnly: false,
	}
	http.SetCookie(c.Writer, cookie)
}
