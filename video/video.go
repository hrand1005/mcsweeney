package video

import (
	"fmt"
	"log"
	"os"
	"os/exec"
)

type Video interface {
	Append(string)
	WriteToFile(string) error
}

type video struct {
	inputs []string
}

func New() *video {
	return &video{}
}

func (v *video) Append(input string) {
	v.inputs = append(v.inputs, input)
}

func (v *video) WriteToFile(outfile string) error {
	log.Printf("Concating input files:\n%v\n", v.inputs)
	args := make([]string, 0, 10)
	complexFilterString := ""
	for i, v := range v.inputs {
		args = append(args, "-i", v)
		complexFilterString += fmt.Sprintf("[%v:v:0][%v:a:0]", i, i)
	}
	complexFilterString += fmt.Sprintf("concat=n=%v:v=1:a=1[outv][outa]", len(v.inputs))
	args = append(args, "-filter_complex", complexFilterString, "-map", "[outv]", "-map", "[outa]", outfile, "-y")
	cmd := exec.Command("ffmpeg", args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	log.Printf("EXECUTING COMMAND:\n%s\n", cmd.String())

	return cmd.Run()
}
