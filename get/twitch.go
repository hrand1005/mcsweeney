package get

import (
	"fmt"
	"github.com/nicklaw5/helix"
	"mcsweeney/config"
	"mcsweeney/content"
	"mcsweeney/db"
	"strings"
	"time"
	//"sync"
)

type TwitchGetter struct {
	client *helix.Client
	query  *helix.ClipsParams
	token  string
}

func NewTwitchGetter(c *config.Config) (*TwitchGetter, error) {
	client, err := helix.NewClient(&helix.Options{
		ClientID: c.ClientID,
	})
	if err != nil {
		return nil, fmt.Errorf("Couldn't create TwitchGetter: %v", err)
	}
	query, err := buildQuery(c.GameID, c.First, c.StartTime)
	if err != nil {
		return nil, fmt.Errorf("Couldn't build query for TwitchGetter: %v", err)
	}

	return &TwitchGetter{
		client: client,
		query:  query,
		token:  c.Token,
	}, nil
}

func (t *TwitchGetter) GetContent(db db.ContentDB) ([]*content.ContentObj, error) {
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

	newContent := make([]*content.ContentObj, 0, len(dirtyClips))
	for _, v := range dirtyClips {
		contentObj, err := convertClipToContentObj(&v)
		exists, err := db.Exists(contentObj)
		if err != nil {
			return nil, err
		}
		if !exists {
			newContent = append(newContent, contentObj)
			err = db.Insert(contentObj)
			if err != nil {
				return nil, err
			}
		}
	}
	if len(newContent) == 0 {
		return nil, fmt.Errorf("No new clips retrieved.")
	}
	fmt.Printf("Downloaded %v new clips.\n", len(newContent))

	return newContent, nil
}

func buildQuery(gameId string, first int, start string) (*helix.ClipsParams, error) {
	var startTimeFormatted time.Time
	switch start {
	case "yesterday":
		startTimeFormatted = time.Now().AddDate(0, 0, -1)
	}

	return &helix.ClipsParams{
		GameID:    gameId,
		First:     first,
		StartedAt: helix.Time{startTimeFormatted},
	}, nil
}

func convertClipToContentObj(clip *helix.Clip) (*content.ContentObj, error) {
	c := &content.ContentObj{}

	c.CreatorName = clip.BroadcasterName
	c.Duration = clip.Duration
	c.Title = clip.Title
	c.Url = strings.SplitN(clip.ThumbnailURL, "-preview", 2)[0] + ".mp4"

	return c, nil
}
