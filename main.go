package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"
	"strconv"

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
	//10k list
	//yourPlaylist := "6qaVfh57zV2Y23B139X1Tn"
	yourPlaylist := "6ko0RCsHny1iOJSF5hbmQ7"
	//yourPlaylist := "3Wd692HY1qm450HUpXLDfE"
	//small list
	//yourPlaylist := "5SzZRpqqpxxhpURIDgiPyZ"
	tmpDir = tmpDir + "/" + yourPlaylist
	os.Mkdir(tmpDir, 775)
	generateCollage(yourPlaylist, client)

}

func generateCollage(playlistID string, client *spotify.Client) {

	playlist, err := client.GetPlaylistItems(context.Background(), spotify.ID(playlistID))
	if err != nil {
		log.Fatal(err)
	}

	for i := 0; i <= playlist.Total/100; i++ {
		playlist, err := client.GetPlaylistItems(context.Background(), spotify.ID(playlistID), spotify.Offset(i*100))
		if err != nil {
			log.Fatal(err)
		}

		//get track list from playlist
		//TODO use offset big playlists
		items := playlist.Items

		//Download imgs and ignore duplicates
		var dl string
		for _, item := range items {
			dl = item.Track.Track.Album.Images[2].URL
			_, dlErr := os.Stat(tmpDir + "/" + dl[25:] + ".jpg")
			if dlErr == nil {
				fmt.Println("File exists, skipping")
			} else if os.IsNotExist(dlErr) {
				fmt.Println("Downloading " + dl)
				downloadImage(dl)
			} else {
				fmt.Println("dl Err:", dlErr)
			}
		}

	}
	//removes excess images to create perfect squares
	//TODO please rework

	wccmd := exec.Command("bash", "-c", "ls -l "+tmpDir+"/* | wc -l")
	wcoutput, wcerr := wccmd.Output()
	if err != nil {
		if exitError, ok := err.(*exec.ExitError); ok {
			// The command exited with a non-zero status code
			errMsg := string(exitError.Stderr) // Error messages from stderr
			log.Fatal("Command failed with error:", errMsg)
		} else {
			// Other types of errors
			log.Fatal("Command execution failed:", wcerr)
		}
	} else {
		fmt.Println("Command executed successfully!")
		fmt.Println("Output:", string(wcoutput))
	}

	wcoutInt, err := strconv.Atoi(string(wcoutput)[:len(string(wcoutput))-1])
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	fmt.Println(wcoutInt)

	for i := 1; i <= 100; i++ {
		fmt.Println(i * i)
		if (i * i) > wcoutInt {
			n := wcoutInt - (i-1)*(i-1)
			fmt.Println(n)
			for j := 0; j < n; j++ {
				//TODO this is slow and stupid
				rmCmd := exec.Command("bash", "-c", "ls "+tmpDir+"/* | tail -n 1")
				out, err := rmCmd.Output()
				if err != nil {
					fmt.Println(err)
				} else {
					fmt.Println(string(out))
					err := os.Remove(string(out)[:len(string(out))-1])
					if err != nil {
						fmt.Println(err)
					}
				}
			}
			break
		}
	}

	//Create the montage/collage
	//Note: Up the disk limit on ImageMagicks policy in /etc/ImageMagic-6/policy.xml
	cmd := exec.Command("bash", "-c", "montage "+tmpDir+"/* -geometry +0+0 "+generatedDir+"/"+playlistID+".jpg")

	output, err := cmd.Output()
	if err != nil {
		if exitError, ok := err.(*exec.ExitError); ok {
			// The command exited with a non-zero status code
			errMsg := string(exitError.Stderr) // Error messages from stderr
			log.Fatal("Command failed with error:", errMsg)
		} else {
			log.Fatal("Command execution failed:", err)
		}
	} else {
		fmt.Println("Command executed successfully!")
		fmt.Println("Output:", string(output))
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
	fmt.Fprintf(w, "Login Completed!")
	ch <- client
}
