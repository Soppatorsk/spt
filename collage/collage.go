package collage

import (
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/zmb3/spotify/v2"
)

var (
	generatedDir = "img"
	tmpDir       = "tmp"
)

func GenerateCollage(playlistID string, client *spotify.Client) string {
	os.Mkdir(generatedDir, 775)
	os.Mkdir(tmpDir, 775)
	tmpDir = tmpDir + "/" + playlistID
	os.Mkdir(tmpDir, 775)

	playlist, err := client.GetPlaylistItems(context.Background(), spotify.ID(playlistID))
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Downloading images...")
	for i := 0; i <= playlist.Total/100; i++ {
		playlist, err := client.GetPlaylistItems(context.Background(), spotify.ID(playlistID), spotify.Offset(i*100))
		if err != nil {
			log.Fatal(err)
		}

		//get track list from playlist
		items := playlist.Items

		//Download imgs and ignore duplicates
		for _, item := range items {
			if len(item.Track.Track.Album.Images) > 2 {
				dl := item.Track.Track.Album.Images[2].URL

				_, dlErr := os.Stat(tmpDir + "/" + dl[25:] + ".jpg")
				if dlErr == nil {
					//
				} else if os.IsNotExist(dlErr) {
					fmt.Println("Downloading " + dl)
					downloadImage(dl, tmpDir)
				} else {
					fmt.Println("dl Err:", dlErr)
				}
			} else {
				//fmt.Println("Probably user local file, skipping")
			}
		}

	}
	//removes excess images to create perfect squares

	fmt.Println("Calculating and fixing for perfect square...")

	//TODO equivalent os call?
	wcOutput, err := exec.Command("bash", "-c", "ls -l "+tmpDir+"/* | wc -l").Output()
	if err != nil {
		log.Fatal("Command execution failed:", err)
	}

	wcCount, err := strconv.Atoi(strings.TrimSpace(string(wcOutput)))

	for i := 1; i <= 100; i++ {
		if i*i > wcCount {
			n := wcCount - (i-1)*(i-1)
			for j := 0; j < n; j++ {
				files, err := filepath.Glob(tmpDir + "/*")
				if err != nil || len(files) == 0 {
					fmt.Println("Error:", err)
					break
				}
				filePath := files[len(files)-1]
				err = os.Remove(filePath)
				if err != nil {
					fmt.Println(err)
				}
			}
			break
		}
	}

	fmt.Println("Creating collage...")
	//Create the montage/collage
	//Note: Up the disk limit on ImageMagicks policy in /etc/ImageMagic-6/policy.xml
	finalImage := generatedDir + "/" + playlistID + ".jpg"
	fmt.Println(finalImage)
	cmd := exec.Command("bash", "-c", "montage "+tmpDir+"/* -geometry 64x64+0+0 "+finalImage)

	output, err := cmd.Output()

	if err != nil {
		if exitError, ok := err.(*exec.ExitError); ok {
			// The command exited with a non-zero status code
			errMsg := string(exitError.Stderr) // Error messages from stderr
			log.Fatal("Command failed with error:", errMsg)
		} else {
			log.Fatal("Command execution failed:", err)
		}
	}
	fmt.Println("Command executed successfully!" + string(output))

	return "/" + finalImage
}

func downloadImage(url string, tmpDir string) error {
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
