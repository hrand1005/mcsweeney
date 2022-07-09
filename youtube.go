package main

import (
	"context"
	"fmt"
	// "io"
	"encoding/json"
	"net"
	"net/http"
	"os"
	"os/exec"
	"runtime"

	"golang.org/x/oauth2"
	"google.golang.org/api/googleapi"
	"google.golang.org/api/option"
	"google.golang.org/api/youtube/v3"
)

const (
	YoutubeClientID     = "YOUTUBE_CLIENT_ID"
	YoutubeClientSecret = "YOUTUBE_CLIENT_SECRET"
	YoutubeTokenFile = "YOUTUBE_TOKEN_FILE"
	YoutubeAuthURI      = "https://accounts.google.com/o/oauth2/auth"
	YoutubeTokenURI     = "https://oauth2.googleapis.com/token"
	// spin up local web server for authentication
	LocalWebServer = "localhost:8090"
)

type YoutubeClient struct {
	*youtube.Service
}

func NewYoutubeClient() (*YoutubeClient, error) {
	config := oauth2.Config{
		ClientID:     os.Getenv(YoutubeClientID),
		ClientSecret: os.Getenv(YoutubeClientSecret),
		Scopes:       []string{youtube.YoutubeUploadScope},
		RedirectURL:  "http://" + LocalWebServer,
		Endpoint: oauth2.Endpoint{
			AuthURL:  YoutubeAuthURI,
			TokenURL: YoutubeTokenURI,
		},
	}

	token, err := getYoutubeAppToken(config)
	if err != nil {
		return nil, err
	}

	ctx := context.Background()
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

func getTokenFromWebServer(c oauth2.Config, url string) (*oauth2.Token, error) {
	codeChan, err := startWebServer()
	if err != nil {
		return nil, fmt.Errorf("failed to start webserver: %v", err)
	}

	err = openURL(url)
	if err != nil {
		return nil, fmt.Errorf("failed to open url %q in webserver: %v", url, err)
	}

	fmt.Println("Waiting for authorization...")
	code := <- codeChan

	// TODO: use real context?
	return c.Exchange(oauth2.NoContext, code)
}

func openURL(url string) error {
	var err error
	switch runtime.GOOS {
	case "linux": 
		err = exec.Command("xdg-open", url).Start()
	default: 
		err = fmt.Errorf("Environment not supported: %v", runtime.GOOS)
	}

	return err
}

func startWebServer() (chan string, error) {
	listener, err := net.Listen("tcp", LocalWebServer)
	if err != nil {
		return nil, err
	}
	codeChan := make(chan string)

	go http.Serve(listener, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		code := r.FormValue("code")
		codeChan <- code
		listener.Close()
		w.Header().Set("Content-Type", "text/plain")
		fmt.Fprintf(w, "Recieved code: %v\r\nYou may now safely close this browser window.", code)
	}))

	return codeChan, nil
}

func getYoutubeAppToken(config oauth2.Config) (*oauth2.Token, error) {
	url := config.AuthCodeURL("state-token", oauth2.AccessTypeOffline)
	fmt.Printf("Visit the following URL for the auth dialogue: %v", url)

	var token *oauth2.Token
	// try getting token from the env variables
	token, err := readTokenFromFile(os.Getenv(YoutubeTokenFile))
	if err != nil {
		fmt.Printf("Failed to read token from file, err: %v\nLaunching web server to get new token", err)
		token, err = getTokenFromWebServer(config, url)
		if err != nil {
			fmt.Printf("Failed to get token from web server: %v\n", err)
			return nil, fmt.Errorf("failed to get token form web server: %v", err)
		}
		writeTokenToFile(os.Getenv(YoutubeTokenFile), token)
	}

	return token, nil
}

func readTokenFromFile(tokenFile string) (*oauth2.Token, error) {
	f, err := os.Open(tokenFile)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	t := &oauth2.Token{}
	err = json.NewDecoder(f).Decode(t)

	return t, err
}

func writeTokenToFile(file string, t *oauth2.Token) error {
	f, err := os.OpenFile(file, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		return err
	}
	defer f.Close()

	return json.NewEncoder(f).Encode(t)
}
