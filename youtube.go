package main

import (
	"context"
	"fmt"
	"os"

	"google.golang.org/api/googleapi"
	"google.golang.org/api/option"
	"google.golang.org/api/youtube/v3"
)

const (
	YoutubeAPIKey     = "YOUTUBE_API_TOKEN"
)

type YoutubeClient struct {
  *youtube.Service
}

func NewYoutubeClient() (*YoutubeClient, error) {
  service, err := youtube.NewService(context.Background(), option.WithAPIKey(os.Getenv(YoutubeAPIKey)))
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

  f, err := os.Open(path)
  if err != nil {
    return googleapi.ServerResponse{}, err
  }
  defer f.Close()

  r, err := call.Media(f).Do()
  if err != nil {
    fmt.Printf("UploadVideo: encountered error: %v", err)
    return googleapi.ServerResponse{}, err
  }
  fmt.Printf("Got server response: %v", r)
  fmt.Printf("Http Status code: %v", r.HTTPStatusCode)
  return r.ServerResponse, err
}
