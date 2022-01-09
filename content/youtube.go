package content

import (
	"google.golang.org/api/youtube/v3"
	"net/http"
	"os"
	"strings"
)

type Privacy string

const (
	PRIVATE Privacy = "private"
	PUBLIC  Privacy = "public"
)

// YoutubeSharer maintains client info for sharing new content via the youtube
// api.
type YoutubeSharer struct {
	client *http.Client
}

// NewYoutubeSharer constructs a new YoutubeSharer using a credentials json file
// that can be created using the google developer console.
func NewYoutubeSharer(credentials string) (*YoutubeSharer, error) {
	client := GetClient(credentials, youtube.YoutubeUploadScope)
	return &YoutubeSharer{
		client: client,
	}, nil
}

// Share implements the sharer interface for YoutubeSharer. It shares a payload
// to the channel corresponding to the YoutubeSharer's credentials.
func (y *YoutubeSharer) Share(p Payload) (int, error) {
	if p.Path == "" {
		return 0, ErrEmptyPath
	}
	service, err := youtube.New(y.client)
	if err != nil {
		return 0, err
	}
	// initialize a youtube object for upload conforming to the youtube api
	upload := &youtube.Video{
		Snippet: &youtube.VideoSnippet{
			Title:       p.Title,
			Description: enforceYoutubeConstraints(p.Description),
			CategoryId:  p.CategoryID,
			Tags:        strings.Split(p.Keywords, ","),
		},
		Status: &youtube.VideoStatus{PrivacyStatus: string(p.Privacy)},
	}
	insertArgs := []string{"snippet", "status"}
	call := service.Videos.Insert(insertArgs, upload)

	file, err := os.Open(p.Path)
	defer file.Close()
	if err != nil {
		return 0, err
	}
	r, err := call.Media(file).Do()

	return r.ServerResponse.HTTPStatusCode, err
}

func enforceYoutubeConstraints(s string) string {
	s = strings.ReplaceAll(s, `<`, ``)
	return strings.ReplaceAll(s, `>`, ``)
}
