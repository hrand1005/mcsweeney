package content

import (
	"fmt"
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

const (
	drawtextFont       = `drawtext=fontfile=/usr/share/fonts/TTF/DejaVuSans.ttf:`
	drawtextProperties = `fontcolor=white:fontsize=24:box=1:boxcolor=black@0.5:boxborderw=5:x=0:y=0`
)

// TODO: Decouple from download -- there should be an abstraction here
func (c *ContentObj) ApplyOverlay(path string) error {
	fmt.Println("Downloading new clip: ", c.Url)
	// create overlay
	overlayText := fmt.Sprintf("text='%s\n%s':", c.Title, c.CreatorName)
	overlayArg := drawtextFont + overlayText + drawtextProperties

	// create paths
	filename := strings.SplitN(c.Url, "twitch.tv", 2)[1]
	outFile := path + filename
	c.Path = outFile

	fmt.Println("Applying overlay:\n", overlayArg)
	// create and execute command
	args := []string{"-i", c.Url, "-vf", overlayArg, "-codec:a", "copy", c.Path}
	ffmpegCmd := exec.Command("ffmpeg", args...)
	err := ffmpegCmd.Run()
	if err != nil {
		fmt.Errorf("Failed to execute ffmpeg cmd: %v\n", err)
	}

	c.Status = PROCESSED

	return nil
}

// TODO: decouple encoding from compiling step
func Compile(contentObjs []*ContentObj) (*ContentObj, error) {
	f, err := os.Create("compile.txt")
	if err != nil {
		return nil, err
	}
	defer f.Close()

	for _, v := range contentObjs {
		filename := filepath.Base(v.Path)
		basename := filename[:len(filename)-4]
		processedPath := "tmp/processed/" + basename + ".mkv"
		cmd := exec.Command("ffmpeg", "-i", v.Path, "-c:v", "libx264", "-preset", "slow", "-crf", "22", "-c:a", "copy", processedPath)
		fmt.Println("Encoding content to ", processedPath)
		err := cmd.Run()
		if err != nil {
			return nil, fmt.Errorf("Failed to download and encode to path %s: %v\n", v.Path, err)
		}

		v.Path = processedPath

		// write to txt file
		w := fmt.Sprintf("file '%s'\n", v.Path)
		_, err = f.WriteString(w)
		if err != nil {
			return nil, err
		}
	}

	fmt.Println("Compiling content...")
	outfile := "compiled-vid.mp4"
	args := []string{"-f", "concat", "-safe", "0", "-i", "compile.txt", outfile}
	cmd := exec.Command("ffmpeg", args...)
	err = cmd.Run()
	if err != nil {
		return nil, fmt.Errorf("Failed to execute ffmpeg cmd: %v\n", err)
	}

	return &ContentObj{
		Path: outfile,
	}, nil
}
