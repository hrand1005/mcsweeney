package content_test

import (
	"mcsweeney/content"
	"testing"
)

// TestDescriberVisitIntro calls Describer.VisitIntro and
// provides an Intro element, checking that the Describer's String()
// method returns the expected string representation.
func TestDescriberVisitIntro(t *testing.T) {
	tests := []struct {
		name    string
		intro   *content.Intro
		visitor *content.Describer
		want    string
	}{
		{
			name:    "Nominal description generation for intro.",
			intro:   &content.Intro{Description: "Intro description."},
			visitor: &content.Describer{},
			want:    "Intro description.",
		},
		{
			name:    "No description.",
			intro:   &content.Intro{},
			visitor: &content.Describer{},
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

// TestDescriberVisitOutro calls Describer.VisitOutro and
// provides an Outro element, checking that the Describer's String()
// method returns the expected string representation.
func TestDescriberVisitOutro(t *testing.T) {
	tests := []struct {
		name    string
		outro   *content.Outro
		visitor *content.Describer
		want    string
	}{
		{
			name:    "Nominal description generation for outro.",
			outro:   &content.Outro{Description: "Outro description."},
			visitor: &content.Describer{},
			want:    "Outro description.",
		},
		{
			name:    "No description.",
			outro:   &content.Outro{},
			visitor: &content.Describer{},
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

// TestDescriberVisitClip calls Describer.VisitClip and
// provides an Clip element, checking that the Describer's String()
// method returns the expected string representation.
func TestDescriberVisitClip(t *testing.T) {
	tests := []struct {
		name    string
		clip    *content.Clip
		visitor *content.Describer
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
			visitor: &content.Describer{},
			want:    "\n\n[0:00] 'Test Title'\nStreamed by TestBroadcaster at \nClipped by TestAuthor\n",
		},
		{
			name:    "Empty clip.",
			clip:    &content.Clip{},
			visitor: &content.Describer{},
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

// TestDescriberVisitMany calls multiple visit methods on
// Describer in sequence and checks that the resulting string
// returned by String() properly represents the visited element sequence.
func TestDescriberVisitMany(t *testing.T) {
	tests := []struct {
		name    string
		intros  []*content.Intro
		clips   []*content.Clip
		outros  []*content.Outro
		visitor *content.Describer
		want    string
	}{
		{
			name: "Description generation for Intro, three clips, and outro.",
			intros: []*content.Intro{
				&content.Intro{
					Description: "Intro description.",
					Duration:    4.0,
				},
			},
			clips: []*content.Clip{
				&content.Clip{
					Author:      "TestAuthor",
					Broadcaster: "TestBroadcaster",
					Duration:    1.0,
					Title:       "Test Title",
				},
				&content.Clip{},
				&content.Clip{
					Author:      "TestTwitchAuthor",
					Broadcaster: "TestTwitchBroadcaster",
					Duration:    20.0,
					Platform:    content.TWITCH,
					Title:       "Test Twitch Title",
				},
			},
			outros: []*content.Outro{
				&content.Outro{
					Description: "Outro description.",
					Duration:    3.0,
				},
			},
			visitor: &content.Describer{},
			want: "Intro description." +
				"\n\n[0:04] 'Test Title'\nStreamed by TestBroadcaster at \nClipped by TestAuthor\n" +
				"\n\n[0:05] 'Test Twitch Title'\nStreamed by TestTwitchBroadcaster at https://twitch.tv/TestTwitchBroadcaster\nClipped by TestTwitchAuthor\n" +
				"Outro description.",
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
		if tc.visitor.String() != tc.want {
			t.Fatalf("Got: %s\nWanted: %s", tc.visitor.String(), tc.want)
		}
	}
}
