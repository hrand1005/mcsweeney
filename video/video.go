package video

import (
	ffmpeg "github.com/u2takey/ffmpeg-go"
)

type video struct {
	file string
}

func New(file string) video {
	return video{
		file: file,
	}
}

func (v video) WriteToFile(outfile string) error {
	return ffmpeg.Input(v.file).Output(outfile, ffmpeg.KwArgs{"c:v": "libx265"}).OverWriteOutput().ErrorToStdOut().Run()
}
