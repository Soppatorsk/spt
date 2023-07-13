package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/Soppatorsk/spt/collage"
	"github.com/zmb3/spotify/v2"
	spotifyauth "github.com/zmb3/spotify/v2/auth"
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
	state        = "abc123" //TODO generate
)

type Collage struct {
	playlistID string
	tmpDir     string //TODO setting tmpdir in the function mixes up goroutines?
	client     *spotify.Client
}

func (c *Collage) GenerateCollage() string {
	img := collage.GenerateCollage(c.playlistID, c.tmpDir, c.client)
	return img
}

func main() {

	// first start an HTTP server
	url := auth.AuthURL(state)
	fmt.Println(url)

	http.HandleFunc("/callback", completeAuth)

	go func() {
		err := http.ListenAndServe(port, nil)
		if err != nil {
			log.Fatal(err)
		}
	}()

	// wait for auth to complete
	client := <-ch

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "POST" {
			r.ParseForm()
			playlistID := trimPlaylistLink(r.Form.Get("playlistID"))
			fmt.Println(playlistID)

			// Create a collage object
			c := &Collage{
				playlistID: playlistID,
				tmpDir:     "tmp/" + playlistID,
				client:     client,
			}
			go func() {
				imgURL := c.GenerateCollage()
				fmt.Println(imgURL)
			}()
		}

		http.ServeFile(w, r, "web/index.html")

		log.Println("Got request for:", r.URL.String())
	})

	//TODO input
	//10k list
	//yourPlaylist := "6qaVfh57zV2Y23B139X1Tn"
	//yourPlaylist := "6ko0RCsHny1iOJSF5hbmQ7"
	//small list
	/*
		yourPlaylist := "5SzZRpqqpxxhpURIDgiPyZ"

			go func() {

				img := collage.GenerateCollage(yourPlaylist, client)
				hostname, err := os.Hostname()
				if err != nil {
					fmt.Println(err)
				}
				fmt.Println(hostname + "/" + img)
			}()

	*/
	fmt.Println("im out lol")
	var input string
	fmt.Scanln(&input)
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

	fmt.Fprintf(w, "Login Completed!")

	ch <- client
}

func trimPlaylistLink(link string) string {
	// Find the index of "/playlist/"
	startIndex := strings.Index(link, "/playlist/") + len("/playlist/")
	if startIndex == -1 || startIndex >= len(link) {
		// Invalid playlist link, return empty string or error handle as per your requirement
		return ""
	}

	// Find the index of "?si="
	endIndex := strings.Index(link[startIndex:], "?si=")
	if endIndex == -1 {
		// If "?si=" is not found, use the end of the link
		endIndex = len(link)
	} else {
		// Add startIndex to the endIndex because of the substring indexing
		endIndex += startIndex
	}

	// Extract the playlist ID substring
	playlistID := link[startIndex:endIndex]
	return playlistID
}
