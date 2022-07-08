package main

import (
	"context"
	"fmt"
	// "io"
	"os"

	"golang.org/x/oauth2"
	"google.golang.org/api/googleapi"
	"google.golang.org/api/option"
	"google.golang.org/api/youtube/v3"
)

const (
	YoutubeClientID     = "YOUTUBE_CLIENT_ID"
	YoutubeClientSecret = "YOUTUBE_CLIENT_SECRET"
	YoutubeAuthURI      = "https://accounts.google.com/o/oauth2/auth"
	YoutubeTokenURI     = "https://oauth2.googleapis.com/token"
	// use this option if not using a web server to redirect
	RedirectURI = "urn:ietf:wg:oauth:2.0:oob"
)

type YoutubeClient struct {
	*youtube.Service
}

func NewYoutubeClient() (*YoutubeClient, error) {
	ctx := context.Background()

	config := oauth2.Config{
		ClientID:     os.Getenv(YoutubeClientID),
		ClientSecret: os.Getenv(YoutubeClientSecret),
		Scopes:       []string{youtube.YoutubeUploadScope},
		RedirectURL:  RedirectURI,
		Endpoint: oauth2.Endpoint{
			AuthURL:  YoutubeAuthURI,
			TokenURL: YoutubeTokenURI,
		},
	}

	url := config.AuthCodeURL("state-token", oauth2.AccessTypeOffline)
	fmt.Printf("Visit the following URL for the auth dialogue: %v", url)

	var code string
	if _, err := fmt.Scan(&code); err != nil {
		return nil, err
	}

	token, err := config.Exchange(ctx, code)

	service, err := youtube.NewService(ctx, option.WithTokenSource(config.TokenSource(ctx, token)))
	if err != nil {
		return nil, err
	}

	return &YoutubeClient{
		service,
	}, nil
}

func (y *YoutubeClient) UploadVideo(path string, video *youtube.Video) (googleapi.ServerResponse, error) {
	insertArgs := []string{"snippet", "status"}
	call := y.Videos.Insert(insertArgs, video)

	fmt.Printf("VideoInsertCall: %v", call)

	f, err := os.Open(path)
	if err != nil {
		return googleapi.ServerResponse{}, err
	}
	defer f.Close()

	callWithMedia := call.Media(f)

	fmt.Printf("VideoInsertCall with media: %v", callWithMedia)

	r, err := callWithMedia.Do()
	if err != nil {
		fmt.Printf("UploadVideo: encountered error: %v", err)
		return googleapi.ServerResponse{}, err
	}
	fmt.Printf("Got server response: %v", r)
	fmt.Printf("Http Status code: %v", r.HTTPStatusCode)
	return r.ServerResponse, err
}
