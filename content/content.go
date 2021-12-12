package content

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
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

const (
	drawtextFont       = `drawtext=fontfile=/usr/share/fonts/TTF/DejaVuSans.ttf:`
	drawtextProperties = `fontcolor=white:fontsize=24:box=1:boxcolor=black@0.5:boxborderw=5:x=0:y=0`
)

// TODO: replace this with a ffmpeg library dear god
// TODO: goroutines!
func (c *ContentObj) ApplyOverlay(outDir string) error {
	// create overlay
	overlayText := fmt.Sprintf("text='%s\n%s':", c.Title, c.CreatorName)
	overlayArg := drawtextFont + overlayText + drawtextProperties

	// create paths
	filename := filepath.Base(c.Path)
	processedPath := outDir + "/" + filename

	// create and execute command
	args := []string{"-i", c.Path, "-vf", overlayArg, "-codec:a", "copy", processedPath}
	ffmpegCmd := exec.Command("ffmpeg", args...)
	fmt.Printf("Attempting to run ffmpeg command: %s", ffmpegCmd.String())
	err := ffmpegCmd.Run()
	if err != nil {
		fmt.Printf("Failed to execute ffmpeg cmd: %v\n", err)
	}

	c.Path = processedPath
	c.Status = PROCESSED

	return nil
}

func Compile(contentObjs []*ContentObj) error {
	f, err := os.Create("clips.txt")
	if err != nil {
		return err
	}
	defer f.Close()

	for _, v := range contentObjs {
		// write to txt file
		w := fmt.Sprintf("file '%s'\n", v.Path)
		_, err = f.WriteString(w)
		if err != nil {
			return err
		}
	}

	args := []string{"-f", "concat", "-safe", "0", "-i", "clips.txt", "compiled-vid.mp4"}
	cmd := exec.Command("ffmpeg", args...)
	err = cmd.Run()
	if err != nil {
		fmt.Printf("Failed to execute ffmpeg cmd: %v\n", err)
	}

	return nil
}
