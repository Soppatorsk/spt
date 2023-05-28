package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"

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

func main() {
	// first start an HTTP server
	http.HandleFunc("/callback", completeAuth)
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		log.Println("Got request for:", r.URL.String())
	})
	go func() {
		err := http.ListenAndServe(port, nil)
		if err != nil {
			log.Fatal(err)
		}
	}()

	url := auth.AuthURL(state)
	fmt.Println(url)

	// wait for auth to complete
	client := <-ch

	// use the client to make calls that require authorization
	user, err := client.CurrentUser(context.Background())
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("You are logged in as:", user.ID)
	fmt.Println("You are logged in as:", user.Followers.Count)

	testPlaylist := "4H972R0YIOTGJdWOGLiPJJ" //TODO input
	testTrack := "6IwSGBUjtiPsDiXPR0yTSS"    //TODO input
	playlist, err := client.GetPlaylistTracks(context.Background(), spotify.ID(testPlaylist))
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(playlist.Endpoint)

	track, err := client.GetTrack(context.Background(), spotify.ID(testTrack))
	img := track.Album.Images[0]
	fmt.Println(img.URL)

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

func printImg(w http.ResponseWriter, r *http.Request, img spotify.Image) {

	fmt.Fprintf(w, "<img src=\""+img.URL+"\">")
}

/*
func main() {
	//Environment variables

	port := ":3000"

	fmt.Println(client_id + "\n" + client_secret + "\n" + "port " + port)
	http.HandleFunc("/", helloHandler)
	http.ListenAndServe(port, nil)
}

func helloHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Hello, World!");
}

*/
