package content

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
)

type ContentStatus int

const (
	UNKNOWN   ContentStatus = 0
	RAW       ContentStatus = 1
	PROCESSED ContentStatus = 2
)

type ContentObj struct {
	CreatorName string
	Title       string
	Description string
	Path        string
	Status      ContentStatus
	Url         string
}

func (c *ContentObj) Download(path string) error {
	fmt.Println("Downloading new clip: ", c.Url)

	resp, err := http.Get(c.Url)
	defer resp.Body.Close()
	if err != nil {
		return err
	}

	filename := strings.SplitN(c.Url, "twitch.tv", 2)[1]
	outFile := path + filename
	c.Path = outFile

	out, err := os.Create(outFile)
	defer out.Close()
	if err != nil {
		return err
	}

	_, err = io.Copy(out, resp.Body)
	c.Status = RAW

	return err
}
