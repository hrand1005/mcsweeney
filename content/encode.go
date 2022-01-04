package content

import (
	"fmt"
	"os"
	"os/exec"
	//"strings"
)

type Encoder struct {
	Path    string
	started bool
	count   int
}

// VisitClip implements the visitor interface for Encoder.
// Encodes the given file using ffmpeg to consistent format, and writes the path
// of the outfile to the filepath e.Path.
func (e *Encoder) VisitClip(c *Clip) {
	if c.Path == "" {
		return
	}
	outfile := e.createOutfile(c.Path)
	e.writeToListfile(outfile)
	encode(c.Path, outfile)
	return
}

// VisitIntro implements the visitor interface for Encoder.
// Encodes the given file using ffmpeg to consistent format, and writes the path
// of the outfile to the filepath e.Path.
func (e *Encoder) VisitIntro(i *Intro) {
	if i.Path == "" {
		return
	}
	outfile := e.createOutfile(i.Path)
	e.writeToListfile(outfile)
	encode(i.Path, outfile)
	return
}

// VisitOutro implements the visitor interface for Encoder.
// Encodes the given file using ffmpeg to consistent format, and writes the path
// of the outfile to the filepath e.Path.
func (e *Encoder) VisitOutro(o *Outro) {
	if o.Path == "" {
		return
	}
	outfile := e.createOutfile(o.Path)
	e.writeToListfile(outfile)
	encode(o.Path, outfile)
	return
}

func (e *Encoder) createOutfile(p string) (outfile string) {
	outfile = fmt.Sprintf("%v.mkv", e.count)
	e.count++
	return
}

func (e *Encoder) writeToListfile(s string) {
	// creates file if not exists, else opens existing file
	f, _ := os.OpenFile(e.Path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	defer f.Close()

	w := fmt.Sprintf("file '%s'\n", s)
	f.WriteString(w)
	return
}

const (
	//ffmpegEncoder string = "libx264"
	ffmpegPreset string = "slow"
)

func encode(infile string, outfile string) {
	cmd := exec.Command("ffmpeg", "-i", infile, "-preset", ffmpegPreset, "-crf", "10", "-c:a", "copy", outfile)
	fmt.Printf("About to encode...\n%s\n", cmd.String())
	cmd.Run()
}
