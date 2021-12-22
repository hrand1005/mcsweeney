package content

import (
	"fmt"
	"mcsweeney/config"
	"os"
	"os/exec"
	"strings"
)

func (c *Content) ApplyOverlay(contentObjs []*Content, options config.Options) error {
	filters := generateOverlayWithFadeArgs(contentObjs, options.Overlay)
	args := []string{"-i", c.Path, "-vf", filters, "-codec:a", "copy", "finished-vid.mp4"}

	ffmpegCmd := exec.Command("ffmpeg", args...)
	err := ffmpegCmd.Run()
	if err != nil {
		return fmt.Errorf("Failed to execute ffmpeg cmd\n%s\nerr: %v\n", ffmpegCmd.String(), err)
	}

	// update content path
	c.Path = "finished-vid.mp4"

	return nil
}

// TODO: decouple encoding from compiling step
// Compile takes a slice of Content objects, encodes them consistently and then
// concatenates them to create a new content object. As part of this, it credits
// it's subobjects in the description.
func Compile(contentObjs []*Content, outfile string) (*Content, error) {
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

	return &Content{
		Path:        outfile,
		Description: buildCredits(contentObjs),
	}, nil
}

// TODO: Clean this up, new generate ffmpeg command library?
const (
	xPos = `x=10:`
	yPos = `y=(h-text_h)-10`
)

func generateOverlayWithFadeArgs(contentObjs []*Content, args config.Overlay) (allFilters string) {
	fade := float64(args.Fade)
	duration := float64(args.Duration)
	font := fmt.Sprintf("drawtext=fontfile=%s:", args.Font)
	fontColor := fmt.Sprintf("fontcolor=%s:", args.Color)
	fontSize := fmt.Sprintf("fontsize=%s:", args.Size)

	var cursor float64
	for i, v := range contentObjs {
		// TODO: find workaround, escaping with `\'` doesn't work
		title := strings.ReplaceAll(v.Title, `'`, ``)
		creator := strings.ReplaceAll(v.CreatorName, `'`, ``)
		overlayText := fmt.Sprintf("text='%s\n%s':", title, creator)
		fmt.Printf("\nAppling overlay text:\n%s\n", overlayText)
		alpha := fmt.Sprintf(`alpha='if(lt(t,%f),0,if(lt(t,%f),(t-%f)/1,if(lt(t,%f),1,if(lt(t,%f),(1-(t-%f))/1,0))))':`, cursor+1.0, cursor+1.0+fade, cursor+1.0, cursor+duration, cursor+duration+fade, cursor+duration)
		fullOverlay := font + overlayText + fontSize + fontColor + alpha + xPos + yPos
		allFilters += fullOverlay
		if i < len(contentObjs)-1 {
			allFilters += ","
		}

		cursor += v.Duration
	}

	return allFilters
}

func buildCredits(contentObjs []*Content) (credits string) {
	cursor := 0.0
	for _, v := range contentObjs {
		// simple youtube timestamp up to 59:59
		minutes := int(cursor) / 60
		seconds := int(cursor) % 60
		var timestamp string
		if seconds < 10 {
			timestamp = fmt.Sprintf("%v:0%v", minutes, seconds)
		} else {
			timestamp = fmt.Sprintf("%v:%v", minutes, seconds)
		}
		credits += fmt.Sprintf("\n\n%s '%s'\nStreamed by %s at %s\nClipped by %s\n", timestamp, v.Title, v.CreatorName, v.Channel, v.ClippedBy)
		cursor += v.Duration
	}

	return
}
