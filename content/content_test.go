package content_test

import (
	"mcsweeney/content"
	"testing"
)

// TestVideoAppend calls Video.Append() and provides a component interface,
// checking for valid return values.
func TestVideoAppend(t *testing.T) {
	tests := []struct {
		name      string
		component content.Component
		video     *content.Video
		wantErr   error
	}{
		{
			name:      "Append leaf component to empty composite.",
			component: &content.Intro{},
			video:     &content.Video{},
			wantErr:   nil,
		},
		{
			name:      "Append composite component to empty composite.",
			component: &content.Video{},
			video:     &content.Video{},
			wantErr:   nil,
		},
	}
	for _, tc := range tests {
		gotErr := tc.video.Append(tc.component)
		if gotErr != tc.wantErr {
			t.Fatalf("Got: %v\nWanted: %v", gotErr, tc.wantErr)
		}
	}
}

// TestVideoPrepend calls Video.Prepend() and provides a component interface,
// checking for valid return values.
func TestVideoPrepend(t *testing.T) {
	tests := []struct {
		name      string
		component content.Component
		video     *content.Video
		wantErr   error
	}{
		{
			name:      "Prepend leaf component to empty composite.",
			component: &content.Intro{},
			video:     &content.Video{},
			wantErr:   nil,
		},
		{
			name:      "Prepend composite component to empty composite.",
			component: &content.Video{},
			video:     &content.Video{},
			wantErr:   nil,
		},
	}
	for _, tc := range tests {
		gotErr := tc.video.Prepend(tc.component)
		if gotErr != tc.wantErr {
			t.Fatalf("Got: %v\nWanted: %v", gotErr, tc.wantErr)
		}
	}
}

/*
func TestCompositeRemove

func TestCompositeChild

func TestComponentAccept
*/

// The component interface defines the common methods of all children in the composite
// content structure.
// The interface is defined by the following methods
//
// Accept(v *Visitor)
