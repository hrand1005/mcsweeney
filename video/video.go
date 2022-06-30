// Video defines utilities for video processing.
package video

import (
	"fmt"
	"log"
	"os"
	"os/exec"
)

// define ffmpeg filter arguments to be used in concatentation
const (
	defaultAspectRatio = "16/9"
	defaultFrameRate   = "60"
	defaultScale       = "1920x1080"
)

// ConcatenateMP4Files takes a slice of MP4 files as input and writes them to
// a file with the given 'outfile' name. Uses default values defined above for
// concatenation.
func ConcatenateMP4Files(inputs []string, outfile string) error {
	args := make([]string, 0, 10)
	complexFilterString := ""
	complexFilterEndString := ""
	for i, v := range inputs {
		args = append(args, "-i", v)
		complexFilterString += fmt.Sprintf("[%[1]v:v:0]scale=%s,setdar=%s[v%[1]d];", i, defaultScale, defaultAspectRatio)
		complexFilterEndString += fmt.Sprintf("[v%[1]d][%[1]v:a:0]", i)
	}
	complexFilterString += complexFilterEndString + fmt.Sprintf("concat=n=%v:v=1:a=1[outv][outa]", len(inputs))
	args = append(args, "-filter_complex", complexFilterString, "-map", "[outv]", "-map", "[outa]", "-s", defaultAspectRatio, "-r", defaultFrameRate, "-c:v", "libx264", outfile, "-y")
	cmd := exec.Command("ffmpeg", args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	log.Printf("EXECUTING COMMAND:\n%s\n", cmd.String())

	return cmd.Run()
}

const intermediateFileList = "intermediate.txt"

func EncodeAndConcatMP4Files(inputs []string, outfile string) error {
	f, err := os.Create(intermediateFileList)
	if err != nil {
		return err
	}
	defer f.Close()

	for i, v := range inputs {
		intFile := fmt.Sprintf("intermediate%v.mkv", i)
		f.WriteString(
			fmt.Sprintf("file '%s'\n", intFile),
		)
		cmd := exec.Command("ffmpeg", "-i", v, "-c:v", "libx264", "-preset", "slow", "-crf", "22", "-c:a", "ac3", intFile, "-y")
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		if err := cmd.Run(); err != nil {
			log.Printf("Encountered error executing ffmpeg command: %v\n", err)
		}
	}

	cmd := exec.Command("ffmpeg", "-f", "concat", "-safe", "0", "-i", intermediateFileList, "-c", "copy", outfile, "-y")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return cmd.Run()
}
