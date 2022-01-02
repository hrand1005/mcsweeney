package content

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
)

type Encoder struct {
	Path    string
	started bool
}

// VisitClip implements the visitor interface for Encoder.
// Encodes the given file using ffmpeg to consistent format, and writes the path
// of the outfile to the filepath e.Path.
func (e *Encoder) VisitClip(c *Clip) {
	if c.Path == "" {
		return
	}
	outfile := createOutfile(c.Path)
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
	outfile := createOutfile(i.Path)
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
	outfile := createOutfile(o.Path)
	e.writeToListfile(outfile)
	encode(o.Path, outfile)
	return
}

func createOutfile(p string) (outfile string) {
	outfile = strings.ReplaceAll(p, "/", "")
	outfile = strings.TrimSuffix(outfile, ".mp4")
	return outfile + ".mkv"
}

func (e *Encoder) writeToListfile(s string) {
	// if this is the first write, then create the outfile
	var f *os.File
	if !e.started {
		f, _ = os.Create(e.Path)
		e.started = true
	}
	defer f.Close()

	w := fmt.Sprintf("file '%s'\n", s)
	f.WriteString(w)
	return
}

const (
	ffmpegEncoder string = "libx264"
	ffmpegPreset  string = "slow"
)

func encode(infile string, outfile string) {
	cmd := exec.Command("ffmpeg", "-i", infile, "-c:v", ffmpegEncoder, "-preset", ffmpegPreset, "-crf", "22", "-c:a", "copy", outfile)
	cmd.Run()
}
