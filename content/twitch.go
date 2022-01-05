package content

import (
	"bufio"
	"errors"
	"fmt"
	"github.com/nicklaw5/helix"
	"os"
	"strings"
	"time"
)

// ErrNoClipsFound is returned when no clips can currently be found for the
// getter's query parameters.
var ErrNoClipsFound = errors.New("No clips could be retrieved for this getter.")

// TwitchGetter maintains client and query info for retrieving new twitch clips
// via the twitch api.
type TwitchGetter struct {
	client *helix.Client
	query  *helix.ClipsParams
	token  string
}

// NewTwitchGetter constructs a new TwitchGetter using a credentials filepath
// and a query object. The credentials file have the clientID on the first line
// and the token on the second line.
func NewTwitchGetter(credentials string, query Query) (*TwitchGetter, error) {
	clientID, token, err := loadTwitchCredentials(credentials)
	if err != nil {
		return nil, err
	}
	client, err := helix.NewClient(&helix.Options{
		ClientID: clientID,
	})
	if err != nil {
		return nil, fmt.Errorf("Error creating client.\nClientID: %s\nToken: %s\n: %v", clientID, token, err)
	}
	// build a twitchQuery conforming to the twitchAPI package
	twitchQuery := buildQuery(query)

	return &TwitchGetter{
		client: client,
		query:  twitchQuery,
		token:  token,
	}, nil
}

// Get implements the getter interface for TwitchGetter. It retrieves clips
// according to the configured query using the configured clientID and token,
// and maintains a cursor on retrieved clips so that subsequent calls don't
// produce duplicates.
func (t *TwitchGetter) Get() ([]*Clip, error) {
	t.client.SetUserAccessToken(t.token)
	defer t.client.SetUserAccessToken("")
	// use getter query and client to retrieve clips
	twitchResp, err := t.client.GetClips(t.query)
	if err != nil {
		return nil, err
	}
	// updates the 'cursor' for where to start retrieving data
	t.query.After = twitchResp.Data.Pagination.Cursor

	twitchClips := twitchResp.Data.Clips
	if err != nil {
		return nil, fmt.Errorf("Couldn't get clips: %v", err)
	}
	if len(twitchClips) == 0 {
		return nil, ErrNoClipsFound
	}

	// convert helix twitch clips to clip objects suitable for content package
	clips := make([]*Clip, 0, len(twitchClips))
	for _, v := range twitchClips {
		clips = append(clips, &Clip{
			Author:      v.CreatorName,
			Broadcaster: v.BroadcasterName,
			Duration:    v.Duration,
			Language:    v.Language,
			Path:        strings.SplitN(v.ThumbnailURL, "-preview", 2)[0] + ".mp4",
			Platform:    TWITCH,
			Title:       v.Title,
		})
	}

	return clips, nil
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

func loadTwitchCredentials(path string) (clientID string, token string, err error) {
	file, err := os.Open(path)
	if err != nil {
		return "", "", fmt.Errorf("Couldn't load twitch credentials from %s, err: %v", path, err)
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
