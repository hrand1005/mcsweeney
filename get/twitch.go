package get

import (
	"fmt"
	"github.com/nicklaw5/helix"
	"io"
	"mcsweeney/config"
	"mcsweeney/db"
	"net/http"
	"os"
	"strings"
	//"sync"
	"time"
)

// TODO: remove duplicate
const (
	RawVidsDir       = "tmp/raw"
	ProcessedVidsDir = "tmp/processed"
)

type TwitchGetter struct {
	client *helix.Client
	db     db.ContentDB
	query  *helix.ClipsParams
	token  string
}

func NewTwitchGetter(c config.Config, db db.ContentDB) (*TwitchGetter, error) {
	client, err := helix.NewClient(&helix.Options{
		ClientID: c.ClientID,
	})
	if err != nil {
		return nil, fmt.Errorf("Couldn't create TwitchGetter: %v", err)
	}
	query := &helix.ClipsParams{
		GameID: c.GameID,
		First:  c.First,
	}

	return &TwitchGetter{
		client: client,
		db:     db,
		query:  query,
		token:  c.Token,
	}, nil
}

// TODO: change this to return a content interface
func (t *TwitchGetter) GetContent() ([]helix.Clip, error) {
	t.client.SetUserAccessToken(t.token)
	defer t.client.SetUserAccessToken("")

	// Execute query for clips, TODO: more error checking here?
	twitchResp, err := t.client.GetClips(t.query)
	if err != nil {
		return nil, err
	}

	dirtyClips := twitchResp.Data.Clips
	if err != nil || len(dirtyClips) == 0 {
		return nil, fmt.Errorf("Couldn't get clips: %v", err)
	}

	err = downloadNewClips(dirtyClips)
	if err != nil {
		return nil, err
	}

	cleanClips := make([]helix.Clip, len(dirtyClips))
	for i, v := range dirtyClips {
		exists, err := t.db.Exists(v.URL)
		if err != nil {
			return nil, err
		}
		if !exists {
			cleanClips[i] = v
			err = t.db.Insert(v.URL)
			if err != nil {
				return nil, err
			}
		}
	}

	return cleanClips, nil
}

func downloadNewClips(manyClips []helix.Clip) error {
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
	defer resp.Body.Close()
	if err != nil {
		return err
	}

	filename := strings.SplitN(mp4URL, "twitch.tv", 2)[1]
	outFile := RawVidsDir + filename

	out, err := os.Create(outFile)
	defer out.Close()
	if err != nil {
		return err
	}

	_, err = io.Copy(out, resp.Body)

	return err
}
