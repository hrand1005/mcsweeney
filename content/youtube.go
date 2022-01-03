package content

import (
	"fmt"
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

type YoutubeSharer struct {
	client *http.Client
}

func NewYoutubeSharer(credentials string) (*YoutubeSharer, error) {
	client := GetClient(credentials, youtube.YoutubeUploadScope)
	return &YoutubeSharer{
		client: client,
	}, nil
}

func (y *YoutubeSharer) Share(p Payload) error {
	//TODO: perform checks on the inputs
	if p.Path == "" {
		return fmt.Errorf("cannot upload nil file")
	}
	service, err := youtube.New(y.client)
	if err != nil {
		return fmt.Errorf("Couldn't create youtube service: %v", err)
	}

	upload := &youtube.Video{
		Snippet: &youtube.VideoSnippet{
			Title:       p.Title,
			Description: p.Description,
			//TODO: this might be nice :)
			CategoryId: p.CategoryID,
			Tags:       strings.Split(p.Keywords, ","),
		},
		Status: &youtube.VideoStatus{PrivacyStatus: string(p.Privacy)},
	}

	insertArgs := []string{"snippet", "status"}
	call := service.Videos.Insert(insertArgs, upload)

	file, err := os.Open(p.Path)
	defer file.Close()
	if err != nil {
		return fmt.Errorf("Couldn't open file: %s, %v", p.Path, err)
	}

	response, err := call.Media(file).Do()
	if err != nil {
		return fmt.Errorf("Couldn't upload file: %v", err)
	}

	fmt.Printf("%s uploaded successfully!", p.Path)
	fmt.Println("Response: ", response)

	return nil
}
