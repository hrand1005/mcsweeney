// Video defines utilities for video processing.
package video

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"sync"
)

// define ffmpeg filter arguments to be used in concatentation
const (
	defaultAspectRatio = "16/9"
	defaultFrameRate   = "60"
	defaultScale       = "1920x1080"
)

// DEPRECATED DUE TO INFLEXIBLE CONCATENATION OF DIFFERENT CODECS / RESOLUTIONS / SCALES
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
		if err := EncodeMP4ToMKV(v, intFile); err != nil {
			log.Printf("Encountered error executing ffmpeg command: %v\n", err)
		}
		f.WriteString(
			fmt.Sprintf("file '%s'\n", intFile),
		)
	}

	return ConcatMKVFromFileToMP4(intermediateFileList, outfile)
}

// EncodeToMKV encodes the given mp4 file to mkv intermediate file defined by
// outfile, which is expected to have the .mkv extension in the name. Returns the
// result of the command execution.
func EncodeMP4ToMKV(input string, outfile string) error {
	cmd := exec.Command("ffmpeg", "-i", input, "-c:v", "libx264", "-preset", "slow", "-crf", "22", "-c:a", "ac3", outfile, "-y")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return cmd.Run()
}

// ConcatInfileToMP4 takes in an ffmpeg input file for concatenation, and outputs the
// result to outfile, which is expected to have the .mp4 file extension in the name.
// Returns the result of the command execution.
func ConcatMKVFromFileToMP4(infile, outfile string) error {
	cmd := exec.Command("ffmpeg", "-f", "concat", "-safe", "0", "-i", infile, "-c", "copy", outfile, "-y")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return cmd.Run()
}

type MP4ToMKVEncoder struct {
	// .txt file listing all encoded mp4s
	mu       sync.Mutex
	listFile string
}

// TODO: enforce encoding capacity through job pool or buffered channels
func NewMP4ToMKVEncoder(listFile string, encodeCapacity int) *MP4ToMKVEncoder {
	return &MP4ToMKVEncoder{
		listFile: listFile,
	}
}

func (e *MP4ToMKVEncoder) Encode(inputFile <-chan string, report chan<- string, done <-chan bool) {
	i := 0
	for {
		select {
		case infile := <-inputFile:
			outfile := fmt.Sprintf("intermediate%v.mkv", i)
			go func(mp4File, mkvFile string) {
				log.Printf("Encoding %s to %s...\n", mp4File, mkvFile)
				cmd := exec.Command("ffmpeg", "-i", mp4File, "-c:v", "libx264", "-preset", "slow", "-crf", "22", "-c:a", "ac3", mkvFile, "-y")
				cmd.Stdout = os.Stdout
				cmd.Stderr = os.Stderr
				err := cmd.Run()
				if err != nil {
					report <- fmt.Sprintf("FAIL: encountered error encoding %s: %v\n", mp4File, err)
				} else {
					// TODO: Enforce a consistent ordering, or report ordering of clips written to files
					report <- fmt.Sprintf("SUCCESS: encoded %s to %s\n", mp4File, mkvFile)
					e.mu.Lock()
					defer e.mu.Unlock()
					WriteStringToFile(fmt.Sprintf("file '%s'\n", mkvFile), e.listFile)
				}
			}(infile, outfile)
		case <-done:
			log.Printf("Recieved done signal, exiting...")
			return
		}
		i += 1
	}
}

func WriteStringToFile(s, filename string) {
	f, _ := os.OpenFile(filename, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0600)
	f.WriteString(s)
}
