package content

import (
	"fmt"
	"os"
	"os/exec"
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
	Description string
	Duration    float64
	Path        string
	Status      ContentStatus
	Title       string
	Url         string
}

func ApplyOverlayWithFade(contentObjs []*ContentObj, contentPath string) error {
	filters := generateOveralyWithFadeArgs(contentObjs)

	// create and execute command
	args := []string{"-i", contentPath, "-vf", filters, "-codec:a", "copy", "finished-vid.mp4"}
	ffmpegCmd := exec.Command("ffmpeg", args...)
	fmt.Println("Running ffmpeg command:\n%s", ffmpegCmd.String())
	err := ffmpegCmd.Run()
	if err != nil {
		fmt.Errorf("Failed to execute ffmpeg cmd: %v\n", err)
	}

	return nil
}

const (
	duration           = 4
	drawtextTimeA      = `drawtext=enable='between(t,`
	drawtextTimeB      = `)':`
	drawtextFont       = `fontfile=/usr/share/fonts/TTF/OptimusPrinceps.ttf:`
	drawtextProperties = `fontcolor=white:fontsize=24:box=1:boxcolor=black@0.5:boxborderw=5:x=0:y=0`
)

// TODO: Decouple from download -- there should be an abstraction here
func ApplyOverlay(contentObjs []*ContentObj, contentPath string) error {
	var cursor float64
	var allFilters string
	for i, v := range contentObjs {
		// create drawtextTime
		timeA := cursor + 1.0
		timeB := cursor + 5.0
		drawtextTime := fmt.Sprintf("%s%f,%f%s", drawtextTimeA, timeA, timeB, drawtextTimeB)

		// create overlay
		overlayText := fmt.Sprintf("text='%s\n%s':", v.Title, v.CreatorName)
		fullOverlay := drawtextTime + drawtextFont + overlayText + drawtextProperties
		allFilters += fullOverlay
		if i < len(contentObjs)-1 {
			allFilters += ","
		}

		// move cursor
		cursor += v.Duration
	}

	// create and execute command
	args := []string{"-i", contentPath, "-vf", allFilters, "-codec:a", "copy", "finished-vid.mp4"}
	ffmpegCmd := exec.Command("ffmpeg", args...)
	fmt.Println("Running ffmpeg command:\n%s", ffmpegCmd.String())
	err := ffmpegCmd.Run()
	if err != nil {
		fmt.Errorf("Failed to execute ffmpeg cmd: %v\n", err)
	}

	return nil
}

// TODO: decouple encoding from compiling step
func Compile(contentObjs []*ContentObj, outfile string) (*ContentObj, error) {
	f, err := os.Create("compile.txt")
	if err != nil {
		return nil, err
	}
	defer f.Close()

	for _, v := range contentObjs {
		filename := strings.SplitN(v.Url, "twitch.tv", 2)[1]
		basename := filename[:len(filename)-4]
		processedPath := "tmp/processed/" + basename + ".mkv"
		cmd := exec.Command("ffmpeg", "-i", v.Url, "-c:v", "libx264", "-preset", "slow", "-crf", "22", "-c:a", "copy", processedPath)
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

// TODO: Clean this up, new generate ffmpeg command library?
const (
	tduration = 3.0
	font      = `drawtext=fontfile=/usr/share/fonts/TTF/OptimusPrinceps.ttf:`
	fontSize  = `fontsize=32:`
	fontColor = `fontcolor=ffffff:`
	xPos      = `x=0:`
	yPos      = `y=(h-text_h)`
)

func generateOveralyWithFadeArgs(contentObjs []*ContentObj) (allFilters string) {
	fade := 1.0
	var cursor float64
	for i, v := range contentObjs {
		overlayText := fmt.Sprintf("text='%s\n%s':", v.Title, v.CreatorName)
		alpha := fmt.Sprintf(`alpha='if(lt(t,%f),0,if(lt(t,%f),(t-1)/1,if(lt(t,%f),1,if(lt(t,%f),(1-(t-%f))/1,0))))':`, cursor+1.0, cursor+1.0+fade, cursor+tduration, cursor+tduration+fade, cursor+tduration)
		fullOverlay := font + overlayText + fontSize + fontColor + alpha + xPos + yPos
		allFilters += fullOverlay
		if i < len(contentObjs)-1 {
			allFilters += ","
		}

		cursor += v.Duration
	}

	return allFilters
}
