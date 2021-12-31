package content_test

import (
	"mcsweeney/content"
	"testing"
)

// TestDescriptionGeneratorVisitIntro calls DescriptionGenerator.VisitIntro and
// provides an Intro element, checking that the DescriptionGenerator's String()
// method returns the expected string representation.
func TestDescriptionGeneratorVisitIntro(t *testing.T) {
	tests := []struct {
		name    string
		intro   *content.Intro
		visitor *content.DescriptionGenerator
		want    string
	}{
		{
			name:    "Nominal description generation for intro.",
			intro:   &content.Intro{Description: "Intro description."},
			visitor: &content.DescriptionGenerator{},
			want:    "Intro description.",
		},
		{
			name:    "No description.",
			intro:   &content.Intro{},
			visitor: &content.DescriptionGenerator{},
			want:    "",
		},
	}
	for _, tc := range tests {
		tc.visitor.VisitIntro(tc.intro)
		if tc.visitor.String() != tc.want {
			t.Fatalf("Got: %s\nWanted: %s", tc.visitor.String(), tc.want)
		}
	}
}

// TestDescriptionGeneratorVisitOutro calls DescriptionGenerator.VisitOutro and
// provides an Outro element, checking that the DescriptionGenerator's String()
// method returns the expected string representation.
func TestDescriptionGeneratorVisitOutro(t *testing.T) {
	tests := []struct {
		name    string
		outro   *content.Outro
		visitor *content.DescriptionGenerator
		want    string
	}{
		{
			name:    "Nominal description generation for outro.",
			outro:   &content.Outro{Description: "Outro description."},
			visitor: &content.DescriptionGenerator{},
			want:    "Outro description.",
		},
		{
			name:    "No description.",
			outro:   &content.Outro{},
			visitor: &content.DescriptionGenerator{},
			want:    "",
		},
	}
	for _, tc := range tests {
		tc.visitor.VisitOutro(tc.outro)
		if tc.visitor.String() != tc.want {
			t.Fatalf("Got: %s\nWanted: %s", tc.visitor.String(), tc.want)
		}
	}
}

// TestDescriptionGeneratorVisitClip calls DescriptionGenerator.VisitClip and
// provides an Clip element, checking that the DescriptionGenerator's String()
// method returns the expected string representation.
func TestDescriptionGeneratorVisitClip(t *testing.T) {
	tests := []struct {
		name    string
		clip    *content.Clip
		visitor *content.DescriptionGenerator
		want    string
	}{
		{
			name: "Nominal description generation for single clip.",
			clip: &content.Clip{
				Author:      "TestAuthor",
				Broadcaster: "TestBroadcaster",
				Duration:    1.0,
				Title:       "Test Title",
			},
			visitor: &content.DescriptionGenerator{},
			want:    "\n\n[0:00] 'Test Title'\nStreamed by TestBroadcaster at \nClipped by TestAuthor\n",
		},
		{
			name:    "Empty clip.",
			clip:    &content.Clip{},
			visitor: &content.DescriptionGenerator{},
			want:    "",
		},
	}
	for _, tc := range tests {
		tc.visitor.VisitClip(tc.clip)
		if tc.visitor.String() != tc.want {
			t.Fatalf("Got: %s\nWanted: %s", tc.visitor.String(), tc.want)
		}
	}
}
