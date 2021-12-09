package share

import (
	"errors"
	"fmt"
	"gotest.tools/assert"
	"mcsweeney/config"
	"reflect"
	"testing"
)

// TODO: Enforce valid db? This is a question of design
func TestNewContentGetter(t *testing.T) {
	testCases := []struct {
		name       string
		config     config.Config
		wantGetter ContentGetter
		wantErr    error
	}{
		{
			name: "Success: New TwitchGetter",
			config: config.Config{
				Source: "twitch",
			},
			wantGetter: new(TwitchGetter),
			wantErr:    nil,
		},
		{
			name: "Failure: Getter not found",
			config: config.Config{
				Source: "not-implemented",
			},
			wantGetter: nil,
			wantErr:    fmt.Errorf("No such content-getter 'not-implemented'"),
		},
	}

	for _, tc := range testCases {
		gotGetter, gotErr := NewContentGetter(tc.config, nil)
		assert.Equal(t, reflect.TypeOf(gotGetter), reflect.TypeOf(tc.wantGetter))
		assert.Equal(t, errors.Unwrap(gotErr), errors.Unwrap(tc.wantErr))
	}
}
