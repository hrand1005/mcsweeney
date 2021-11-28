package twitch

import (
	"fmt"
	"github.com/nicklaw5/helix"
	"io"
	"net/http"
	"os"
	"strings"
)

const (
	RawVidsDir = "tmp/"
)

func GetClips(clientID string, token string, gameID string) error {
	client, err := helix.NewClient(&helix.Options{
		ClientID: clientID,
	})
	if err != nil {
		return err
	}

	client.SetUserAccessToken(token)
	defer client.SetUserAccessToken("")

	// Define query for clips
	clipParams := &helix.ClipsParams{
		GameID: gameID,
	}

	// Execute query for clips
	twitchResp, err := client.GetClips(clipParams)
	if err != nil {
		return err
	}

	err = DownloadNewClips(twitchResp.Data.Clips)
	if err != nil {
		return err
	}

	return nil
}

// TODO: Validate clips here, then download all new
func DownloadNewClips(manyClips []helix.Clip) error {
	for _, v := range manyClips {
		// TODO: Verify clip here
		fmt.Println("Attempting to download a clip...")
		err := downloadClip(&v)
		if err != nil {
			fmt.Println("Failed to download a clip: ", err)
		}
	}

	return nil
}

func downloadClip(clip *helix.Clip) error {
	thumbURL := clip.ThumbnailURL
	mp4URL := strings.SplitN(thumbURL, "-preview", 2)[0] + ".mp4"
	fmt.Println("MP4 URL: ", mp4URL)

	resp, err := http.Get(mp4URL)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	filename := strings.SplitN(mp4URL, "twitch.tv", 2)[1]
	outFile := RawVidsDir + filename
	out, err := os.Create(outFile)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, resp.Body)

	return err
}
