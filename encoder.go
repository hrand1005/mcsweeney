package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
)

const listFile = "intermediate.txt"

type MP4ToMKVEncoder struct {
	workerPool chan string
}

// NewMP4ToMKVEncoder returns an MP4ToMKVEncoder with (todo) the capability of
// encoding the given capacity concurrently.
func NewMP4ToMKVEncoder(capacity int) *MP4ToMKVEncoder {
	return &MP4ToMKVEncoder{
		workerPool: make(chan string, capacity),
	}
}

// EncodeReport contains the results of an attempted video encoding.
type EncodeReport struct {
	Input  string
	Output string
	Err    error
}

// Encode begins a goroutine with the given input and done channels to
// asynchronously accept and  encode mp4s to mkvs.
// Returns an EncodeReport channel which notifies the client when an encoding
// is complete. It is the responsibility of the client to send the done signal
// when encoding is complete.
func (e *MP4ToMKVEncoder) Encode(inputFile <-chan string, done <-chan bool) <-chan EncodeReport {
	reportChan := make(chan EncodeReport)
	go func() {
		// increment to create unique file names
		i := 0
		for {
			select {
			case infile := <-inputFile:
				outfile := fmt.Sprintf("intermediate%v.mkv", i)

				// ensure that encoding happens only when there is
				// space in the worker pool
				e.workerPool <- infile
				go func(in, out string) {
					e.encode(infile, outfile, reportChan)
					<-e.workerPool
				}(infile, outfile)

			case <-done:
				log.Printf("Recieved done signal, exiting...")
				return
			}
			i += 1
		}
	}()

	return reportChan
}

// encode encodes the given mp4File string to the mkvFile output, generates an EncodeReport,
// and pushes it to the report channel.
func (e *MP4ToMKVEncoder) encode(mp4File, mkvFile string, reportChan chan<- EncodeReport) {
	log.Printf("Encoding %s to %s...\n", mp4File, mkvFile)
	cmd := exec.Command("ffmpeg", "-i", mp4File, "-c:v", "libx264", "-preset", "slow", "-crf", "22", "-c:a", "ac3", mkvFile, "-y")
	// DEBUGGING:
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	if err != nil {
		reportChan <- EncodeReport{
			Input: mp4File,
			Err:   fmt.Errorf("FAIL: encountered error encoding %s: %v\n", mp4File, err),
		}
	} else {
		reportChan <- EncodeReport{
			Input:  mp4File,
			Output: mkvFile,
		}
	}
}
