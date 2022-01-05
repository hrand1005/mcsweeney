package content_test

import (
	"mcsweeney/content"
	"testing"
)

// TestOverlayerVisitIntro calls Overlayer.VisitIntro and
// provides an Intro element, checking that the Overlayer's String()
// method returns the expected string representation.
func TestOverlayerVisitIntro(t *testing.T) {
	tests := []struct {
		name    string
		intro   *content.Intro
		visitor *content.Overlayer
		want    string
		wantErr error
	}{
		{
			name: "Nominal overlay for intro, empty string.",
			intro: &content.Intro{
				Duration: 1.0,
			},
			visitor: &content.Overlayer{},
			want:    "",
			wantErr: nil,
		},
	}
	for _, tc := range tests {
		tc.visitor.VisitIntro(tc.intro)
		if tc.visitor.String() != tc.want {
			t.Fatalf("Got: %s\nWanted: %s", tc.visitor.String(), tc.want)
		}
	}
}

// TestOverlayerVisitOutro calls Overlayer.VisitOutro and
// provides an Outro element, checking that the Overlayer's String()
// method returns the expected string representation.
func TestOverlayerVisitOutro(t *testing.T) {
	tests := []struct {
		name    string
		outro   *content.Outro
		visitor *content.Overlayer
		want    string
		wantErr error
	}{
		{
			name: "Nominal overlay for outro, empty string.",
			outro: &content.Outro{
				Duration: 1.0,
			},
			visitor: &content.Overlayer{},
			want:    "",
			wantErr: nil,
		},
	}
	for _, tc := range tests {
		gotErr := tc.visitor.VisitOutro(tc.outro)
		if tc.visitor.String() != tc.want || gotErr != tc.wantErr {
			t.Fatalf("Got: %s\nWanted: %s\nGotErr: %s\nWantErr: %s\n", tc.visitor.String(), tc.want, gotErr, tc.wantErr)
		}
	}
}

// TestOverlayerVisitClip calls Overlayer.VisitClip and
// provides an Clip element, checking that the Overlayer's String()
// method returns the expected string representation.
func TestOverlayerVisitClip(t *testing.T) {
	tests := []struct {
		name    string
		clip    *content.Clip
		visitor *content.Overlayer
		want    string
		wantErr error
	}{
		{
			name: "Overlay clip with background set in Overlayer.",
			clip: &content.Clip{
				Title:       "TestTitle",
				Broadcaster: "TestBroadcaster",
				Duration:    1.0,
			},
			visitor: &content.Overlayer{
				Background: "/path/to/Background.png",
				Font:       "/path/to/Font.ttf",
			},
			want: "-i /path/to/Background.png -filter_complex overlay=x='if(lt(t,0.000000)," +
				"NAN,if(lt(t,0.275000),-w+(t-0.000000)*2000.000000,if(lt(t,3.275000)," +
				"-w+550.000000,-w+550.000000-(t-3.275000)*2000.000000)))':y=(H-h)/2," +
				"drawtext=fontfile=/path/to/Font.ttf:text=TestTitle\nTestBroadcaster:" +
				"fontsize=52:fontcolor=ffffff:alpha='if(lt(t,0.500000),0,if(lt(t,1.000000)," +
				"(t-0.500000)/1,if(lt(t,2.500000),1,if(lt(t,3.000000),(1-(t-2.500000))/1," +
				"0))))':x=20:y=(h-text_h)/2",
			wantErr: nil,
		},
		{
			name: "No overlay args returns ErrEmptyOverlay",
			clip: &content.Clip{
				Duration: 1.0,
			},
			visitor: &content.Overlayer{
				Font: "/path/to/Font.ttf",
			},
			want:    "",
			wantErr: content.ErrEmptyOverlay,
		},
		{
			name: "No duration returns ErrNoDuration",
			clip: &content.Clip{
				Duration: 0.0,
			},
			visitor: &content.Overlayer{
				Font: "/path/to/Font.ttf",
			},
			want:    "",
			wantErr: content.ErrNoDuration,
		},
	}
	for _, tc := range tests {
		gotErr := tc.visitor.VisitClip(tc.clip)
		if tc.visitor.String() != tc.want || gotErr != tc.wantErr {
			t.Fatalf("Got: %s\nWanted: %s\nGotErr: %s\nWantErr: %s\n", tc.visitor.String(), tc.want, gotErr, tc.wantErr)
		}
	}
}

