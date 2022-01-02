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
