package content

import (
	"fmt"
	"math"
	"mcsweeney/config"
	"os"
	"os/exec"
	"strings"
	"sync"
	"time"
)

func (c *Content) ApplyOverlay(contentObjs []*Content, options config.Options) error {
	// generate strings for overlay background and overlay text filters
	bgFilters := generateOverlayBackground(contentObjs, options.Overlay)
	tFilters := generateOverlayWithFadeArgs(contentObjs, options.Overlay)

	// create ffmpeg command, using complex filter args for overlay and text
	var overlayCount int = len(contentObjs)
	// TODO: fix this
	if contentObjs[0].Title == "Intro" {
		overlayCount -= 1
	}
	bargs := make([]string, 0, overlayCount*2+4)
	bargs = append(bargs, "-i", c.Path)
	for i := 0; i < overlayCount; i++ {
		bargs = append(bargs, "-i", options.Overlay.Background)
	}
	bargs = append(bargs, "-filter_complex", bgFilters+","+tFilters, "finished-vid.mp4")
	bgCmd := exec.Command("ffmpeg", bargs...)
	err := bgCmd.Run()
	if err != nil {
		return fmt.Errorf("Failed to execute ffmpeg cmd\n%s\nerr: %v\n", bgCmd.String(), err)
	}

	// update content path
	c.Path = "finished-vid.mp4"

	return nil
}

// TODO: decouple encoding from concatenation step
// Concatenate takes a slice of Content objects, encodes them consistently and then
// concatenates them to create a new content object. As part of this, it credits
// it's subobjects in the description.
func Concatenate(contentObjs []*Content, outfile string) (*Content, error) {
	start := time.Now()
	f, err := os.Create("compile.txt")
	if err != nil {
		return nil, err
	}
	defer f.Close()

	var wg sync.WaitGroup
	encodeFunc := func(v *Content) {
		defer wg.Done()
		//filename := strings.ReplaceAll(v.Url, ".", "")
		path := "tmp/" + strings.ReplaceAll(v.Url, "/", "") + ".mkv"
		cmd := exec.Command("ffmpeg", "-i", v.Url, "-c:v", "libx264", "-preset", "slow", "-crf", "22", "-c:a", "copy", path)
		fmt.Println("Encoding content to ", path)
		/*err := */ cmd.Run()
		/*
		   if err != nil {
		       return nil, fmt.Errorf("Failed to download and encode to path %s: %v\n", v.Path, err)
		   }*/
		v.Path = path
	}

	for _, v := range contentObjs {
		wg.Add(1)
		go encodeFunc(v)
	}

	wg.Wait()

	// write path names to file
	for _, v := range contentObjs {
		w := fmt.Sprintf("file '%s'\n", v.Path)
		_, err = f.WriteString(w)
		if err != nil {
			return nil, err
		}
	}

	fmt.Printf("Time to concatenate: %v", time.Since(start))

	fmt.Println("Concatenating content...")
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
	// TODO: make x and y relative to video format
	xPos       = 20
	yPos       = 500
	slideSpeed = float64(2000)
)

func generateOverlayWithFadeArgs(contentObjs []*Content, args config.Overlay) (allFilters string) {
	fade := args.Fade
	duration := args.Duration
	font := fmt.Sprintf("drawtext=fontfile=%s:", args.Font)
	fontColor := fmt.Sprintf("fontcolor=%s:", args.Color)
	fontSize := fmt.Sprintf("fontsize=%s:", args.Size)
	xPosition := fmt.Sprintf(`x=%v:`, xPos)
	yPosition := fmt.Sprintf(`y=%v`, yPos)

	var cursor float64
	for i, v := range contentObjs {
		if v.Title != "Intro" {
			title := formatOverlayString(v.Title)
			creator := formatOverlayString(v.CreatorName)
			overlayText := fmt.Sprintf("text=%s\n%s:", title, creator)
			fmt.Printf("\nAppling overlay text:\n%s\n", overlayText)
			alpha := fmt.Sprintf(`alpha='if(lt(t,%f),0,if(lt(t,%f),(t-%f)/1,if(lt(t,%f),1,if(lt(t,%f),(1-(t-%f))/1,0))))':`, cursor+0.5, cursor+0.5+fade, cursor+0.5, cursor+duration, cursor+duration+fade, cursor+duration)
			fullOverlay := font + overlayText + fontSize + fontColor + alpha + xPosition + yPosition
			allFilters += fullOverlay
			if i < len(contentObjs)-1 {
				allFilters += ","
			}
		}

		cursor += v.Duration
	}

	return allFilters
}

func generateOverlayBackground(contentObjs []*Content, args config.Overlay) (bgFilter string) {
	duration := args.Duration + 0.5
	yPosition := fmt.Sprintf(`y=%v`, yPos-20)

	var cursor float64
	for i, v := range contentObjs {
		if v.Title != "Intro" {
			// calculates a rough estimate for bg length based on content title
			// base on font size?
			tLength := float64(len(v.Title) * 14)
			cLength := float64(len(v.CreatorName) * 14)
			bgLength := math.Max(tLength, cLength)
			slide := bgLength / slideSpeed
			bgFilter += fmt.Sprintf(`overlay=x='if(lt(t,%f),NAN,if(lt(t,%f),-w+(t-%f)*%f,if(lt(t,%f),-w+%f,-w+%f-(t-%f)*%f)))':%s`, cursor, cursor+slide, cursor, slideSpeed, cursor+slide+duration, slide*slideSpeed, slide*slideSpeed, cursor+slide+duration, slideSpeed, yPosition)
			if i < len(contentObjs)-1 {
				bgFilter += ","
			}
		}

		cursor += v.Duration
	}

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
	formatted = strings.ReplaceAll(formatted, `:`, `\\\:`)
	return strings.ReplaceAll(formatted, `,`, `\\\,`)
}
