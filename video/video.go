// Video defines utilities for video processing.
package video

import (
	"fmt"
	"log"
	"os"
	"os/exec"
)

// ConcatenateMP4Files takes a slice of MP4 files as input and writes them to
// a file with the given 'outfile' name
func ConcatenateMP4Files(inputs []string, outfile string) error {
	args := make([]string, 0, 10)
	complexFilterString := ""
	for i, v := range inputs {
		args = append(args, "-i", v)
		complexFilterString += fmt.Sprintf("[%v:v:0][%v:a:0]", i, i)
	}
	complexFilterString += fmt.Sprintf("concat=n=%v:v=1:a=1[outv][outa]", len(inputs))
	args = append(args, "-filter_complex", complexFilterString, "-map", "[outv]", "-map", "[outa]", outfile, "-y")
	cmd := exec.Command("ffmpeg", args...)

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	log.Printf("EXECUTING COMMAND:\n%s\n", cmd.String())

	return cmd.Run()
}
