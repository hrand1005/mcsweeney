package main

import (
	"fmt"
	"os"
	"os/exec"
)

// MKVToMP4Concatenator maintains an ordering of mkv files to be
// concatenated into one composite mp4 file.
type MKVToMP4Concatenator struct {
	inputFiles []string
}

// NewMKVToMP4Concatenator returns a blank concatenator with no input files.
func NewMKVToMP4Concatenator() *MKVToMP4Concatenator {
	return &MKVToMP4Concatenator{}
}

// AppendMKVFile appends the given file during concatenation.
// The provided file is expected to have the *.mkv extension
func (c *MKVToMP4Concatenator) AppendMKVFile(f string) {
	c.inputFiles = append(c.inputFiles, f)
}

// PrependMKVFile prepends the given file during concatenation.
// The provided file is expected to have the *.mkv extension
func (c *MKVToMP4Concatenator) PrependMKVFile(f string) {
	c.inputFiles = append([]string{f}, c.inputFiles...)
}

// Concatenate takes in an outfile and concatenates the input mkv files to a
// single mp4 file. The provided outfile is expected to have the *.mp4 extension
func (c *MKVToMP4Concatenator) Concatenate(outfile string) error {
	f, err := os.Create(listFile)
	if err != nil {
		return err
	}
	defer f.Close()

	for _, v := range c.inputFiles {
		f.WriteString(
			fmt.Sprintf("file '%s'\n", v),
		)
	}

	return ConcatMKVFromFileToMP4(listFile, outfile)

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
