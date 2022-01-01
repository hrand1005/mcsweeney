package content_test

import (
	"mcsweeney/content"
	"testing"
)

// TestOverlayGeneratorVisitIntro calls OverlayGenerator.VisitIntro and
// provides an Intro element, checking that the OverlayGenerator's String()
// method returns the expected string representation.
// TODO: Decide what overlays should be generated for intros
// Currently VisitIntro does not change OverlayGenerator's string
// representation. However, it should still have effects on OverlayGenerator's
// internal state, such as incrementing the OverlayGenerator's internal cursor.
func TestOverlayGeneratorVisitIntro(t *testing.T) {
	tests := []struct {
		name    string
		intro   *content.Intro
		visitor *content.OverlayGenerator
		want    string
	}{
		{
			name: "Nominal overlay for intro, empty string.",
			intro: &content.Intro{
				Duration: 1.0,
			},
			visitor: &content.OverlayGenerator{},
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

// TestOverlayGeneratorVisitOutro calls OverlayGenerator.VisitOutro and
// provides an Outro element, checking that the OverlayGenerator's String()
// method returns the expected string representation.
// TODO: Decide what overlays should be generated for intros
// Currently VisitIntro does not change OverlayGenerator's string
// representation. However, it should still have effects on OverlayGenerator's
// internal state, such as incrementing the OverlayGenerator's internal cursor.
func TestOverlayGeneratorVisitOutro(t *testing.T) {
	tests := []struct {
		name    string
		outro   *content.Outro
		visitor *content.OverlayGenerator
		want    string
	}{
		{
			name: "Nominal overlay for outro, empty string.",
			outro: &content.Outro{
				Duration: 1.0,
			},
			visitor: &content.OverlayGenerator{},
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

// TestOverlayGeneratorVisitClip calls OverlayGenerator.VisitClip and
// provides an Clip element, checking that the OverlayGenerator's String()
// method returns the expected string representation.
func TestOverlayGeneratorVisitClip(t *testing.T) {
	tests := []struct {
		name    string
		clip    *content.Clip
		visitor *content.OverlayGenerator
		want    string
	}{
		{
			name: "Overlay clip with background set in OverlayGenerator.",
			clip: &content.Clip{
				Title:       "TestTitle",
				Broadcaster: "TestBroadcaster",
				Duration:    1.0,
			},
			visitor: &content.OverlayGenerator{
				Background: "/path/to/Background.png",
				Font:       "/path/to/Font.ttf",
			},
			want: "-i /path/to/Background.png -filter_complex overlay=x='if(lt(t,0.000000)," +
				"NAN,if(lt(t,0.155000),-w+(t-0.000000)*2000.000000,if(lt(t,3.155000)," +
				"-w+310.000000,-w+310.000000-(t-3.155000)*2000.000000)))':y=470," +
				"drawtext=fontfile=/path/to/Font.ttf:text=TestTitle\nTestBroadcaster:" +
				"fontsize=26:fontcolor=ffffff:alpha='if(lt(t,0.500000),0,if(lt(t,1.000000)," +
				"(t-0.500000)/1,if(lt(t,2.500000),1,if(lt(t,3.000000),(1-(t-2.500000))/1," +
				"0))))':x=20:y=500",
		},
	}
	for _, tc := range tests {
		tc.visitor.VisitClip(tc.clip)
		if tc.visitor.String() != tc.want {
			t.Fatalf("Got: %s\nWanted: %s", tc.visitor.String(), tc.want)
		}
	}
}

// TestOverlayGeneratorVisitMany calls multiple visit methods on
// OverlayGenerator in sequence and checks that the resulting string returned by
// String() properly represents the visited element sequence.
//func TestOverlayGeneratorVisitMany(t *testing.T) {

//}
