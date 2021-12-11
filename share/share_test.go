package share_test

import (
	"errors"
	"fmt"
	"gotest.tools/assert"
	"mcsweeney/config"
	"mcsweeney/share"
	"reflect"
	"testing"
)

func TestNewContentSharer(t *testing.T) {
	testCases := []struct {
		name       string
		config     config.Config
		wantSharer share.ContentSharer
		wantErr    error
	}{
		{
			name: "Success: New YoutubeSharer",
			config: config.Config{
				Destination: "youtube",
			},
			wantSharer: new(share.YoutubeSharer),
			wantErr:    nil,
		},
		{
			name: "Failure: sharer not found",
			config: config.Config{
				Destination: "not-implemented",
			},
			wantSharer: nil,
			wantErr:    fmt.Errorf("No such content-sharer 'not-implemented'"),
		},
	}

	for _, tc := range testCases {
		gotSharer, gotErr := share.NewContentSharer(tc.config, "path")
		assert.Equal(t, reflect.TypeOf(gotSharer), reflect.TypeOf(tc.wantSharer))
		assert.Equal(t, errors.Unwrap(gotErr), errors.Unwrap(tc.wantErr))
	}
}
