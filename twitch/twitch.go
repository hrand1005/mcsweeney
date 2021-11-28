package twitch


import (
	"fmt"
	"github.com/nicklaw5/helix"
	"io"
	"net/http"
	"os"
	"strings"
    //"sync"
    "time"
)


const (
	RawVidsDir = "tmp/raw"
)


func GetClips(clientID string, token string, gameID string, first int) error {
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
        First: first,
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


func DownloadNewClips(manyClips []helix.Clip) error {
    start := time.Now()
    fmt.Println("Download start...")

    //var wg sync.WaitGroup

    // TODO: spawn goroutines for each download
    // NOTE: Preliminary testing indicates not much of a difference around 14
    // clips downloaded using this commented out method. 
    fmt.Printf("Trying to download %v clips.\n", len(manyClips))
	for _, v := range manyClips {
		// TODO: Verify clip here
        //wg.Add(1)
        //v := v
		fmt.Println("Attempting to download a clip...")
        //go func() {
            //defer wg.Done()
        err := downloadClip(&v)
        if err != nil {
            fmt.Println("Failed to download a clip: ", err)
        }
        //}()
	}

    //wg.Wait()
    finish := time.Now()
    elapsed := finish.Sub(start)
    fmt.Printf("finished in %v\n", elapsed)

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