// TestOverlayerVisitMany calls multiple visit methods on
// Overlayer in sequence and checks that the resulting string returned by
// String() properly represents the visited element sequence.
func TestOverlayerVisitMany(t *testing.T) {
	tests := []struct {
		name    string
		intros  []*content.Intro
		clips   []*content.Clip
		outros  []*content.Outro
		visitor *content.Overlayer
		want    string
		wantErr error
	}{
		{
			name: "One Intro, one Clip, one Outro in sequence, success case.",
			intros: []*content.Intro{
				&content.Intro{
					Duration: 1.0,
				},
			},
			clips: []*content.Clip{
				&content.Clip{
					Title:       "TestTitle",
					Broadcaster: "TestBroadcaster",
					Duration:    1.0,
				},
			},
			outros: []*content.Outro{
				&content.Outro{
					Duration: 1.0,
				},
			},
			visitor: &content.Overlayer{
				Background: "/path/to/Background.png",
				Font:       "/path/to/Font.ttf",
			},
			want: "-i /path/to/Background.png -filter_complex overlay=x='if(lt(t,1.000000)," +
				"NAN,if(lt(t,1.275000),-w+(t-1.000000)*2000.000000,if(lt(t,4.275000)," +
				"-w+550.000000,-w+550.000000-(t-4.275000)*2000.000000)))':y=(H-h)/2," +
				"drawtext=fontfile=/path/to/Font.ttf:text=TestTitle\nTestBroadcaster:" +
				"fontsize=52:fontcolor=ffffff:alpha='if(lt(t,1.500000),0,if(lt(t,2.000000)," +
				"(t-1.500000)/1,if(lt(t,3.500000),1,if(lt(t,4.000000),(1-(t-3.500000))/1," +
				"0))))':x=20:y=(h-text_h)/2",
			wantErr: nil,
		},
		{
			name:   "Multiple clips success case.",
			intros: []*content.Intro{},
			clips: []*content.Clip{
				&content.Clip{
					Title:       "TestTitle",
					Broadcaster: "TestBroadcaster",
					Duration:    1.0,
				},
				&content.Clip{
					Title:       "TestTitle2",
					Broadcaster: "TestBroadcaster2",
					Duration:    1.0,
				},
			},
			outros: []*content.Outro{},
			visitor: &content.Overlayer{
				Background: "/path/to/Background.png",
				Font:       "/path/to/Font.ttf",
			},
			want: "-i /path/to/Background.png -i /path/to/Background.png -filter_complex " +
				"overlay=x='if(lt(t,0.000000)," +
				"NAN,if(lt(t,0.275000),-w+(t-0.000000)*2000.000000,if(lt(t,3.275000)," +
				"-w+550.000000,-w+550.000000-(t-3.275000)*2000.000000)))':y=(H-h)/2," +
				"drawtext=fontfile=/path/to/Font.ttf:text=TestTitle\nTestBroadcaster:" +
				"fontsize=52:fontcolor=ffffff:alpha='if(lt(t,0.500000),0,if(lt(t,1.000000)," +
				"(t-0.500000)/1,if(lt(t,2.500000),1,if(lt(t,3.000000),(1-(t-2.500000))/1," +
				"0))))':x=20:y=(h-text_h)/2," +
				"overlay=x='if(lt(t,1.000000)," +
				"NAN,if(lt(t,1.291000),-w+(t-1.000000)*2000.000000,if(lt(t,4.291000)," +
				"-w+582.000000,-w+582.000000-(t-4.291000)*2000.000000)))':y=(H-h)/2," +
				"drawtext=fontfile=/path/to/Font.ttf:text=TestTitle2\nTestBroadcaster2:" +
				"fontsize=52:fontcolor=ffffff:alpha='if(lt(t,1.500000),0,if(lt(t,2.000000)," +
				"(t-1.500000)/1,if(lt(t,3.500000),1,if(lt(t,4.000000),(1-(t-3.500000))/1," +
				"0))))':x=20:y=(h-text_h)/2",
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
