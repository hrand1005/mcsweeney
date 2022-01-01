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
    tests := []struct{
        name string
        intro *content.Intro
        visitor *content.OverlayGenerator
        want string
    }{
        {
            name: "Nominal overlay for intro, empty string.",
            intro: &content.Intro{
                Duration: 1.0,
            },
            visitor: &content.OverlayGenerator{},
            want: "",
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
    tests := []struct{
        name string
        outro *content.Outro
        visitor *content.OverlayGenerator
        want string
    }{
        {
            name: "Nominal overlay for outro, empty string.",
            outro: &content.Outro{
                Duration: 1.0,
            },
            visitor: &content.OverlayGenerator{},
            want: "",
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
//func TestOverlayGeneratorVisitClip(t *testing.T) {

//}

// TestOverlayGeneratorVisitMany calls multiple visit methods on
// OverlayGenerator in sequence and checks that the resulting string returned by
// String() properly represents the visited element sequence.
//func TestOverlayGeneratorVisitMany(t *testing.T) {

//}
