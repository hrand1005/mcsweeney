package content

import (
	"fmt"
	"mcsweeney/config"
	"os"
	"os/exec"
	"strings"
)

func (c *Content) ApplyOverlay(contentObjs []*Content, options config.Options) error {
	bgFilter := generateOverlayBackground(contentObjs, options.Overlay)
	filters := generateOverlayWithFadeArgs(contentObjs, options.Overlay)
	bargs := make([]string, 0, len(contentObjs)*2+4)
	bargs = append(bargs, "-i", c.Path)
	for range contentObjs {
		bargs = append(bargs, "-i", options.Overlay.Background)
	}
	bargs = append(bargs, "-filter_complex", bgFilter+","+filters, "finished-vid.mp4")
	bgCmd := exec.Command("ffmpeg", bargs...)
	err := bgCmd.Run()
	if err != nil {
		return fmt.Errorf("Failed to execute ffmpeg cmd\n%s\nerr: %v\n", bgCmd.String(), err)
	}

	//args := []string{"-i", "output.mp4", "-vf", filters, "-codec:a", "copy", "finished-vid.mp4"}

	/*
		ffmpegCmd := exec.Command("ffmpeg", args...)
		err = ffmpegCmd.Run()
		if err != nil {
			return fmt.Errorf("Failed to execute ffmpeg cmd\n%s\nerr: %v\n", ffmpegCmd.String(), err)
		}
	*/

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
	//box = `box=1:boxcolor=black@0.5:boxborderw=5:`
	xPos       = `x=20:`
	yPos       = `y=800`
	slideSpeed = float64(2000)
)

func generateOverlayWithFadeArgs(contentObjs []*Content, args config.Overlay) (allFilters string) {
	fade := args.Fade
	duration := args.Duration
	font := fmt.Sprintf("drawtext=fontfile=%s:", args.Font)
	fontColor := fmt.Sprintf("fontcolor=%s:", args.Color)
	fontSize := fmt.Sprintf("fontsize=%s:", args.Size)

	var cursor float64
	for i, v := range contentObjs {
		title := formatOverlayString(v.Title)
		creator := formatOverlayString(v.CreatorName)
		overlayText := fmt.Sprintf("text=%s\n%s:", title, creator)
		fmt.Printf("\nAppling overlay text:\n%s\n", overlayText)
		alpha := fmt.Sprintf(`alpha='if(lt(t,%f),0,if(lt(t,%f),(t-%f)/1,if(lt(t,%f),1,if(lt(t,%f),(1-(t-%f))/1,0))))':`, cursor+0.5, cursor+0.5+fade, cursor+0.5, cursor+duration, cursor+duration+fade, cursor+duration)
		fullOverlay := font + overlayText + fontSize + fontColor + alpha + xPos + yPos
		allFilters += fullOverlay
		if i < len(contentObjs)-1 {
			allFilters += ","
		}

		cursor += v.Duration
	}

	return allFilters
}

func generateOverlayBackground(contentObjs []*Content, args config.Overlay) (bgFilter string) {
	duration := args.Duration + 0.5
	ypos := `y=780`

	var cursor float64
	for i, v := range contentObjs {
		// calculates a rough estimate for bg length based on content title
		var bgLength float64 = float64(len(v.Title) * 19)
		var slide float64 = bgLength / slideSpeed
		//TODO: implement max distance, generate distance from length of title
		bgFilter += fmt.Sprintf(`overlay=x='if(lt(t,%f),NAN,if(lt(t,%f),-w+(t-%f)*%f,if(lt(t,%f),-w+%f,-w+%f-(t-%f)*%f)))':%s`, cursor, cursor+slide, cursor, slideSpeed, cursor+slide+duration, slide*slideSpeed, slide*slideSpeed, cursor+slide+duration, slideSpeed, ypos)
		if i < len(contentObjs)-1 {
			bgFilter += ","
		}

		cursor += v.Duration
	}
	fmt.Println("Bg filter: ", bgFilter)

	return bgFilter
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

func formatOverlayString(raw string) string {
	formatted := strings.ReplaceAll(raw, `'`, `\\\'`)
	return strings.ReplaceAll(formatted, `,`, `\\\,`)
}
