package main

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/Soppatorsk/spt/collage"
	"github.com/gin-gonic/gin"
	"github.com/zmb3/spotify/v2"
	spotifyauth "github.com/zmb3/spotify/v2/auth"
	"golang.org/x/oauth2"
)

const (
	port        = ":3001"
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

var playlists = []playlist{
	{ID: "0", CollageURL: "img/test.jpg"},
	{ID: "3I3qR4YotZyh9R5BFSUx89", CollageURL: "img/3I3qR4YotZyh9R5BFSUx89.jpg"},
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
	c.String(http.StatusOK, url)

	token := <-tk
	c.String(http.StatusOK, token.AccessToken)
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
	fmt.Println(id)
	fmt.Println(token)

	oauthToken := &oauth2.Token{
		AccessToken: token,
	}
	//playlistid := "3okg2NywBjVFkFM9LNrWA2"
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

func main() {

	http.HandleFunc("/callback", completeAuth)

	go func() {
		err := http.ListenAndServe(":3000", nil)
		if err != nil {
			log.Fatal(err)
		}
	}()
	router := gin.Default()
	router.GET("/playlists", getPlaylists)
	router.GET("/playlists/:id", getPlaylistById)
	router.GET("/auth/", getAuth)
	router.POST("/playlists/", createPlaylist)
	router.Run("localhost" + port)
}

func completeAuth(w http.ResponseWriter, r *http.Request) {

	tok, err := auth.Token(r.Context(), state, r)
	if err != nil {
		http.Error(w, "Couldn't get token", http.StatusForbidden)
		log.Fatal(err)
	}
	if st := r.FormValue("state"); st != state {
		http.NotFound(w, r)
		log.Fatalf("State mismatch: %s != %s\n", st, state)
	}

	// use the token to get an authenticated client
	client := spotify.New(auth.Client(r.Context(), tok))

	fmt.Fprintf(w, " \nLogin Completed! You can close this window")
	fmt.Fprintf(w, "\n"+tok.AccessToken)
	tk <- tok
	ch <- client
}
