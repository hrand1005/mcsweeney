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
		wantErr error
	}{
		{
			name:    "Nominal description generation for intro.",
			intro:   &content.Intro{Description: "Intro description."},
			visitor: &content.Describer{},
			want:    "Intro description.",
			wantErr: nil,
		},
		{
			name:    "No description.",
			intro:   &content.Intro{},
			visitor: &content.Describer{},
			want:    "",
			wantErr: nil,
		},
	}
	for _, tc := range tests {
		gotErr := tc.visitor.VisitIntro(tc.intro)
		if tc.visitor.String() != tc.want || gotErr != tc.wantErr {
			t.Fatalf("Got: %s\nWanted: %s\nGotErr: %s, WantErr: %s", tc.visitor, tc.want, gotErr, tc.wantErr)
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
		wantErr error
	}{
		{
			name:    "Nominal description generation for outro.",
			outro:   &content.Outro{Description: "Outro description."},
			visitor: &content.Describer{},
			want:    "Outro description.",
			wantErr: nil,
		},
		{
			name:    "No description.",
			outro:   &content.Outro{},
			visitor: &content.Describer{},
			want:    "",
			wantErr: nil,
		},
	}
	for _, tc := range tests {
		gotErr := tc.visitor.VisitOutro(tc.outro)
		if tc.visitor.String() != tc.want || gotErr != tc.wantErr {
			t.Fatalf("Got: %s\nWanted: %s\nGotErr: %s, WantErr: %s", tc.visitor, tc.want, gotErr, tc.wantErr)
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
		wantErr error
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
			wantErr: nil,
		},
		{
			name:    "Empty clip returns ErrNoDuration.",
			clip:    &content.Clip{},
			visitor: &content.Describer{},
			want:    "",
			wantErr: content.ErrNoDuration,
		},
	}
	for _, tc := range tests {
		gotErr := tc.visitor.VisitClip(tc.clip)
		if tc.visitor.String() != tc.want || gotErr != tc.wantErr {
			t.Fatalf("Got: %s\nWanted: %s\nGotErr: %s, WantErr: %s", tc.visitor, tc.want, gotErr, tc.wantErr)
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
		wantErr error
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
				&content.Clip{
					Author:      "TestAuthor2",
					Broadcaster: "TestBroadcaster2",
					Duration:    0.5,
					Title:       "Test Title 2",
				},
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
				"\n\n[0:05] 'Test Title 2'\nStreamed by TestBroadcaster2 at \nClipped by TestAuthor2\n" +
				"\n\n[0:05] 'Test Twitch Title'\nStreamed by TestTwitchBroadcaster at https://twitch.tv/TestTwitchBroadcaster\nClipped by TestTwitchAuthor\n" +
				"Outro description.",
			wantErr: nil,
		},
	}
	for _, tc := range tests {
		// visit all elements
		for _, intro := range tc.intros {
			gotErr := tc.visitor.VisitIntro(intro)
			if gotErr != tc.wantErr {
				t.Fatalf("GotErr: %s\nWantErr: %s\n", gotErr, tc.wantErr)
			}
		}
		for _, clip := range tc.clips {
			gotErr := tc.visitor.VisitClip(clip)
			if gotErr != tc.wantErr {
				t.Fatalf("GotErr: %s\nWantErr: %s\n", gotErr, tc.wantErr)
			}
		}
		for _, outro := range tc.outros {
			gotErr := tc.visitor.VisitOutro(outro)
			if gotErr != tc.wantErr {
				t.Fatalf("GotErr: %s\nWantErr: %s\n", gotErr, tc.wantErr)
			}
		}
		if tc.visitor.String() != tc.want {
			t.Fatalf("Got: %s\nWanted: %s", tc.visitor.String(), tc.want)
		}
	}
}
