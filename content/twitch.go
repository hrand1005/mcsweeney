package content

import (
	"bufio"
	"fmt"
	"github.com/nicklaw5/helix"
	"os"
	"strings"
	"time"
)

type TwitchGetter struct {
	client *helix.Client
	query  *helix.ClipsParams
	token  string
}

func NewTwitchGetter(credentials string, query Query) (*TwitchGetter, error) {
	clientID, token, err := loadTwitchCredentials(credentials)
	if err != nil {
		return nil, fmt.Errorf("Couldn't load credentials, err: %v", err)
	}
	client, err := helix.NewClient(&helix.Options{
		ClientID: clientID,
	})
	if err != nil {
		fmt.Printf("ClientID: %s\nToken: %s\n", clientID, token)
		return nil, fmt.Errorf("Couldn't create client: %v", err)
	}

	twitchQuery := buildQuery(query)

	return &TwitchGetter{
		client: client,
		query:  twitchQuery,
		token:  token,
	}, nil
}

func (t *TwitchGetter) Get() ([]*Content, error) {
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

	contentObjs := make([]*Content, 0, len(clips))
	for _, v := range clips {
		contentObjs = append(contentObjs, convertClipToContentObj(&v))
	}

	return contentObjs, nil
}

func buildQuery(query Query) *helix.ClipsParams {
	// start time for query, specified in config by 'days'
	startTime := time.Now().AddDate(0, 0, -1*query.Days)

	return &helix.ClipsParams{
		GameID:    query.GameID,
		First:     query.First,
		StartedAt: helix.Time{Time: startTime},
	}
}

func convertClipToContentObj(clip *helix.Clip) (c *Content) {
	c = &Content{}
	c.Channel = "https://twitch.tv/" + clip.BroadcasterName
	c.ClippedBy = clip.CreatorName
	c.CreatorName = clip.BroadcasterName
	c.Duration = clip.Duration
	c.Language = clip.Language
	c.Title = clip.Title
	c.Url = strings.SplitN(clip.ThumbnailURL, "-preview", 2)[0] + ".mp4"

	return
}

func loadTwitchCredentials(path string) (clientID string, token string, err error) {
	file, err := os.Open(path)
	if err != nil {
		return "", "", fmt.Errorf("Couldn't load credentials from %s, err: %v", path, err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanLines)

	// first line should contain the clientID
	scanner.Scan()
	clientID = scanner.Text()

	// second line should contain the token
	scanner.Scan()
	token = scanner.Text()

	if err = scanner.Err(); err != nil {
		return "", "", fmt.Errorf("Couldn't scan items, err: %v", err)
	}

	return
}