package content

import (
	"errors"
	"fmt"
	"strings"
)

// ErrEmptyOverlay is returned when an overlay with no defined background
// attempts to overlay an element with no title or broadcaster.
var ErrEmptyOverlay = errors.New("No image or text available to overlay element.")

// Overlayer generates video overlays for the video components it visits.
type Overlayer struct {
	Font       string
	Background string
	args       []string
	cursor     float64
	overlay    []string
}

// VisitIntro implements the visitor interface for Overlayer.
// Increments internal cursor, but does not alter the overlay string.
func (o *Overlayer) VisitIntro(i *Intro) error {
	o.cursor += i.Duration
	return nil
}

// VisitOutro implements the visitor interface for Overlayer.
// Increments internal cursor, but does not alter the overlay string.
func (o *Overlayer) VisitOutro(u *Outro) error {
	o.cursor += u.Duration
	return nil
}

// These constants are restrictive, but simplify the client's interface to
// create overlays.
const (
	// define constants for overlay background
	OverlayDuration float64 = 3
	SlideSpeed      float64 = 2000
	XPosition       int     = 20
	// define constants specific to the text overlay and fade
	FontColor string  = "ffffff"
	FontSize  int     = 52
	TextFade  float64 = 0.5
	TextDelay float64 = 0.5
)

// VisitClip implements the visitor interface for Overlayer. It generates
// an overlay for the clip using the clips Title and Broadcaster fields, as well
// as a number of configurable options for the Overlayer instance.
func (o *Overlayer) VisitClip(c *Clip) error {
	// overlays can't be applied to clips with no duration
	if c.Duration == 0.0 {
		return ErrNoDuration
	}
	// create background and text for overlay, at least one must be valid
	bgOverlay := o.createOverlayBackground(c)
	textOverlay := o.createOverlayText(c)
	if bgOverlay == "" && textOverlay == "" {
		return ErrEmptyOverlay
	}
	// update input args, overlay string, and cursor
	o.args = append(o.args, "-i", o.Background)
	o.overlay = append(o.overlay, bgOverlay+textOverlay)
	o.cursor += c.Duration

	return nil
}

// String returns a string of the generated overlay. The overlay is the
// aggregate of all visited elements, also reflecting visit order.
func (o *Overlayer) String() string {
	if len(o.overlay) == 0 {
		return ""
	}
	// add filter complex argument before exporting
	argsString := strings.Join(o.args, " ") + " -filter_complex "
	overlayString := strings.Join(o.overlay, ",")

	return argsString + overlayString
}

// Slice returns a slice representation of the generated overlay.
func (o *Overlayer) Slice() []string {
	if len(o.args) == 0 {
		return nil
	}
	// add filter complex argument before exporting
	args := append(o.args, "-filter_complex")
	overlayString := strings.Join(o.overlay, ",")

	return append(args, overlayString)
}

func (o *Overlayer) createOverlayBackground(c *Clip) string {
	if o.Background == "" {
		return ""
	}
	// determine the length of the overlay's background, multiply by arbitrary
	// size coefficient (16), plus x offset to create margins for text
	bgLen := float64(max(len(c.Title), len(c.Broadcaster)))*32.0 + 3.5*float64(XPosition)
	// determine the duration of the background's slide animation
	sDur := bgLen / SlideSpeed
	// generate strings to animate the background x and y position over time
	yString := `y=(H-h)/2`
	xString := fmt.Sprintf(`x='if(lt(t,%f),NAN,if(lt(t,%f),-w+(t-%f)*%f,if(lt(t,%f),-w+%f,-w+%f-(t-%f)*%f)))'`, o.cursor, o.cursor+sDur, o.cursor, SlideSpeed, o.cursor+sDur+OverlayDuration, bgLen, bgLen, o.cursor+sDur+OverlayDuration, SlideSpeed)

	return fmt.Sprintf(`overlay=%s:%s,`, xString, yString)
}

func (o *Overlayer) createOverlayText(c *Clip) string {
	if c.Title == "" && c.Broadcaster == "" {
		return ""
	}
	// generate arguments for ffmpeg filter
	fontString := fmt.Sprintf(`drawtext=fontfile=%s`, o.Font)
	textString := fmt.Sprintf("text=%s\n%s", escapeText(c.Title), escapeText(c.Broadcaster))
	sizeString := fmt.Sprintf(`fontsize=%v`, FontSize)
	colorString := fmt.Sprintf(`fontcolor=%s`, FontColor)
	// generate string to fade text in over background with a delay
	fadeString := fmt.Sprintf(`alpha='if(lt(t,%f),0,if(lt(t,%f),(t-%f)/1,if(lt(t,%f),1,if(lt(t,%f),(1-(t-%f))/1,0))))'`, o.cursor+TextDelay, o.cursor+TextDelay+TextFade, o.cursor+TextDelay, o.cursor+OverlayDuration-TextDelay, o.cursor+OverlayDuration-TextDelay+TextFade, o.cursor+OverlayDuration-TextDelay)
	xString := fmt.Sprintf(`x=%v`, XPosition)
	yString := `y=(h-text_h)/2`

	return fmt.Sprintf(`%s:%s:%s:%s:%s:%s:%s`, fontString, textString, sizeString, colorString, fadeString, xString, yString)
}

func max(a, b int) int {
	if a >= b {
		return a
	}
	return b
}

func escapeText(s string) string {
	s = strings.ReplaceAll(s, `'`, `\\\'`)
	s = strings.ReplaceAll(s, `"`, `\\\"`)
	s = strings.ReplaceAll(s, `:`, `\\\:`)
	return strings.ReplaceAll(s, `,`, `\\\,`)
}
