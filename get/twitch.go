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
)

// TODO: remove duplicate
const (
	RawVidsDir       = "tmp/raw"
	ProcessedVidsDir = "tmp/processed"
)

type TwitchGetter struct {
	client *helix.Client
	query  *helix.ClipsParams
	token  string
}

func NewTwitchGetter(c config.Config) (*TwitchGetter, error) {
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
		query:  query,
		token:  c.Token,
	}, nil
}

// TODO: change this to return a content interface
func (t *TwitchGetter) GetContent(db db.ContentDB) ([]helix.Clip, error) {
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

	cleanClips := make([]helix.Clip, 0, len(dirtyClips))
	for _, v := range dirtyClips {
		exists, err := db.Exists(v.URL)
		if err != nil {
			return nil, err
		}
		if !exists {
			cleanClips = append(cleanClips, v)
			err = db.Insert(v.URL)
			if err != nil {
				return nil, err
			}
			// TODO: spawn goroutines here?
			err = downloadClip(&v)
			if err != nil {
				return nil, err
			}
		}
	}
	if len(cleanClips) == 0 {
		return nil, fmt.Errorf("No new clips retrieved.")
	}
	fmt.Printf("Downloaded %v new clips.\n", len(cleanClips))

	return cleanClips, nil
}

func downloadClip(clip *helix.Clip) error {
	thumbURL := clip.ThumbnailURL
	mp4URL := strings.SplitN(thumbURL, "-preview", 2)[0] + ".mp4"
	fmt.Println("Downloading new clip: ", mp4URL)

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
