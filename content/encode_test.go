package content_test

import (
	"bufio"
	"errors"
	"mcsweeney/content"
	"os"
	"strings"
	"testing"
)

const (
	EncoderPath string = "encoded.txt"
)

// TestEncoderPath calls encoder's visit methods without initializing a Path in
// the Encoder, checking that the appropriate error is returned.
func TestEncoderEmpty(t *testing.T) {
	encoder := &content.Encoder{}
	err := encoder.VisitIntro(&content.Intro{})
	if err != content.ErrEmptyPath {
		t.Fatalf("Got: %s, Wanted: %s\n", err, content.ErrEmptyPath)
	}
	err = encoder.VisitClip(&content.Clip{})
	if err != content.ErrEmptyPath {
		t.Fatalf("Got: %s, Wanted: %s\n", err, content.ErrEmptyPath)
	}
	err = encoder.VisitOutro(&content.Outro{})
	if err != content.ErrEmptyPath {
		t.Fatalf("Got: %s, Wanted: %s\n", err, content.ErrEmptyPath)
	}
}

// TestEncoderVisitIntro calls Encoder.VisitIntro and provides an Intro element,
// checking that the outfile is written to the file returned by Path, and that
// the contents of the file are correct.
func TestEncoderVisitIntro(t *testing.T) {
	tests := []struct {
		name    string
		intro   *content.Intro
		visitor *content.Encoder
		want    string
		wantErr error
	}{
		{
			name:    "Nominal encoding for intro.",
			intro:   &content.Intro{Path: "fakes/intro.mp4"},
			visitor: &content.Encoder{Path: EncoderPath},
			want:    "file '0.mkv'",
			wantErr: nil,
		},
		{
			name:    "No encoding.",
			intro:   &content.Intro{},
			visitor: &content.Encoder{Path: EncoderPath},
			want:    "",
			wantErr: os.ErrNotExist,
		},
	}
	for _, tc := range tests {
		tc.visitor.VisitIntro(tc.intro)

		// check that outfile can be opened
		f, err := os.Open(tc.visitor.Path)
		if err != nil {
			if errors.Is(err, tc.wantErr) {
				break
			}
			os.Remove(tc.visitor.Path)
			t.Fatalf("Error opening file: %v", err)
		}
		defer f.Close()

		// scan file for first row contents
		s := bufio.NewScanner(f)
		s.Split(bufio.ScanLines)
		s.Scan()
		got := s.Text()

		// check that scan didn't result in an error
		if err = s.Err(); err != nil {
			os.Remove(tc.visitor.Path)
			t.Fatalf("Error scanning file: %v", err)
		}

		// check for correct file contents
		if got != tc.want {
			os.Remove(tc.visitor.Path)
			t.Fatalf("Got: %s\nWanted: %s", got, tc.want)
		}
		os.Remove(tc.visitor.Path)
	}
	os.Remove(EncoderPath)
}

// TestEncoderVisitOutro calls Encoder.VisitOutro and provides an Outro element,
// checking that the outfile is written to the file returned by Path, and that
// the contents of the file are correct.
func TestEncoderVisitOutro(t *testing.T) {
	tests := []struct {
		name    string
		outro   *content.Outro
		visitor *content.Encoder
		want    string
		wantErr error
	}{
		{
			name:    "Nominal encoding for outro.",
			outro:   &content.Outro{Path: "fakes/outro.mp4"},
			visitor: &content.Encoder{Path: EncoderPath},
			want:    "file '0.mkv'",
			wantErr: nil,
		},
		{
			name:    "No encoding.",
			outro:   &content.Outro{},
			visitor: &content.Encoder{Path: EncoderPath},
			want:    "",
			wantErr: os.ErrNotExist,
		},
	}
	for _, tc := range tests {
		tc.visitor.VisitOutro(tc.outro)

		// check that outfile can be opened
		f, err := os.Open(tc.visitor.Path)
		if err != nil {
			if errors.Is(err, tc.wantErr) {
				break
			}
			os.Remove(tc.visitor.Path)
			t.Fatalf("Error opening file: %v", err)
		}
		defer f.Close()

		// scan file for first row contents
		s := bufio.NewScanner(f)
		s.Split(bufio.ScanLines)
		s.Scan()
		got := strings.TrimSpace(s.Text())

		// check that scan didn't result in an error
		if err = s.Err(); err != nil {
			os.Remove(tc.visitor.Path)
			t.Fatalf("Error scanning file: %v", err)
		}

		// check for correct file contents
		if got != tc.want {
			os.Remove(tc.visitor.Path)
			t.Fatalf("Got: %s\nWanted: %s", got, tc.want)
		}
		os.Remove(tc.visitor.Path)
	}
	os.Remove(EncoderPath)
}

