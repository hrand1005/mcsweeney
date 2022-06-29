package video

import (
	ffmpeg "github.com/u2takey/ffmpeg-go"
	"log"
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
	inputStreams := make([]*ffmpeg.Stream, 0, len(v.inputs))
	for _, v := range v.inputs {
		inputStreams = append(inputStreams, ffmpeg.Input(v))
	}
	return ffmpeg.Concat(inputStreams).Output(outfile, ffmpeg.KwArgs{"c:v": "libx265"}).OverWriteOutput().ErrorToStdOut().Run()
}
