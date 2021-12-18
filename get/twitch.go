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

func NewTwitchGetter(credentials config.Credentials, query config.Query) (*TwitchGetter, error) {
	client, err := helix.NewClient(&helix.Options{
		ClientID: credentials.ClientID,
	})
	if err != nil {
		return nil, fmt.Errorf("Couldn't create client: %v", err)
	}
	twitchQuery, err := buildQuery(query)
	if err != nil {
		return nil, fmt.Errorf("Couldn't build query for TwitchGetter: %v", err)
	}

	return &TwitchGetter{
		client: client,
		query:  twitchQuery,
		token:  credentials.Token,
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

func buildQuery(query config.Query) (*helix.ClipsParams, error) {
	// start time for query, specified in config by 'days'
	startTime := time.Now().AddDate(0, 0, -1*query.Days)

	return &helix.ClipsParams{
		GameID:    query.GameID,
		First:     query.First,
		StartedAt: helix.Time{startTime},
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
