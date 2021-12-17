package get

import (
	"fmt"
	"github.com/nicklaw5/helix"
	"mcsweeney/config"
	"mcsweeney/content"
	"strings"
	"time"
)

type TwitchGetter struct {
	client *helix.Client
	query  *helix.ClipsParams
	token  string
}

func NewTwitchGetter(clientID string, queryArgs config.Query, token string) (*TwitchGetter, error) {
	client, err := helix.NewClient(&helix.Options{
		ClientID: clientID,
	})
	if err != nil {
		return nil, fmt.Errorf("Couldn't create client: %v", err)
	}
	query, err := buildQuery(queryArgs)
	if err != nil {
		return nil, fmt.Errorf("Couldn't build query for TwitchGetter: %v", err)
	}

	return &TwitchGetter{
		client: client,
		query:  query,
		token:  token,
	}, nil
}

func (t *TwitchGetter) GetContent() ([]*content.ContentObj, error) {
	t.client.SetUserAccessToken(t.token)
	defer t.client.SetUserAccessToken("")

	twitchResp, err := t.client.GetClips(t.query)
	if err != nil {
		return nil, err
	}
	// updates the 'cursor' for where to start retrieving data
	t.query.After = twitchResp.Data.Pagination.Cursor

	clips := twitchResp.Data.Clips
	if err != nil || len(clips) == 0 {
		return nil, fmt.Errorf("Couldn't get clips: %v", err)
	}

	contentObjs := make([]*content.ContentObj, 0, len(clips))
	for _, v := range clips {
		contentObjs = append(contentObjs, convertClipToContentObj(&v))
	}

	return contentObjs, nil
}

func buildQuery(queryArgs config.Query) (*helix.ClipsParams, error) {
	var startTimeFormatted time.Time
	switch queryArgs.StartTime {
	case "yesterday":
		startTimeFormatted = time.Now().AddDate(0, 0, -1)
	}

	return &helix.ClipsParams{
		GameID:    queryArgs.GameID,
		First:     queryArgs.First,
		StartedAt: helix.Time{startTimeFormatted},
	}, nil
}

func convertClipToContentObj(clip *helix.Clip) *content.ContentObj {
	c := &content.ContentObj{}
	c.CreatorName = clip.BroadcasterName
	c.Duration = clip.Duration
	c.Title = clip.Title
	c.Url = strings.SplitN(clip.ThumbnailURL, "-preview", 2)[0] + ".mp4"

	return c
}