// TestEncoderVisitClip calls Encoder.VisitClip and provides an Clip element,
// checking that the outfile is written to the file returned by Path, and that
// the contents of the file are correct.
func TestEncoderVisitClip(t *testing.T) {
	tests := []struct {
		name    string
		clip    *content.Clip
		visitor *content.Encoder
		want    string
		wantErr error
	}{
		{
			name:    "Nominal encoding for clip.",
			clip:    &content.Clip{Path: "fakes/clip.mp4"},
			visitor: &content.Encoder{Path: EncoderPath},
			want:    "file '0.mkv'",
			wantErr: nil,
		},
		{
			name:    "No encoding.",
			clip:    &content.Clip{},
			visitor: &content.Encoder{Path: EncoderPath},
			want:    "",
			wantErr: os.ErrNotExist,
		},
	}
	for _, tc := range tests {
		tc.visitor.VisitClip(tc.clip)

		// check that outfile can be opened
		f, err := os.Open(tc.visitor.Path)
		if err != nil {
			if errors.Is(err, tc.wantErr) {
				break
			}
			os.Remove(tc.visitor.Path)
			t.Fatalf("Error opening file: %v", err)
		}
		defer f.Close()

		// scan file for first row contents
		s := bufio.NewScanner(f)
		s.Split(bufio.ScanLines)
		s.Scan()
		got := s.Text()

		// check that scan didn't result in an error
		if err = s.Err(); err != nil {
			os.Remove(tc.visitor.Path)
			t.Fatalf("Error scanning file: %v", err)
		}

		// check for correct file contents
		if got != tc.want {
			os.Remove(tc.visitor.Path)
			t.Fatalf("Got: %s\nWanted: %s", got, tc.want)
		}
		os.Remove(tc.visitor.Path)
	}
	os.Remove(EncoderPath)
}

// TestEncoderVisitMany calls a multiple visits in sequence an providing various
// components, checking that the outfile is written to the file returned by Path,
// and that the contents of the file are correct.
func TestEncoderVisitMany(t *testing.T) {
	tests := []struct {
		name    string
		intros  []*content.Intro
		clips   []*content.Clip
		outros  []*content.Outro
		visitor *content.Encoder
		want    string
		wantErr error
	}{
		{
			name: "Encode one intro, three clips, and an outro.",
			intros: []*content.Intro{
				&content.Intro{Path: "fakes/intro.mp4"},
			},
			clips: []*content.Clip{
				&content.Clip{Path: "fakes/clip.mp4"},
				&content.Clip{Path: "fakes/clip2.mp4"},
				&content.Clip{Path: "fakes/clip3.mp4"},
			},
			outros: []*content.Outro{
				&content.Outro{Path: "fakes/outro.mp4"},
			},
			visitor: &content.Encoder{Path: EncoderPath},
			want: "file '0.mkv'\n" +
				"file '1.mkv'\n" +
				"file '2.mkv'\n" +
				"file '3.mkv'\n" +
				"file '4.mkv'\n",
			wantErr: nil,
		},
	}
	for _, tc := range tests {
		// visit all elements
		for _, intro := range tc.intros {
			tc.visitor.VisitIntro(intro)
		}
		for _, clip := range tc.clips {
			tc.visitor.VisitClip(clip)
		}
		for _, outro := range tc.outros {
			tc.visitor.VisitOutro(outro)
		}

		// read file contents
		gotBytes, err := os.ReadFile(tc.visitor.Path)
		if err != nil {
			if errors.Is(err, tc.wantErr) {
				break
			}
			os.Remove(tc.visitor.Path)
			t.Fatalf("Error opening file: %v", err)
		}

		// convert bytes to string for comparison
		got := string(gotBytes)

		// check for correct file contents
		if got != tc.want {
			os.Remove(tc.visitor.Path)
			t.Fatalf("Got: %s\nWanted: %s", got, tc.want)
		}
		os.Remove(tc.visitor.Path)
	}
	os.Remove(EncoderPath)

}
