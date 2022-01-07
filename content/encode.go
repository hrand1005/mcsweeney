package content

import (
	"fmt"
	"os"
	"os/exec"
	"sync"
)

// Encoder defines a type for encoding video components.
type Encoder struct {
	Path    string
	started bool
	count   int
	Wg      sync.WaitGroup
}

// VisitIntro implements the visitor interface for Encoder.
// Encodes the given file using ffmpeg to consistent format, and writes the path
// of the outfile to the filepath e.Path.
func (e *Encoder) VisitIntro(i *Intro) error {
	if i.Path == "" {
		return ErrEmptyPath
	}
	outfile := e.createOutfile(i.Path)
	err := e.writeToListfile(outfile)
	if err != nil {
		return err
	}

	e.Wg.Add(1)
	go e.encode(i.Path, outfile)
	return nil
}

// VisitOutro implements the visitor interface for Encoder.
// Encodes the given file using ffmpeg to consistent format, and writes the path
// of the outfile to the filepath e.Path.
func (e *Encoder) VisitOutro(o *Outro) error {
	if o.Path == "" {
		return ErrEmptyPath
	}
	outfile := e.createOutfile(o.Path)
	err := e.writeToListfile(outfile)
	if err != nil {
		return err
	}

	e.Wg.Add(1)
	go e.encode(o.Path, outfile)
	return nil
}

// VisitClip implements the visitor interface for Encoder.
// Encodes the given file using ffmpeg to consistent format, and writes the path
// of the outfile to the filepath e.Path.
func (e *Encoder) VisitClip(c *Clip) error {
	if c.Path == "" {
		return ErrEmptyPath
	}
	outfile := e.createOutfile(c.Path)
	err := e.writeToListfile(outfile)
	if err != nil {
		return err
	}

	e.Wg.Add(1)
	go e.encode(c.Path, outfile)
	return nil
}

func (e *Encoder) createOutfile(p string) (outfile string) {
	outfile = fmt.Sprintf("%v.mkv", e.count)
	e.count++
	return
}

func (e *Encoder) writeToListfile(s string) error {
	// creates file if not exists, else opens existing file
	f, err := os.OpenFile(e.Path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer f.Close()

	// write a string to the file in ffmpeg concat format
	w := fmt.Sprintf("file '%s'\n", s)
	_, err = f.WriteString(w)
	if err != nil {
		return err
	}

	return nil
}

const (
	ffmpegPreset string = "slow"
)

func (e *Encoder) encode(infile, outfile string) error {
	defer e.Wg.Done()
	cmd := exec.Command("ffmpeg", "-i", infile, "-preset", ffmpegPreset, "-crf", "10", "-c:a", "copy", outfile)
	fmt.Printf("Encoding:\n%s\n", cmd.String())
	cmd.Run()
	/*err := cmd.Run()
	if err != nil {
		return err
	}*/

	return nil
}
