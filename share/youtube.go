package share

import (
	"fmt"
	"google.golang.org/api/youtube/v3"
	"mcsweeney/config"
	"os"
	"strings"
)

type YoutubeSharer struct {
	filename    string
	title       string
	description string
	category    string
	keywords    string
	privacy     string
}

func NewYoutubeSharer(c config.Config, path string) (*YoutubeSharer, error) {
	// TODO: validate args?
	// TODO: should sharer object be reusable? ie should we not do this?
	return &YoutubeSharer{
		filename:    path,
		title:       c.Title,
		description: c.Description,
		keywords:    c.Keywords,
		privacy:     c.Privacy,
	}, nil
}

func (y *YoutubeSharer) Share() error {
	if y.filename == "" {
		return fmt.Errorf("cannot upload nil file")
	}
	//TODO: perform checks on the inputs

	// TODO: figure out this
	client := GetClient(youtube.YoutubeUploadScope)

	service, err := youtube.New(client)
	if err != nil {
		return fmt.Errorf("Couldn't create youtube service: %v", err)
	}

	upload := &youtube.Video{
		Snippet: &youtube.VideoSnippet{
			Title:       y.title,
			Description: y.description,
			//TODO: this might be nice :)
			//CategoryId:  y.category,
			Tags: strings.Split(y.keywords, ","),
		},
		Status: &youtube.VideoStatus{PrivacyStatus: y.privacy},
	}

	insertArgs := []string{"snippet", "status"}
	call := service.Videos.Insert(insertArgs, upload)

	file, err := os.Open(y.filename)
	defer file.Close()
	if err != nil {
		return fmt.Errorf("Couldn't open file: %s, %v", y.filename, err)
	}

	response, err := call.Media(file).Do()
	if err != nil {
		return fmt.Errorf("Couldn't upload file: %v", err)
	}

	fmt.Printf("%s uploaded successfully!\nTitle: %s\n", y.filename, y.title)
	fmt.Printf("Response:\n%v", response)

	return nil
}
