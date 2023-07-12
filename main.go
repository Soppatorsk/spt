package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/zmb3/spotify/v2"
	spotifyauth "github.com/zmb3/spotify/v2/auth"
)

const (
	port        = ":3000"
	redirectURI = "http://localhost:3000/callback"

	generatedDir = "collages"
)

var (
	tmpDir = "tmp"

	clientID     = os.Getenv("SPOTIFY_ID")
	clientSecret = os.Getenv("SPOTIFY_SECRET")
	auth         = spotifyauth.New(spotifyauth.WithRedirectURL(redirectURI), spotifyauth.WithScopes(spotifyauth.ScopeUserReadPrivate))
	ch           = make(chan *spotify.Client)
	state        = "abc123" //TODO generate
)

func main() {
	// first start an HTTP server
	url := auth.AuthURL(state)
	fmt.Println(url)

	http.HandleFunc("/callback", completeAuth)
	http.HandleFunc("/test", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "public/test.html")
		log.Println("Got request for:", r.URL.String())
	})

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		//http.ServeFile(w, r, "public/index.html")
		//TODO cant use %s???
		fmt.Fprint(w, "<a href=\"", url, "\">Click here</a>")
		log.Println("Got request for:", r.URL.String())
	})
	go func() {
		err := http.ListenAndServe(port, nil)
		if err != nil {
			log.Fatal(err)
		}
	}()

	// wait for auth to complete
	client := <-ch
	//TODO input
	//yourPlaylist := "6qaVfh57zV2Y23B139X1Tn"
	yourPlaylist := "5SzZRpqqpxxhpURIDgiPyZ"
	tmpDir = tmpDir + "/" + yourPlaylist
	os.Mkdir(tmpDir, 775)
	generateCollage(yourPlaylist, client)

}

func generateCollage(playlistID string, client *spotify.Client) {

	playlist, err := client.GetPlaylistItems(context.Background(), spotify.ID(playlistID))
	if err != nil {
		log.Fatal(err)
	}

	//get track list from playlist
	//TODO use offset big playlists
	items := playlist.Items
	for _, item := range items {
		downloadImage(item.Track.Track.Album.Images[0].URL)
	}

	files, err := filepath.Glob(tmpDir + "/*")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(append(files))
	cmd := exec.Command("montage", append(files, "-geometry", "256x256+0+0", generatedDir+"/"+playlistID+".jpg")...)
	err = cmd.Run()
	if err != nil {
		log.Fatal(err)
	}

	//TODO remove files in tmp. Cronjob?
}

func downloadImage(url string) error {
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	out, err := os.Create(tmpDir + "/" + url[25:] + ".jpg")
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, resp.Body)
	return err
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
	//http.ServeFile(w, r, "public/login_complete.html")
	fmt.Fprint(w, `
	<html>
		<body>
		<form>
  	    	<input type="text" value="Enter a playlist URL">
   	    	<input type="submit" value="Submit">
    	</form>
		</body>
	</html>	
	`)
	//fmt.Fprintf(w, "Login Completed!")

	ch <- client
}
