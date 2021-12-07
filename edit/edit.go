package edit

import (
	"fmt"
	"github.com/nicklaw5/helix"
	"os"
	"os/exec"
	"strings"
	"time"
)

const (
	RawVidsDir         = "tmp/raw"
	ProcessedVidsDir   = "tmp/processed"
	drawtextFont       = `drawtext=fontfile=/usr/share/fonts/TTF/DejaVuSans.ttf:`
	drawtextProperties = `fontcolor=white:fontsize=24:box=1:boxcolor=black@0.5:boxborderw=5:x=0:y=0`
)

// TODO: replace this with a ffmpeg library dear god
// TODO: goroutines!
func ApplyOverlay(clips []helix.Clip) error {
	f, err := os.Create("clips.txt")
	if err != nil {
		return err
	}
	defer f.Close()

	for _, v := range clips {
		// create overlay
		overlayText := fmt.Sprintf("text='%s\n%s':", v.Title, v.BroadcasterName)
		overlayArg := drawtextFont + overlayText + drawtextProperties

		// create paths
		filename := getClipPath(&v)
		rawPath := RawVidsDir + filename
		processedPath := ProcessedVidsDir + filename

		// create and execute command
		args := []string{"-i", rawPath, "-vf", overlayArg, "-codec:a", "copy", processedPath}
		ffmpegCmd := exec.Command("ffmpeg", args...)
		err := ffmpegCmd.Run()
		if err != nil {
			fmt.Printf("Failed to execute ffmpeg cmd: %v\n", err)
		}

		// write to txt file
		w := fmt.Sprintf("file '%s'\n", processedPath)
		_, err = f.WriteString(w)
		if err != nil {
			return err
		}
	}

	return nil
}

func Compile() error {
	cmdName := "ffmpeg"
	args := []string{"-f", "concat", "-safe", "0", "-i", "clips.txt", "compiled-vid.mp4"}
	cmd := exec.Command(cmdName, args...)
	err := cmd.Run()
	if err != nil {
		fmt.Printf("Failed to execute ffmpeg cmd: %v\n", err)
	}

	return nil
}

//TODO: get rid of this
func getClipPath(clip *helix.Clip) string {
	thumbURL := clip.ThumbnailURL
	mp4URL := strings.SplitN(thumbURL, "-preview", 2)[0] + ".mp4"
	filename := strings.SplitN(mp4URL, "twitch.tv", 2)[1]

	return filename
}

// some of that experimental stuff
type clipFunc func([]helix.Clip) error

func clipFuncTimer(f clipFunc) clipFunc {
	return func(c []helix.Clip) error {
		defer func(t time.Time) {
			fmt.Printf("clipFunc elapsed in %v\n", time.Since(t))
		}(time.Now())

		return f(c)
	}
}
