package content

import (
	"fmt"
	"google.golang.org/api/youtube/v3"
	"net/http"
	"os"
	//"strings"
)

type Privacy string

const (
	PRIVATE Privacy = "private"
	PUBLIC  Privacy = "public"
)

type YoutubeSharer struct {
	client *http.Client
}

func newYoutubeSharer(credentials string) (*YoutubeSharer, error) {
	client := GetClient(credentials, youtube.YoutubeUploadScope)
	return &YoutubeSharer{
		client: client,
	}, nil
}

func (y *YoutubeSharer) Share(v Component) error {
	//TODO: perform checks on the inputs
	if v.Path() == "" {
		return fmt.Errorf("cannot upload nil file")
	}
	service, err := youtube.New(y.client)
	if err != nil {
		return fmt.Errorf("Couldn't create youtube service: %v", err)
	}

	upload := &youtube.Video{
		Snippet: &youtube.VideoSnippet{
			//Title:       v.Title,
			//Description: v.Description,
			//TODO: this might be nice :)
			//CategoryId:  v.CategoryID,
			//Tags: strings.Split(v.Keywords, ","),
		},
		//Status: &youtube.VideoStatus{PrivacyStatus: string(v.Privacy)},
	}

	insertArgs := []string{"snippet", "status"}
	call := service.Videos.Insert(insertArgs, upload)

	file, err := os.Open(v.Path())
	defer file.Close()
	if err != nil {
		return fmt.Errorf("Couldn't open file: %s, %v", v.Path(), err)
	}

	response, err := call.Media(file).Do()
	if err != nil {
		return fmt.Errorf("Couldn't upload file: %v", err)
	}

	fmt.Printf("%s uploaded successfully!", v.Path())
	fmt.Println("Response: ", response)

	return nil
}
