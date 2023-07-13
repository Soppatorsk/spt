package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"text/template"

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
	fmt.Println("Running...")
	// first start an HTTP server
	/*
		url := auth.AuthURL(state)
		fmt.Println(url)
	*/

	http.HandleFunc("/callback", completeAuth)
	http.HandleFunc("/test", test)
	http.HandleFunc("/img/", serveImage)

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "web/index.html")

		log.Println("Got request for:", r.URL.String())
	})
	go func() {
		err := http.ListenAndServe(port, nil)
		if err != nil {
			log.Fatal(err)
		}
	}()

	//client := <-ch
	//fmt.Println(client)
	var input string
	fmt.Scanln(&input)
}

func serveImage(w http.ResponseWriter, r *http.Request) {

	imagePath := "./img/" + r.URL.Path[len("/img/"):]
	fmt.Println(imagePath)
	http.ServeFile(w, r, imagePath)
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

	fmt.Fprintf(w, "Login Completed! You can close this window")

	ch <- client
}

func test(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		client := <-ch
		r.ParseForm()
		playlistID := trimPlaylistLink(r.Form.Get("playlistID"))
		//TODO input validation
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
	url := auth.AuthURL(state)

	tmpl, err := template.ParseFiles("web/test.html")
	if err != nil {
		log.Fatal(err)
	}

	imageDir := "img"
	imageFiles, err := filepath.Glob(filepath.Join(imageDir, "*.jpg")) // Change the pattern if needed
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Define the data to pass to the template
	data := struct {
		URL    string
		Images []string
	}{
		URL:    url,
		Images: imageFiles,
	}

	// Render the template with the data
	err = tmpl.Execute(w, data)
	if err != nil {
		log.Fatal(err)
	}

}
func trimPlaylistLink(link string) string {
	fmt.Println("trim")
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
